package config

import (
	"YH-FireWall/internal/core/group"
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	UpdateDate string         `json:"update_date"`
	Groups     []group.Config `json:"groups"`
}

// Load 读取 JSON 文件到内存，并初始化 Index
func Load(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		// 如果文件不存在，可以认为初始化为空
		if !os.IsNotExist(err) {
			return nil, err
		} else {
			// 创建文件
			return nil, err
		}
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, err
}

// Store 将内存规则写回 JSON 文件
func (c *Config) Store(filepath string) error {
	// 更新时间
	c.UpdateDate = time.Now().Format("2006-01-02 15:04:05")

	// 序列化
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}
