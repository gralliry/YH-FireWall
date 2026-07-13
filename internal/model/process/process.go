package process

import (
	_process "github.com/shirou/gopsutil/v4/process"
)

type Process struct {
	Pid  int32
	Fd   int32
	Exe  string
	Name string
	Cmd  string
	User string
}

func New(pc *_process.Process) *Process {
	// 获取进程信息
	pid := pc.Pid
	fd, _ := pc.NumFDs()
	exe, _ := pc.Exe()
	name, _ := pc.Name()
	cmd, _ := pc.Cmdline()
	user, _ := pc.Username()
	return &Process{
		Pid:  pid,
		Fd:   fd,
		Exe:  exe,
		Name: name,
		Cmd:  cmd,
		User: user,
	}
}
