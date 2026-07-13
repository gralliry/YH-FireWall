package webserver

import (
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var (
	mutex sync.RWMutex

	server  *fiber.App
	handler Handler
	config  *Config

	initialized bool
	running     bool
)

func Start(handler_ Handler, config_ *Config) error {
	mutex.Lock()
	defer mutex.Unlock()

	if !config_.Enable {
		initialized = false
		return nil
	}

	handler = handler_
	config = config_

	server_ := newServer()
	// 在 goroutine 中启动，确认监听成功后设置标记
	ready := make(chan struct{})
	go func() {
		initialized = true
		running = true
		close(ready)
		if err := server_.Listen(config.Address); err != nil {
			log.Println(err)
			running = false
		}
	}()
	<-ready
	server = server_
	return nil
}

func Close() error {
	mutex.Lock()
	defer mutex.Unlock()

	if !initialized {
		return nil
	}

	if !running {
		return fmt.Errorf("webserver is not running")
	}

	running = false

	if err := server.Shutdown(); err != nil {
		return fmt.Errorf("failed to close webserver: %w", err)
	}
	return nil
}

func Running() bool {
	mutex.RLock()
	defer mutex.RUnlock()

	return initialized && running
}
