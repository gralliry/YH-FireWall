package ctable

import (
	"YH-FireWall/internal/itable"
	"YH-FireWall/internal/model/flow"
	"net"

	"github.com/google/gopacket/layers"
	nnet "github.com/shirou/gopsutil/v4/net"
	nprocess "github.com/shirou/gopsutil/v4/process"
)

func (m *Manager) Push(flow *flow.Flow) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// 添加连接键
	key := flow.Key()
	//
	if conn, exists := m.table.Get(key); exists {
		// 检测连接状态
		isClosed, isExpired := conn.Closed(), conn.Expired()
		switch {
		case isClosed && isExpired:
			// 如果被关闭 且 已过期
			m.table.Del(key)
			return true
		case isClosed && !isExpired:
			// 表示该连接任然未过期
			return false
		default:
			// 如果没有过期，更新
			conn.Active()
			return true
		}
	} else {
		// 添加连接 // 获取方向
		connect := conn.New(flow)
		// 写入表
		table.Set(connect, connect.LKey(), connect.RKey(), connect.Id())
		//
		return true
	}
}

func pushByProcess() {
	// 获取进程信息
	processMap := make(map[int32]*nprocess.Process)
	if processList, err := nprocess.Processes(); err == nil {
		for _, pc := range processList {
			processMap[pc.Pid] = pc
		}
	} else {
		return
	}
	// 获取网络连接列表
	connectionList, err := nnet.Connections("inet")
	if err != nil {
		return
	}
	// 获取网卡ip映射
	for _, conn := range connectionList {
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
		flow := flow.Flow{
			Family:   uint8(conn.Family),
			Protocol: protocol,

			SrcIP:   net.ParseIP(conn.Laddr.IP),
			SrcPort: uint16(conn.Laddr.Port),
			DstIP:   net.ParseIP(conn.Raddr.IP),
			DstPort: uint16(conn.Raddr.Port),
		}
		// 查找网卡
		if inDev, exist := itable.LookupByIp(flow.SrcIP.String()); exist {
			flow.InDev = uint32(inDev)
		}
		if outDev, exist := itable.LookupByIp(flow.DstIP.String()); exist {
			flow.OutDev = uint32(outDev)
		}
		//  key
		lkey := flow.LKey()
		//
		connect, exists := table[lkey]
		if !exists {
			connect = New(&flow)
			// 添加连接信息
			rkey := connect.RKey()
			// 写入表
			table[lkey] = connect
			table[rkey] = connect
			namespace[connect.Id()] = connect
		}
		// 添加连接进程信息
		if pc, e := processMap[conn.Pid]; e {
			// 更新进程信息
			connect.SetProcess(pc)
		}
	}
}
