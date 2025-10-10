package config

import (
	"YH-FireWall/core/rule"
	"time"
)

const Version = "1.0.0"

type Config struct {
	LastUpdateDate string        `json:"last_update_date"`
	Rules          []rule.Config `json:"rules"`
	Queue          Queue         `json:"queue"`
}

func DefaultConfig() *Config {
	return &Config{
		Queue: Queue{
			Num:    0,
			Accept: true,
		},
		LastUpdateDate: time.Now().Format("2006-01-02 15:04:05"),
		Rules:          []rule.Config{},
	}
}
