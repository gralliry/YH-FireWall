package config

import (
	"YH-FireWall/internal/queue"
	"YH-FireWall/internal/rtable"
	"YH-FireWall/internal/server/cmdserver"
	"YH-FireWall/internal/server/webserver"
)

type Config struct {
	Version string           `toml:"version"`
	Queue   queue.Config     `toml:"queue"`
	Web     webserver.Config `toml:"web"`
	Cmd     cmdserver.Config `toml:"cmd"`
	Rule    rtable.Config    `toml:"rule"`
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
