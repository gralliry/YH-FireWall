package webserver

// @title           YH FireWall API
// @version         1.0
// @description     YH FireWall 防火墙管理接口
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@yh-firewall.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

import (
	_ "YH-FireWall/internal/model/conn"
	"YH-FireWall/internal/model/rule"

	"github.com/gofiber/fiber/v2"
)

// handlerPing  godoc
// @Summary     Ping
// @Description 健康检查
// @Tags        system
// @Success     200  {string}  string  "pong"
// @Router      /api/ping [get]
func handlerPing() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("pong")
	}
}

// handlerRuleList  godoc
// @Summary     获取所有规则
// @Description 返回防火墙规则列表
// @Tags        rule
// @Produce     json
// @Success     200  {array}   rule.Data
// @Router      /api/rule [get]
func handlerRuleList(handler Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cfgs := handler.ListRules()
		return c.JSON(cfgs)
	}
}

// handlerRuleCreate  godoc
// @Summary     添加规则
// @Description 添加一条防火墙规则
// @Tags        rule
// @Accept      json
// @Produce     plain
// @Param       option  body      rule.Option  true  "规则配置"
// @Success     200     {string}  string       "规则ID"
// @Failure     400     {string}  string       "错误信息"
// @Router      /api/rule [post]
func handlerRuleCreate(handler Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		option := new(rule.Option)
		if err := c.BodyParser(option); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		id, err := handler.CreateRule(option)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.SendString(id)
	}
}

// handlerRuleUpdate  godoc
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
func handlerRuleUpdate(handler Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
	}
}

// handlerRuleDelete  godoc
// @Summary     删除规则
// @Description 删除指定 ID 的防火墙规则
// @Tags        rule
// @Produce     plain
// @Param       id   path      string  true  "规则ID"
// @Success     200  {string}  string  "ok"
// @Failure     400  {string}  string  "错误信息"
// @Router      /api/rule/{id} [delete]
func handlerRuleDelete(handler Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("id is required")
		}
		if err := handler.DeleteRule(id); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

// handlerConfigGet  godoc
// @Summary     获取配置
// @Description 获取当前防火墙配置
// @Tags        config
// @Produce     plain
// @Success     200  {string}  string  "配置内容"
// @Failure     500  {string}  string  "错误信息"
// @Router      /api/config [get]
func handlerConfigGet(handler Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString(handler.GetConfig())
	}
}

// handlerConfigSet  godoc
// @Summary     更新配置
// @Description 更新防火墙配置
// @Tags        config
// @Accept      plain
// @Produce     plain
// @Param       data  body      string  true  "JSON 格式的配置内容"
// @Success     200   {string}  string  "ok"
// @Failure     500   {string}  string  "错误信息"
// @Router      /api/config [post]
func handlerConfigSet(handler Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		data := string(c.Body())
		if err := handler.SetConfig(data); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

// handlerConnectionClose  godoc
// @Summary     关闭连接
// @Description 强制关闭指定 ID 的网络连接
// @Tags        connection
// @Produce     plain
// @Param       id   path      string  true  "连接ID"
// @Success     200  {string}  string  "ok"
// @Failure     500  {string}  string  "错误信息"
// @Router      /api/connection/{id} [delete]
func handlerConnectionClose(handler Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("id is required")
		}
		if err := handler.CloseConnection(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

// handlerConnectionList  godoc
// @Summary     获取连接列表
// @Description 获取当前所有活跃的网络连接
// @Tags        connection
// @Produce     json
// @Success     200  {array}   conn.Info
// @Router      /api/connection [get]
func handlerConnectionList(handler Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		conns := handler.ListConnections()
		return c.JSON(conns)
	}
}
