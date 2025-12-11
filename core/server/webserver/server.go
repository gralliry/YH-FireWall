package webserver

import (
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"log"
	"net/http"
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

//go:embed static/*
var static embed.FS

func Start(handler Handler, config Config) error {
	// 初始化 Fiber 实例，并关闭默认日志
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true, // 隐藏启动信息
	})

	// 设置跨域中间件
	if config.EnableCORS {
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders: "Origin,Content-Type,Accept,Authorization",
		}))
	}

	// 设置验证中间件
	if config.AuthPassword != "" {
		app.Use(basicauth.New(basicauth.Config{
			Users: map[string]string{config.AuthUsername: config.AuthPassword},
			Realm: "Firewall Web Login",
		}))
	}

	// API 分组
	api := app.Group("/api")
	// 挂载接口
	mount(api, handler)

	// 方法1: 使用 filesystem 中间件
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(static),
		PathPrefix: "static",
		Browse:     true, // 允许目录浏览
		Index:      "index.html",
	}))

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
