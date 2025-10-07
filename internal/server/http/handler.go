package http

import (
	"YH-FireWall/internal/core"
	"YH-FireWall/internal/core/rule"
	"github.com/labstack/echo/v4"
	"net/http"
)

func getRules(c echo.Context) error {
	cfgs := core.GetRules()
	return c.JSON(http.StatusOK, cfgs)
}

func appendRule(c echo.Context) error {
	option := new(rule.Option)
	if err := c.Bind(option); err != nil {
		return err
	}
	if err := core.AppendRule(option); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func updateRule(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	option := new(rule.Option)
	if err := c.Bind(option); err != nil {
		return err
	}
	if err := core.UpdateRule(id, option); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func deleteRule(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := core.DeleteRule(id); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func enableRule(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	if core.EnableRule(id, true) {
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusBadRequest)
	}
}

func disableRule(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.NoContent(http.StatusBadRequest)
	}
	if core.EnableRule(id, false) {
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusBadRequest)
	}
}

func enableGroup(c echo.Context) error {
	group := c.QueryParam("g")
	if core.EnableGroup(group, true) {
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusBadRequest)
	}
}

func disableGroup(c echo.Context) error {
	group := c.QueryParam("g")
	if core.EnableGroup(group, false) {
		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusBadRequest)
	}
}
