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
	Family uint32 `json:"family"`
	Type   uint32 `json:"type"`
	Laddr  string `json:"localaddr"`
	Raddr  string `json:"remoteaddr"`
	Status string `json:"status"`
	Pid    int32  `json:"pid"`
}

func connectionStat2Connection(cs nnet.ConnectionStat) Connection {
	return Connection{
		Family: cs.Family,
		Type:   cs.Type,
		Laddr:  net.JoinHostPort(cs.Laddr.IP, strconv.Itoa(int(cs.Laddr.Port))),
		Raddr:  net.JoinHostPort(cs.Raddr.IP, strconv.Itoa(int(cs.Raddr.Port))),
		Status: cs.Status,
		Pid:    cs.Pid,
	}
}

func GetConnections() ([]Connection, error) {
	// netstat -tunlp // ss -tunlp
	cs, err := nnet.Connections("inet")
	if err != nil {
		return nil, err
	}
	return fp.Map(cs, connectionStat2Connection), nil
}

type Process struct {
	Pid      int32  `json:"pid"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Cmd      string `json:"cmd"`
}

func GetProcesses() ([]Process, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}
	return fp.Map(procs, func(sp *process.Process) Process {
		p := Process{Pid: sp.Pid}
		p.Name, _ = sp.Name()
		p.Username, _ = sp.Username()
		p.Cmd, _ = sp.Cmdline()
		return p
	}), nil
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
