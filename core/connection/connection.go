package connection

import (
	"fmt"
	nnet "github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
	"net"
	"os/exec"
	"strconv"
)

type Connection struct {
	// Family uint32 `json:"family"` // 2: ipv4, 10: ipv6
	Fd       uint32 `json:"fd"`
	Type     uint32 `json:"type"` // 1: tcp, 2: udp
	Pid      int32  `json:"pid"`
	Exe      string `json:"exe"`
	Name     string `json:"name"`
	Cmd      string `json:"cmd"`
	Username string `json:"username"`
	Laddr    string `json:"localaddr"`
	Raddr    string `json:"remoteaddr"`
	Status   string `json:"status"`
}

func GetAll() ([]Connection, error) {
	// 获取进程列表
	processList, err := process.Processes()
	if err != nil {
		return nil, err
	}
	// 获取进程信息
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
			Fd:     conn.Fd,
			Type:   conn.Type,
			Pid:    conn.Pid,
			Laddr:  net.JoinHostPort(conn.Laddr.IP, strconv.Itoa(int(conn.Laddr.Port))),
			Raddr:  net.JoinHostPort(conn.Raddr.IP, strconv.Itoa(int(conn.Raddr.Port))),
			Status: conn.Status,
		}
		nc.Exe, _ = pc.Exe()
		nc.Name, _ = pc.Name()
		nc.Username, _ = pc.Username()
		nc.Cmd, _ = pc.Cmdline()
		connections = append(connections, nc)
	}
	return connections, nil
}

func Close(pid int32, fd uint32) error {
	call := fmt.Sprintf("call (int) close(%d)", fd)
	cmd := exec.Command(
		"gdb",
		"-p", strconv.Itoa(int(pid)),
		"--batch",
		"-ex", call,
		"-ex", "detach",
		"-ex", "quit",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gdb failed: %w, output: %s", err, out)
	}
	fmt.Println(string(out)) // 打印 gdb 执行结果，调试用
	return nil
}
