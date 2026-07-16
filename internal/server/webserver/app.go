package webserver

import (
	"YH-FireWall/internal/model/conn"
	"YH-FireWall/internal/model/itf"
	"YH-FireWall/internal/model/rule"
	"embed"
	"net/http"
	"unsafe"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed static/*
var staticFS embed.FS

type Handler interface {
	Version() string
	//
	AppendRule(ro *rule.Option) (string, error)
	UpdateRule(id string, ro *rule.Option) error
	DeleteRule(id string) error
	SearchRule(id string) *rule.Option
	ListRules() []*rule.Option
	EnableRule(id string, enable bool) bool
	//
	GetConfig() string
	SetConfig(data string) error
	//
	CloseConnection(id string) error
	ListConnections() []conn.Info
	//
	ListInterfaces() []itf.Itf
	//
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

	// ping
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// 获取规则
	api.Get("/rule", func(c *fiber.Ctx) error {
		cfgs := handler.ListRules()
		return c.JSON(cfgs)
	})

	// 添加规则
	api.Post("/rule", func(c *fiber.Ctx) error {
		option := new(rule.Option)
		if err := c.BodyParser(option); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		id, err := handler.AppendRule(option)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.SendString(id)
	})

	// 更新规则
	api.Put("/rule/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("id is required")
		}
		option := new(rule.Option)
		if err := c.BodyParser(option); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		if err := handler.UpdateRule(id, option); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// 删除规则
	api.Delete("/rule/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("id is required")
		}
		if err := handler.DeleteRule(id); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// 获取配置
	api.Get("/config", func(c *fiber.Ctx) error {
		data := handler.GetConfig()
		return c.SendString(data)
	})

	// 设置配置
	api.Post("/config", func(c *fiber.Ctx) error {
		buf := c.Body()
		data := unsafe.String(unsafe.SliceData(buf), len(buf))
		if err := handler.SetConfig(data); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// 关闭连接
	api.Delete("/connection/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("id is required")
		}
		if err := handler.CloseConnection(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// 获取连接
	api.Get("/connection", func(c *fiber.Ctx) error {
		conns := handler.ListConnections()
		return c.JSON(conns)
	})

	// 获取网络接口
	api.Get("/interface", func(c *fiber.Ctx) error {
		interfaces := handler.ListInterfaces()
		return c.JSON(interfaces)
	})

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
