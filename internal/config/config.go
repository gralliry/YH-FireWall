package config

import (
	"YH-FireWall/internal/rtable"
	"YH-FireWall/internal/server/cmdserver"
	"YH-FireWall/internal/server/webserver"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	file    *os.File
	mutex   sync.RWMutex
	content []byte
)

type Config struct {
	LastUpdateTime string           `json:"last_update_time"`
	QueueNo        uint16           `json:"queue_no"`
	Web            webserver.Config `json:"web"`
	Cmd            cmdserver.Config `json:"cmd"`
	RuleTable      rtable.Config    `json:"rule_table"`
}

func Default() *Config {
	return &Config{
		LastUpdateTime: time.Now().Format(time.RFC3339),
		QueueNo:        0,
		Web: webserver.Config{
			Enable:       true,
			Address:      ":8080",
			AuthUsername: "admin",
			AuthPassword: "admin",
			EnableCORS:   true,
		},
		Cmd: cmdserver.Config{
			Enable:     true,
			SocketPath: "/tmp/yfw.sock",
		},
		RuleTable: rtable.Config{
			Path:          "/etc/yfw/rule.json",
			DefaultAccept: true,
		},
	}
}

func Init(configPath string) (err error) {
	// 确保目录存在
	if err = os.MkdirAll(path.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	// 打开文件
	file, err = os.OpenFile(configPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	// 尝试独占锁（非阻塞）
	if err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = file.Close()
		return fmt.Errorf("failed to acquire exclusive lock: %w", err)
	}
	// 重置文件指针
	if _, err = file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stat: %w", err)
	}
	content = make([]byte, info.Size())
	if _, err = file.Read(content); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	return nil
}

func Read() []byte {
	mutex.RLock()
	defer mutex.RUnlock()
	// 拷贝内容
	result := make([]byte, len(content))
	copy(result, content)
	return result

}

func Load() (cfg *Config, err error) {
	buf := Read()
	// 默认配置
	cfg = Default()
	if err = yaml.Unmarshal(buf, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return cfg, nil
}

func Save(buf []byte) error {
	mutex.Lock()
	defer mutex.Unlock()
	//
	var cfg Config
	// 验证
	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	if _, err := file.Write(buf); err != nil {
		return fmt.Errorf("failed to write string to file: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}
	return nil
}

func Close() error {
	var errs []error
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN); err != nil {
		errs = append(errs, fmt.Errorf("failed to unlock file: %w", err))
	}
	if err := file.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close file: %w", err))
	}
	return errors.Join(errs...)
}
