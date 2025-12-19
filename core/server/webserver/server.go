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
	// 启动服务器
	go func() {
		if err := app.Listen(config.Address); err != nil {
			log.Println(err)
		}
	}()

	server = app
	isRunning = true
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
