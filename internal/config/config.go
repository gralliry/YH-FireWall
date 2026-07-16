package config

import (
	"YH-FireWall/internal/queue"
	"YH-FireWall/internal/rtable"
	"YH-FireWall/internal/server/cmdserver"
	"YH-FireWall/internal/server/webserver"
)

type Config struct {
	Version string           `json:"version"`
	Queue   queue.Config     `json:"queue"`
	Web     webserver.Config `json:"web"`
	Cmd     cmdserver.Config `json:"cmd"`
	Rule    rtable.Config    `json:"rule"`
}

func DefaultConfig() *Config {
	return &Config{
		Version: Version,
		Queue:   *queue.DefaultConfig(),
		Web:     *webserver.DefaultConfig(),
		Cmd:     *cmdserver.DefaultConfig(),
		Rule:    *rtable.DefaultConfig(),
	}
}
