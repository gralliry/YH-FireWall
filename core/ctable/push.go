package ctable

import (
	"YH-FireWall/core/connection"
	"YH-FireWall/core/iface"
	"github.com/google/gopacket/layers"
	nnet "github.com/shirou/gopsutil/v4/net"
	nprocess "github.com/shirou/gopsutil/v4/process"
	"net"
)

func Push(
	family uint32, proto layers.IPProtocol,
	srcIP net.IP, srcPort uint16, dstIP net.IP, dstPort uint16,
	inDev, outDev *uint32,
) bool {
	mutex.Lock()
	defer mutex.Unlock()
	// 添加连接键
	lkey := connection.MakeKey(proto, srcIP, srcPort, dstIP, dstPort)
	rkey := connection.MakeKey(proto, dstIP, dstPort, srcIP, srcPort)
	conn, exists := table[lkey]
	if exists {
		// 检测连接状态
		isClosed, isExpired := conn.Status()
		switch {
		case isClosed && isExpired:
			// 如果被关闭 且 已过期
			delete(table, lkey)
			delete(table, rkey)
			delete(namespcae, conn.Id())
			return true
		case isClosed && !isExpired:
			// 表示该连接任然未过期
			return false
		default:
			// 如果没有过期，更新
			conn.UpdateByPush()
			return true
		}
	}
	// 添加连接 // 获取方向
	switch {
	case inDev != nil && outDev != nil:
		// 转发模式
		name, err := net.InterfaceByIndex(int(*outDev))
		if err != nil {
			return false
		}
		conn = connection.NewByPush(family, proto, srcIP, srcPort, connection.Forward, dstIP, dstPort, name.Name)
	case inDev != nil:
		// 入口模式 // 源是外部连接
		name, err := net.InterfaceByIndex(int(*inDev))
		if err != nil {
			return false
		}
		conn = connection.NewByPush(family, proto, dstIP, dstPort, connection.Inbound, srcIP, srcPort, name.Name)
	case outDev != nil:
		// 出口模式 // 源是内部连接
		name, err := net.InterfaceByIndex(int(*outDev))
		if err != nil {
			return false
		}
		conn = connection.NewByPush(family, proto, srcIP, srcPort, connection.Outbound, dstIP, dstPort, name.Name)
	default:
		// 未知数据，直接停止
		return false
	}
	// 写入表
	table[lkey] = conn
	table[rkey] = conn
	namespcae[conn.Id()] = conn
	//
	return true
}

func pushByProcess() {
	// 获取进程列表
	processList, err := nprocess.Processes()
	if err != nil {
		return
	}
	// 获取进程信息
	processMap := make(map[int32]*nprocess.Process)
	for _, pc := range processList {
		processMap[pc.Pid] = pc
	}
	// 获取网络连接列表
	connectionList, err := nnet.Connections("inet")
	if err != nil {
		return
	}
	for _, conn := range connectionList {
		// 2: ipv4, 10: ipv6
		if conn.Family != 2 && conn.Family != 10 {
			continue
		}
		// 1: tcp, 2: udp
		var protocol layers.IPProtocol
		switch conn.Type {
		case 1:
			protocol = layers.IPProtocolTCP
		case 2:
			protocol = layers.IPProtocolUDP
		default:
			continue
		}
		// ip and port
		localIP, remoteIP := net.ParseIP(conn.Laddr.IP), net.ParseIP(conn.Raddr.IP)
		localPort, remotePort := uint16(conn.Laddr.Port), uint16(conn.Raddr.Port)
		//  key
		lkey := connection.MakeKey(protocol, localIP, localPort, remoteIP, remotePort)
		//
		if connect, exists := table[lkey]; exists {
			// 添加连接进程信息
			if connect.IsProcessInfoEmpty() {
				if pc, e := processMap[conn.Pid]; e {
					// 获取进程信息
					exe, _ := pc.Exe()
					name, _ := pc.Name()
					username, _ := pc.Username()
					cmd, _ := pc.Cmdline()
					// 更新进程信息
					connect.UpdateByProcess(conn.Fd, conn.Pid, exe, name, cmd, username, conn.Status)
					continue
				}
			}
		} else {
			// 查找网卡
			localInterface, remoteInterface := iface.FindNameByIp(localIP), iface.FindNameByIp(remoteIP)
			// 添加连接 // 获取方向
			switch {
			case localInterface != "" && remoteInterface != "":
				// 转发模式
				connect = connection.NewByPush(conn.Family, protocol, localIP, localPort, connection.Forward, remoteIP, remotePort, localInterface)
			case localInterface != "":
				// 入口模式 // 源是外部连接
				connect = connection.NewByPush(conn.Family, protocol, localIP, localPort, connection.Inbound, remoteIP, remotePort, localInterface)
			case remoteInterface != "":
				// 出口模式 // 源是内部连接
				connect = connection.NewByPush(conn.Family, protocol, localIP, localPort, connection.Outbound, remoteIP, remotePort, remoteInterface)
			default:
				// 转发模式
				connect = connection.NewByPush(conn.Family, protocol, localIP, localPort, connection.Forward, remoteIP, remotePort, "")
			}
			// 更新连接进程信息
			if pc, e := processMap[conn.Pid]; e {
				// 获取进程信息
				exe, _ := pc.Exe()
				name, _ := pc.Name()
				username, _ := pc.Username()
				cmd, _ := pc.Cmdline()
				// 更新进程信息
				connect.UpdateByProcess(conn.Fd, conn.Pid, exe, name, cmd, username, conn.Status)
			}
			// 添加连接信息
			rkey := connection.MakeKey(protocol, remoteIP, remotePort, localIP, localPort)
			// 写入表
			table[lkey] = connect
			table[rkey] = connect
			namespcae[connect.Id()] = connect
		}
	}
}
