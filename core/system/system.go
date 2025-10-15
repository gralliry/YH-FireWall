package system

import (
	"YH-FireWall/core/pkg/fp"
	nnet "github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
	"net"
	"strconv"
	"strings"
)

type Connection struct {
	// Family uint32 `json:"family"` // 2: ipv4, 10: ipv6
	Type     uint32 `json:"type"` // 1: tcp, 2: udp
	Pid      int32  `json:"pid"`
	Name     string `json:"name"`
	Cmd      string `json:"cmd"`
	Username string `json:"username"`
	Laddr    string `json:"localaddr"`
	Raddr    string `json:"remoteaddr"`
	Status   string `json:"status"`
}

func GetConnections() ([]Connection, error) {
	// 获取进程列表
	processList, err := process.Processes()
	if err != nil {
		return nil, err
	}
	processMap := make(map[int32]*process.Process)
	for _, pc := range processList {
		processMap[pc.Pid] = pc
	}
	// 获取网络连接列表
	connectionList, err := nnet.Connections("inet")
	if err != nil {
		return nil, err
	}
	// 遍历网络连接列表
	connections := make([]Connection, 0)
	for _, conn := range connectionList {
		pc, exists := processMap[conn.Pid]
		if !exists {
			continue
		}
		// 2: ipv4, 10: ipv6
		if conn.Family != 2 && conn.Family != 10 {
			continue
		}
		// 1: tcp, 2: udp
		if conn.Type != 1 && conn.Type != 2 {
			continue
		}
		// 构造连接
		nc := Connection{
			Type:   conn.Type,
			Pid:    conn.Pid,
			Laddr:  net.JoinHostPort(conn.Laddr.IP, strconv.Itoa(int(conn.Laddr.Port))),
			Raddr:  net.JoinHostPort(conn.Raddr.IP, strconv.Itoa(int(conn.Raddr.Port))),
			Status: conn.Status,
		}
		nc.Name, _ = pc.Name()
		nc.Username, _ = pc.Username()
		nc.Cmd, _ = pc.Cmdline()
		connections = append(connections, nc)
	}
	return connections, nil
}

type Interface struct {
	Index int      `json:"index"`
	Name  string   `json:"name"`
	MAC   string   `json:"mac"`
	MTU   int      `json:"mtu"`
	Flags []string `json:"flags"`
	Addrs []string `json:"addrs"`
}

func GetInterface() ([]Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	return fp.Map(interfaces, func(ifs net.Interface) Interface {
		oaddrs, _ := ifs.Addrs()
		return Interface{
			Index: ifs.Index,
			Name:  ifs.Name,
			MAC:   ifs.HardwareAddr.String(),
			MTU:   ifs.MTU,
			Flags: strings.Split(ifs.Flags.String(), "|"),
			Addrs: fp.Map(oaddrs, func(a net.Addr) string {
				return a.String()
			}),
		}
	}), nil
}
