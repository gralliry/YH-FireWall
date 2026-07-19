package webserver

import (
	"YH-FireWall/internal/model/conn"
	"YH-FireWall/internal/model/itf"
	"YH-FireWall/internal/model/rule"
	"YH-FireWall/ui"
	"fmt"
	"io/fs"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/basicauth"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/swaggo/swag"

	_ "YH-FireWall/internal/server/webserver/docs"
)

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
	ListConnections() ([]*conn.Info, error)
	//
	ListInterfaces() ([]itf.Itf, error)
	ListProtocols() []string
}

func newApp(config Config, handler Handler) (*fiber.App, error) {
	app := fiber.New()

	// 设置跨域中间件
	if config.EnableCORS {
		app.Use(cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
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
	app.Get("/docs/swagger.json", func(c fiber.Ctx) error {
		doc, err := swag.ReadDoc()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendString(doc)
	})
	app.Get("/docs", serveSwaggerUI)

	// 前端文件
	if config.StaticDir != "" {
		app.Use(static.New(config.StaticDir, static.Config{
			Browse: false,
		}))
	} else {
		subFS, err := fs.Sub(ui.FS, "dist")
		if err != nil {
			return nil, fmt.Errorf("frontend dist not found: %w (run 'npm run build' first)", err)
		}
		app.Use(static.New("", static.Config{
			FS:     subFS,
			Browse: false,
		}))
	}

	return app, nil
}

const swaggerHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>YH FireWall API</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.17.14/swagger-ui.min.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.17.14/swagger-ui-bundle.min.js"></script>
  <script>
    SwaggerUIBundle({
      url: '/docs/swagger.json',
      dom_id: '#swagger-ui',
      deepLinking: true,
    })
  </script>
</body>
</html>`

func serveSwaggerUI(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(swaggerHTML)
}
