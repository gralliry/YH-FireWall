package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"log"
	"net/http"
)

var e *echo.Echo

func Start(address, username, password string) error {
	// 初始化 echo 实例
	e = echo.New()
	// 隐藏Banner
	e.HideBanner = true
	// 日志级别设置为OFF，关闭echo官方日志输出
	e.Logger.SetOutput(io.Discard)
	// 设置跨域 // 使用 CORS 中间件
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // 允许的域名，可以写具体域名
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))
	// 设置 BasicAuth 中间件
	e.Use(middleware.BasicAuth(func(user, pwd string, c echo.Context) (bool, error) {
		// 这里做认证逻辑，例如固定用户名密码
		if user == username && pwd == password {
			return true, nil
		} else {
			return false, nil
		}
	}))
	// 设置全局错误处理
	e.HTTPErrorHandler = handleGlobalError
	// 必须放前面，提高api匹配优先级
	api := e.Group("/server")
	// 注册api
	api.GET("/ping", ping)
	//
	api.GET("/rule", getRules)
	api.POST("/rule", appendRule)
	api.PUT("/rule", updateRule)
	api.DELETE("/rule", deleteRule)
	// 启动服务器
	go func() {
		if err := e.Start(address); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Web接口服务启动失败: %v", err)
		}
	}()
	return nil
}

func Close() error {
	if e == nil {
		return nil
	}
	return e.Close()
}
