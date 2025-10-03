package web

import (
	"YH-FireWall/internal/mapper"
	"YH-FireWall/internal/rule"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func handleGlobalError(err error, c echo.Context) {
	_ = c.JSON(http.StatusInternalServerError, err.Error())
}

func ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func getAllRules(c echo.Context) error {
	cfgs, err := mapper.GetAllRules()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cfgs)
}

func updateRule(c echo.Context) error {
	cfg := new(rule.Config)
	if err := c.Bind(cfg); err != nil {
		return err
	}
	if err := mapper.UpdateRule(*cfg); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "success")
}

func appendRule(c echo.Context) error {
	cfg := new(rule.Config)
	if err := c.Bind(cfg); err != nil {
		return err
	}
	if err := mapper.AppendRule(*cfg); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "success")
}

func deleteRule(c echo.Context) error {
	idStr := c.QueryParam("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	if err = mapper.DeleteRule(uint32(id)); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "success")
}
