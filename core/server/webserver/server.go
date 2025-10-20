package webserver

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

var (
	server    *fiber.App
	isRunning bool
)

type Config struct {
	Enable     bool   `json:"enable"`
	Address    string `json:"address"`
	Token      string `json:"auth_token"`
	StaticDir  string `json:"static_dir"`
	EnableCORS bool   `json:"enable_cors"`
}

func Start(handler Handler, config Config) error {
	// 初始化 Fiber 实例，并关闭默认日志
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true, // 隐藏启动信息
	})

	// 设置静态文件
	if config.StaticDir != "" {
		app.Static("/", config.StaticDir)
	}

	// 设置跨域中间件
	if config.EnableCORS {
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders: "Origin,Content-Type,Accept,Authorization",
		}))
	}

	// 设置 Token 验证中间件
	if config.Token != "" {
		app.Use(func(c *fiber.Ctx) error {
			auth := c.Get("Authorization")
			if auth != config.Token {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
			return c.Next()
		})
	}

	// API 分组
	api := app.Group("/api")

	// 挂载接口
	mount(api, handler)

	// 启动服务器
	go start(app, config.Address)

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

func start(app *fiber.App, addr string) {
	if err := app.Listen(addr); err != nil {
		log.Println(err)
	}
}
