package proc

import (
	"errors"

	"github.com/shirou/gopsutil/v4/process"
)

type Info struct {
	Pid      int32  `json:"pid"`
	Exe      string `json:"exe"`
	Name     string `json:"name"`
	Cmdline  string `json:"cmdline"`
	Username string `json:"username"`
}

func New(pc *process.Process) (*Info, error) {
	var info Info
	var errs []error
	var err error
	info.Pid = pc.Pid
	set := func(dst *string, fn func() (string, error)) {
		*dst, err = fn()
		if err != nil {
			errs = append(errs, err)
		}
	}
	set(&info.Exe, pc.Exe)
	set(&info.Name, pc.Name)
	set(&info.Cmdline, pc.Cmdline)
	set(&info.Username, pc.Username)
	return &info, errors.Join(errs...)
}

func NewByPID(pid int32) (*Info, error) {
	pc, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	return New(pc)
}
