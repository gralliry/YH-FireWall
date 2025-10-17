package webserver

import (
	"YH-FireWall/core/connection"
	"YH-FireWall/core/rule"
	"YH-FireWall/core/system"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
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
	GetConnections() ([]connection.Connection, error)
	CloseConnection(pid int32, fd uint32) error
	GetInterfaces() ([]system.Interface, error)
}

//// 必须放前面，提高api匹配优先级
//api := e.Group("/api")

func mount(api *echo.Group, handler Handler) {
	//api.GET("/ping", ping)
	api.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	// 规则增删改
	//api.GET("/rule", getRules)
	api.GET("/rule", func(c echo.Context) error {
		cfgs := handler.GetRules()
		return c.JSON(http.StatusOK, cfgs)
	})
	//api.POST("/rule", appendRule)
	api.POST("/rule", func(c echo.Context) error {
		option := new(rule.Option)
		if err := c.Bind(option); err != nil {
			return err
		}
		id, err := handler.AppendRule(option)
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, id)
	})
	//api.PUT("/rule/:id", updateRule)
	api.PUT("/rule/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.String(http.StatusBadRequest, "id is required")
		}
		option := new(rule.Option)
		if err := c.Bind(option); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		if err := handler.UpdateRule(id, option); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})
	//api.DELETE("/rule/:id", deleteRule)
	api.DELETE("/rule/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.String(http.StatusBadRequest, "id is required")
		}
		if err := handler.DeleteRule(id); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})
	// config get/set
	api.GET("/config", func(c echo.Context) error {
		data, err := handler.GetConfig()
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.String(http.StatusOK, data)
	})
	api.POST("/config", func(c echo.Context) error {
		data, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		if err = handler.SetConfig(string(data)); err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	})
	api.GET("/connection", func(c echo.Context) error {
		conns, err := handler.GetConnections()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, conns)
	})
	api.DELETE("/connection/:pid/:fd", func(c echo.Context) error {
		pid := c.Param("pid")
		if pid == "" {
			return c.String(http.StatusBadRequest, "pid is required")
		}
		pidInt32, err := strconv.ParseInt(pid, 10, 32)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid pid")
		}
		fd := c.Param("fd")
		if fd == "" {
			return c.String(http.StatusBadRequest, "fd is required")
		}
		fdUint32, err := strconv.ParseUint(fd, 10, 32)
		if err != nil {
			return c.String(http.StatusBadRequest, "invalid fd")
		}
		if err = handler.CloseConnection(int32(pidInt32), uint32(fdUint32)); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})
	api.GET("/interface", func(c echo.Context) error {
		interfaces, err := handler.GetInterfaces()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, interfaces)
	})
}
