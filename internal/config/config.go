package config

import (
	"YH-FireWall/internal/pkg/flock"
	"YH-FireWall/internal/queue"
	"YH-FireWall/internal/rtable"
	"YH-FireWall/internal/server/cmdserver"
	"YH-FireWall/internal/server/webserver"
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	lf      *flock.LockedFile
	mutex   sync.RWMutex
	content []byte
)

type Config struct {
	Queue *queue.Config     `json:"queue"`
	Web   *webserver.Config `json:"web"`
	Cmd   *cmdserver.Config `json:"cmd"`
	Rule  *rtable.Config    `json:"rule"`
}

func (c *Config) Read(buf []byte) error {
	return yaml.Unmarshal(buf, c)
}

func (c *Config) Write() []byte {
	return nil
}

func Default() *Config {
	return &Config{
		Queue: queue.DefaultConfig(),
		Web:   webserver.DefaultConfig(),
		Cmd:   cmdserver.DefaultConfig(),
		Rule:  rtable.DefaultConfig(),
	}
}

func Init(configPath string) (err error) {
	lf, err = flock.Open(configPath)
	if err != nil {
		return err
	}
	content, err = lf.Read()
	return err
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
	cfg = Default()
	if len(buf) == 0 {
		return cfg, nil
	}
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
	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	if err := lf.Write(buf); err != nil {
		return err
	}
	content = buf
	return nil
}

func Close() error {
	return lf.Close()
}
