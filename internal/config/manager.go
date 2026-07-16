package config

import (
	"YH-FireWall/internal/pkg/lfile"
	"encoding/json"
	"fmt"
	"sync"
	"unsafe"
)

type Manager struct {
	mutex  sync.RWMutex
	file   *lfile.LockedFile
	config *Config
}

func New(path string) (*Manager, error) {
	lf, err := lfile.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	buf, err := lf.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var cfg *Config = DefaultConfig()
	if len(buf) > 0 {
		if err := json.Unmarshal(buf, &cfg); err != nil {
			return nil, fmt.Errorf("failed to decode config file: %w", err)
		}
	}
	return &Manager{
		file:   lf,
		config: cfg,
	}, nil
}

func (m *Manager) Load() Config {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return *m.config
}

func (m *Manager) Read() (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	buf, err := m.file.Read()
	if err != nil {
		return "", err
	}
	return string(buf), err
}

func (m *Manager) Write(data string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	buf := unsafe.Slice(unsafe.StringData(data), len(data))

	var cfg Config
	if err := json.Unmarshal(buf, &cfg); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}
	if err := m.file.Write(buf); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	m.config = &cfg
	return nil
}

func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.file.Close()
}
