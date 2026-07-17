package config

import (
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"

	"YH-FireWall/internal/pkg/lfile"
)

type Manager struct {
	mutex sync.RWMutex
	// 文件只读一次，写入多次
	file *lfile.LockedFile
	// 文件内容由conten控制
	content []byte
}

func New(path string) (*Manager, error) {
	file, err := lfile.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	// 读取数据流
	buf, err := file.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	// 当不为空时，验证buf合法性
	if len(buf) > 0 {
		if err := toml.Unmarshal(buf, new(Config)); err != nil {
			return nil, fmt.Errorf("failed to decode config file: %w", err)
		}
	}
	// 返回
	return &Manager{
		file:    file,
		content: buf,
	}, nil
}

func (m *Manager) Load() Config {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	// 覆盖默认配置
	cfg := DefaultConfig()
	// 直接认定不会出错
	_ = toml.Unmarshal(m.content, &cfg)
	return *cfg
}

func (m *Manager) Read() string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return string(m.content)
}

func (m *Manager) Write(data string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 验证是否满足格式
	if _, err := toml.Decode(data, new(Config)); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}
	// 尝试写入
	content := []byte(data)
	if err := m.file.Write(content); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	// 写入成功
	m.content = content
	return nil
}

func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.file.Close()
}
