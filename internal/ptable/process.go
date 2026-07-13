package ptable

import (
	nprocess "github.com/shirou/gopsutil/v4/process"
)

type Info struct {
	Pid  int32  `json:"fd"`
	Fd   int32  `json:"pid"`
	Exe  string `json:"exe"`
	Name string `json:"name"`
	Cmd  string `json:"cmd"`
	User string `json:"user"`
}

func New(pc *nprocess.Process) *Info {
	// 获取进程信息
	pid := pc.Pid
	fd, _ := pc.NumFDs()
	exe, _ := pc.Exe()
	name, _ := pc.Name()
	cmd, _ := pc.Cmdline()
	user, _ := pc.Username()
	return &Info{
		Pid:  pid,
		Fd:   fd,
		Exe:  exe,
		Name: name,
		Cmd:  cmd,
		User: user,
	}
}

func Infos() ([]Info, error) {
	np, err := nprocess.Processes()
	if err != nil {
		return nil, err
	}
	info := make([]Info, len(np))
	for _, p := range np {
		info = append(info, *New(p))
	}
	return info, nil
}
