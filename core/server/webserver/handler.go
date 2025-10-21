package webserver

import (
	"YH-FireWall/core/connection"
	"YH-FireWall/core/process"
	"YH-FireWall/core/rule"
	"YH-FireWall/core/system"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Handler interface {
	AppendRule(ro *rule.Option) (string, error)
	UpdateRule(id string, ro *rule.Option) error
	DeleteRule(id string) error
	GetRules() []rule.Config
	EnableRule(id string, enable bool) bool
	GetConfig() (string, error)
	SetConfig(raw string) error

	GetConnections() []connection.Config
	CloseConnection(id string) error

	GetProcesses() ([]process.Process, error)
	CloseProcess(pid int32, fd uint32) error
	GetInterfaces() ([]system.Interface, error)
}

func mount(api fiber.Router, handler Handler) {
	// ping
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// 获取规则
	api.Get("/rule", func(c *fiber.Ctx) error {
		cfgs := handler.GetRules()
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
		data, err := handler.GetConfig()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendString(data)
	})

	// 设置配置
	api.Post("/config", func(c *fiber.Ctx) error {
		data := c.Body()
		if err := handler.SetConfig(string(data)); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// 获取连接
	api.Get("/connection", func(c *fiber.Ctx) error {
		conns := handler.GetConnections()
		return c.JSON(conns)
	})

	//
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

	// 获取进程
	api.Get("/process", func(c *fiber.Ctx) error {
		conns, err := handler.GetProcesses()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(conns)
	})

	// 关闭进程连接
	api.Delete("/process/:pid/:fd", func(c *fiber.Ctx) error {
		pid := c.Params("pid")
		if pid == "" {
			return c.Status(fiber.StatusBadRequest).SendString("pid is required")
		}
		pidInt32, err := strconv.ParseInt(pid, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid pid")
		}

		fd := c.Params("fd")
		if fd == "" {
			return c.Status(fiber.StatusBadRequest).SendString("fd is required")
		}
		fdUint32, err := strconv.ParseUint(fd, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid fd")
		}

		if err = handler.CloseProcess(int32(pidInt32), uint32(fdUint32)); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// 获取网络接口
	api.Get("/interface", func(c *fiber.Ctx) error {
		interfaces, err := handler.GetInterfaces()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(interfaces)
	})
}
