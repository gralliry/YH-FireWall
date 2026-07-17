package webserver

import (
	"YH-FireWall/internal/model/conn"
	"YH-FireWall/internal/model/itf"
	"YH-FireWall/internal/model/rule"
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/swagger"
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
	GetConfig() (string, error)
	SetConfig(data string) error
	//
	CloseConnection(id string) error
	ListConnections() []*conn.Info
	//
	ListInterfaces() []*itf.Itf
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
	// @Summary     Ping
	// @Description 健康检查
	// @Tags        system
	// @Success     200  {string}  string  "pong"
	// @Router      /api/ping [get]
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// 获取规则
	// @Summary     获取所有规则
	// @Description 返回防火墙规则列表
	// @Tags        rule
	// @Produce     json
	// @Success     200  {array}   rule.Data
	// @Router      /api/rule [get]
	api.Get("/rule", func(c *fiber.Ctx) error {
		cfgs := handler.ListRules()
		return c.JSON(cfgs)
	})

	// 添加规则
	// @Summary     添加规则
	// @Description 添加一条防火墙规则
	// @Tags        rule
	// @Accept      json
	// @Produce     plain
	// @Param       option  body      rule.Option  true  "规则配置"
	// @Success     200     {string}  string       "规则ID"
	// @Failure     400     {string}  string       "错误信息"
	// @Router      /api/rule [post]
	api.Post("/rule", func(c *fiber.Ctx) error {
		option := new(rule.Option)
		if err := c.BodyParser(option); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		id, err := handler.CreateRule(option)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.SendString(id)
	})

	// 更新规则
	// @Summary     更新规则
	// @Description 更新指定 ID 的防火墙规则
	// @Tags        rule
	// @Accept      json
	// @Produce     plain
	// @Param       id      path      string      true  "规则ID"
	// @Param       option  body      rule.Option true  "规则配置"
	// @Success     200     {string}  string      "ok"
	// @Failure     400     {string}  string      "错误信息"
	// @Router      /api/rule/{id} [put]
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
	// @Summary     删除规则
	// @Description 删除指定 ID 的防火墙规则
	// @Tags        rule
	// @Produce     plain
	// @Param       id   path      string  true  "规则ID"
	// @Success     200  {string}  string  "ok"
	// @Failure     400  {string}  string  "错误信息"
	// @Router      /api/rule/{id} [delete]
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
	// @Summary     获取配置
	// @Description 获取当前防火墙配置
	// @Tags        config
	// @Produce     plain
	// @Success     200  {string}  string  "配置内容"
	// @Failure     500  {string}  string  "错误信息"
	// @Router      /api/config [get]
	api.Get("/config", func(c *fiber.Ctx) error {
		data, err := handler.GetConfig()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendString(data)
	})

	// 设置配置
	// @Summary     更新配置
	// @Description 更新防火墙配置
	// @Tags        config
	// @Accept      plain
	// @Produce     plain
	// @Param       data  body      string  true  "JSON 格式的配置内容"
	// @Success     200   {string}  string  "ok"
	// @Failure     500   {string}  string  "错误信息"
	// @Router      /api/config [post]
	api.Post("/config", func(c *fiber.Ctx) error {
		buf := c.Body()
		data := string(buf)
		if err := handler.SetConfig(data); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// 关闭连接
	// @Summary     关闭连接
	// @Description 强制关闭指定 ID 的网络连接
	// @Tags        connection
	// @Produce     plain
	// @Param       id   path      string  true  "连接ID"
	// @Success     200  {string}  string  "ok"
	// @Failure     500  {string}  string  "错误信息"
	// @Router      /api/connection/{id} [delete]
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
	// @Summary     获取连接列表
	// @Description 获取当前所有活跃的网络连接
	// @Tags        connection
	// @Produce     json
	// @Success     200  {array}   conn.Info
	// @Router      /api/connection [get]
	api.Get("/connection", func(c *fiber.Ctx) error {
		conns := handler.ListConnections()
		return c.JSON(conns)
	})

	// 获取网络接口
	// @Summary     获取网卡列表
	// @Description 获取系统网络接口信息
	// @Tags        system
	// @Produce     json
	// @Success     200  {array}   itf.Itf
	// @Router      /api/interface [get]
	api.Get("/interface", func(c *fiber.Ctx) error {
		interfaces := handler.ListInterfaces()
		return c.JSON(interfaces)
	})

	// Swagger 文档
	// @Summary     Swagger UI
	// @Description Swagger API 文档页面
	// @Router      /swagger/* [get]
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
