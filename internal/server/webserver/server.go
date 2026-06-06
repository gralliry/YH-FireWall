package webserver

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
)

var (
	server    *fiber.App
	isRunning bool
)

type Config struct {
	Enable       bool   `json:"enable"`
	Address      string `json:"address"`
	AuthUsername string `json:"auth_username"`
	AuthPassword string `json:"auth_password"`
	EnableCORS   bool   `json:"enable_cors"`
}

func Start(handler Handler, config Config) error {
	app := newServer(config, handler)
	// 在 goroutine 中启动，确认监听成功后设置标记
	ready := make(chan struct{})
	go func() {
		isRunning = true
		close(ready)
		if err := app.Listen(config.Address); err != nil {
			log.Println(err)
			isRunning = false
		}
	}()
	<-ready
	server = app
	return nil
}

func Close() error {
	if !isRunning {
		return fmt.Errorf("webserver is not running")
	}
	isRunning = false
	if err := server.Shutdown(); err != nil {
		return fmt.Errorf("failed to close webserver: %w", err)
	}
	return nil
}

func IsRunning() bool {
	return isRunning
}
