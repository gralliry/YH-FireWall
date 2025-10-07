package repo

import (
	"YH-FireWall/internal/core/rule"
	"encoding/json"
	"fmt"
	"os"
	"syscall"
	"time"
)

type Config struct {
	QueueNum      uint16        `json:"queue_num"`
	UpdateDate    string        `json:"update_date"`
	Rules         []rule.Config `json:"rules"`
	DefaultAccept bool          `json:"default_accept"`

	filepath string
	file     *os.File
}

func Load(filepath string) (*Config, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}

	if err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("锁定文件失败: %w", err)
	}

	cfg := Config{
		QueueNum:      0,
		UpdateDate:    time.Now().Format("2006-01-02 15:04:05"),
		Rules:         []rule.Config{},
		DefaultAccept: true,

		filepath: filepath,
		file:     f,
	}

	decoder := json.NewDecoder(f)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Store() error {
	if _, err := c.file.Seek(0, 0); err != nil {
		return fmt.Errorf("重置文件指针失败: %w", err)
	}
	if err := c.file.Truncate(0); err != nil {
		return fmt.Errorf("清空文件失败: %w", err)
	}

	encoder := json.NewEncoder(c.file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("写入配置失败: %w", err)
	}

	return c.file.Sync()
}

func (c *Config) Close() error {
	err := syscall.Flock(int(c.file.Fd()), syscall.LOCK_UN)
	_ = c.file.Close()
	return err
}
