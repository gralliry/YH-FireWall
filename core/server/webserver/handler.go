package webserver

import (
	"YH-FireWall/core/rule"
	"YH-FireWall/core/system"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

type Handler interface {
	AppendRule(ro *rule.Option) error
	UpdateRule(id string, ro *rule.Option) error
	DeleteRule(id string) error
	GetRules() []rule.Config
	EnableRule(id string, enable bool) bool
	GetConfig() (string, error)
	SetConfig(raw string) error
	GetConnections() ([]system.Connection, error)
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
		if err := handler.AppendRule(option); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})
	//api.PUT("/rule/:id", updateRule)
	api.PUT("/rule/:id", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.NoContent(http.StatusBadRequest)
		}
		option := new(rule.Option)
		if err := c.Bind(option); err != nil {
			return err
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
			return c.NoContent(http.StatusBadRequest)
		}
		if err := handler.DeleteRule(id); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})
	//// 规则启用禁用
	//api.PUT("/rule/:id/enable", enableRule)
	api.PUT("/rule/:id/enable", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.NoContent(http.StatusBadRequest)
		}
		if !handler.EnableRule(id, true) {
			return c.NoContent(http.StatusBadRequest)
		}
		return c.NoContent(http.StatusOK)
	})
	//api.PUT("/rule/:id/disable", disableRule)
	api.PUT("/rule/:id/disable", func(c echo.Context) error {
		id := c.Param("id")
		if id == "" {
			return c.NoContent(http.StatusBadRequest)
		}
		if !handler.EnableRule(id, false) {
			return c.NoContent(http.StatusBadRequest)
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
	//
	api.GET("/connection", func(c echo.Context) error {
		conns, err := handler.GetConnections()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, conns)
	})
	//
	api.GET("/interface", func(c echo.Context) error {
		interfaces, err := handler.GetInterfaces()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, interfaces)
	})
}
