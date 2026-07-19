//go:generate swag init -g handler.go -o docs --parseDependency --parseDepth 2

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
	_ "YH-FireWall/internal/model/itf"
	"YH-FireWall/internal/model/rule"
	_ "YH-FireWall/internal/server/webserver/docs"

	"github.com/gofiber/fiber/v3"
)

// handlePing  godoc
// @Summary     Ping
// @Description 健康检查
// @Tags        system
// @Success     200  {string}  string  "pong"
// @Router      /ping [get]
func handlePing() fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.SendString("pong")
	}
}

// handleRuleList  godoc
// @Summary     获取所有规则
// @Description 返回防火墙规则列表
// @Tags        rule
// @Produce     json
// @Success     200  {array}   rule.Data
// @Router      /rule [get]
func handleRuleList(handler Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		cfgs := handler.ListRules()
		return c.JSON(cfgs)
	}
}

// handleRuleCreate  godoc
// @Summary     添加规则
// @Description 添加一条防火墙规则
// @Tags        rule
// @Accept      json
// @Produce     plain
// @Param       option  body      rule.Option  true  "规则配置"
// @Success     200     {string}  string       "规则ID"
// @Failure     400     {string}  string       "错误信息"
// @Router      /rule [post]
func handleRuleCreate(handler Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		option := new(rule.Option)
		if err := c.Bind().Body(option); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		id, err := handler.CreateRule(option)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.SendString(id)
	}
}

// handleRuleUpdate  godoc
// @Summary     更新规则
// @Description 更新指定 ID 的防火墙规则
// @Tags        rule
// @Accept      json
// @Produce     plain
// @Param       id      path      string      true  "规则ID"
// @Param       option  body      rule.Option true  "规则配置"
// @Success     200     {string}  string      "ok"
// @Failure     400     {string}  string      "错误信息"
// @Router      /rule/{id} [put]
func handleRuleUpdate(handler Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).SendString("id is required")
		}
		option := new(rule.Option)
		if err := c.Bind().Body(option); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		if err := handler.UpdateRule(id, option); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

// handleRuleDelete  godoc
// @Summary     删除规则
// @Description 删除指定 ID 的防火墙规则
// @Tags        rule
// @Produce     plain
// @Param       id   path      string  true  "规则ID"
// @Success     200  {string}  string  "ok"
// @Failure     400  {string}  string  "错误信息"
// @Router      /rule/{id} [delete]
func handleRuleDelete(handler Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
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

// handleConfigGet  godoc
// @Summary     获取配置
// @Description 获取当前防火墙配置
// @Tags        config
// @Produce     plain
// @Success     200  {string}  string  "配置内容"
// @Failure     500  {string}  string  "错误信息"
// @Router      /config [get]
func handleConfigGet(handler Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.SendString(handler.GetConfig())
	}
}

// handleConfigSet  godoc
// @Summary     更新配置
// @Description 更新防火墙配置
// @Tags        config
// @Accept      plain
// @Produce     plain
// @Param       data  body      string  true  "JSON 格式的配置内容"
// @Success     200   {string}  string  "ok"
// @Failure     500   {string}  string  "错误信息"
// @Router      /config [post]
func handleConfigSet(handler Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		data := string(c.Body())
		if err := handler.SetConfig(data); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	}
}

// handleConnectionClose  godoc
// @Summary     关闭连接
// @Description 强制关闭指定 ID 的网络连接
// @Tags        connection
// @Produce     plain
// @Param       id   path      string  true  "连接ID"
// @Success     200  {string}  string  "ok"
// @Failure     500  {string}  string  "错误信息"
// @Router      /connection/{id} [delete]
func handleConnectionClose(handler Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
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

// handleConnectionList  godoc
// @Summary     获取连接列表
// @Description 获取当前所有活跃的网络连接
// @Tags        connection
// @Produce     json
// @Success     200  {array}   conn.Info
// @Router      /connection [get]
func handleConnectionList(handler Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		conns, err := handler.ListConnections()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(conns)
	}
}

// handleInterfaceList  godoc
// @Summary     获取网卡列表
// @Description 获取当前系统所有网络接口信息
// @Tags        system
// @Produce     json
// @Success     200  {array}  itf.Itf
// @Router      /interface [get]
func handleInterfaceList(handler Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		ifaces, err := handler.ListInterfaces()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(ifaces)
	}
}

// handleProtocolList  godoc
// @Summary     获取协议列表
// @Description 获取防火墙支持的所有 IP 协议名称
// @Tags        system
// @Produce     json
// @Success     200  {array}  string
// @Router      /protocol [get]
func handleProtocolList(handler Handler) fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.JSON(handler.ListProtocols())
	}
}
