package repo

import (
	"YH-FireWall/internal/core/rule"
	"encoding/json"
	"os"
	"path"
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
	// 确保目录存在
	if err := os.MkdirAll(path.Dir(filepath), 0755); err != nil {
		return nil, err
	}
	// 打开文件
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	if err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, err
	}

	cfg := Config{
		QueueNum:      0,
		UpdateDate:    time.Now().Format("2006-01-02 15:04:05"),
		Rules:         []rule.Config{},
		DefaultAccept: true,

		filepath: filepath,
		file:     file,
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() == 0 {
		return &cfg, nil
	}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Store() error {
	if _, err := c.file.Seek(0, 0); err != nil {
		return err
	}
	if err := c.file.Truncate(0); err != nil {
		return err
	}

	encoder := json.NewEncoder(c.file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(c); err != nil {
		return err
	}

	return c.file.Sync()
}

func (c *Config) Close() error {
	err := syscall.Flock(int(c.file.Fd()), syscall.LOCK_UN)
	_ = c.file.Close()
	return err
}
