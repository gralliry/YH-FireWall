package webserver

import (
	"YH-FireWall/internal/model/conn"
	"YH-FireWall/internal/model/rule"
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/swagger"

	_ "YH-FireWall/internal/server/webserver/docs"
)

//go:embed static/*
var staticFS embed.FS

type Handler interface {
	Version() string
	//
	CreateRule(ro *rule.Option) (string, error)
	UpdateRule(id string, ro *rule.Option) error
	DeleteRule(id string) error

	SearchRule(id string) *rule.Data
	ListRules() []*rule.Data

	EnableRule(id string, enable bool) error
	//
	GetConfig() string
	SetConfig(data string) error
	//
	CloseConnection(id string) error
	ListConnections() []*conn.Info
	//
	ListInterfaces() []string
	ListProtocols() []string
}

func newApp(config Config, handler Handler) *fiber.App {
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
	if config.AuthUsername != "" && config.AuthPassword != "" {
		app.Use(basicauth.New(basicauth.Config{
			Users: map[string]string{config.AuthUsername: config.AuthPassword},
			Realm: "Firewall Web Login",
		}))
	}

	// API 分组
	api := app.Group("/api")

	// 路由
	api.Get("/ping", handlePing())
	api.Get("/rule", handleRuleList(handler))
	api.Post("/rule", handleRuleCreate(handler))
	api.Put("/rule/:id", handleRuleUpdate(handler))
	api.Delete("/rule/:id", handleRuleDelete(handler))
	api.Get("/config", handleConfigGet(handler))
	api.Post("/config", handleConfigSet(handler))
	api.Delete("/connection/:id", handleConnectionClose(handler))
	api.Get("/connection", handleConnectionList(handler))
	api.Get("/interface", handleInterfaceList(handler))
	api.Get("/protocol", handleProtocolList(handler))

	// Swagger 文档
	app.Get("/swagger/*", swagger.HandlerDefault)

	// 前端文件
	var root http.FileSystem
	if config.StaticDir != "" {
		root = http.Dir(config.StaticDir)
	} else {
		root = http.FS(staticFS)
	}
	app.All("/*", filesystem.New(filesystem.Config{
		Root:       root,
		PathPrefix: "static",
		Browse:     false,
		Index:      "index.html",
	}))

	return app
}
