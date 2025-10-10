package config

import (
	"YH-FireWall/core/rule"
	"time"
)

const Version = "1.0.0"

type Queue struct {
	Num    uint16 `json:"num"`
	Accept bool   `json:"accept"`
}

type Web struct {
	Address           string `json:"address"`
	BasicAuthUser     string `json:"basic_auth_user"`
	BasicAuthPassword string `json:"basic_auth_password"`
	StaticDir         string `json:"static_dir"`
}

type Unix struct {
	Path string `json:"path"`
}

type Config struct {
	LastUpdateDate string        `json:"last_update_date"`
	Rules          []rule.Config `json:"rules"`
	Queue          Queue         `json:"queue"`
	Web            Web           `json:"web"`
	Unix           Unix          `json:"unix"`
}

func DefaultConfig() *Config {
	return &Config{
		Queue: Queue{
			Num:    0,
			Accept: true,
		},
		Web: Web{
			Address:           ":8080",
			BasicAuthUser:     "",
			BasicAuthPassword: "",
			StaticDir:         "front/dist",
		},
		Unix: Unix{
			Path: "/tmp/firewall.sock",
		},
		LastUpdateDate: time.Now().Format("2006-01-02 15:04:05"),
		Rules:          []rule.Config{},
	}
}
