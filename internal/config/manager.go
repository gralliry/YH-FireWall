package config

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"

	"YH-FireWall/internal/pkg/cfile"
)

type Manager struct {
	mutex sync.RWMutex
	//
	path string
	// 文件只读一次，写入多次
	file   *cfile.CacheFile
	logger *slog.Logger
}

func New(path string, logger *slog.Logger) (*Manager, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}
	file, err := cfile.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	// 读取数据流
	buf := file.Read()
	// 当不为空时，验证buf合法性
	// if len(buf) > 0 {} // 当为json时，期待{}，这里会报错
	if err := toml.Unmarshal(buf, new(Config)); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}
	// 返回
	return &Manager{
		path:   absPath,
		file:   file,
		logger: logger,
	}, nil
}

func (m *Manager) Load() Config {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	// 覆盖默认配置
	cfg := DefaultConfig()
	if err := toml.Unmarshal(m.file.Read(), &cfg); err != nil {
		m.logger.Warn("config: decode failed, using defaults", slog.String("error", err.Error()))
	}
	return *cfg
}

func (m *Manager) Path() string {
	return m.path
}

func (m *Manager) Read() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return string(m.file.Read())
}

func (m *Manager) Write(data string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	content := []byte(data)
	// 验证是否满足格式
	if err := toml.Unmarshal(content, new(Config)); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}
	// 尝试写入
	if err := m.file.Write(content); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.file.Close()
	return nil
}
