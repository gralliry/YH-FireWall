package proc

import (
	"github.com/shirou/gopsutil/v4/process"
)

type Info struct {
	Pid      int32  `json:"pid"`
	Exe      string `json:"exe"`
	Name     string `json:"name"`
	Cmdline  string `json:"cmdline"`
	Username string `json:"username"`
}

func New(pc *process.Process) *Info {
	// 获取进程信息
	var info Info
	info.Pid = pc.Pid
	info.Exe, _ = pc.Exe()
	info.Name, _ = pc.Name()
	info.Cmdline, _ = pc.Cmdline()
	info.Username, _ = pc.Username()
	return &info
}

func NewByPID(pid int32) *Info {
	pc, err := process.NewProcess(pid)
	if err != nil {
		return nil
	}
	return New(pc)
}
