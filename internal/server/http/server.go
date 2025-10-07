package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"log"
	"net/http"
	"sync"
)

var (
	server *echo.Echo
	mutex  sync.RWMutex
)

func Start(addr, user, pswd string) error {
	// 初始化 echo 实例
	e := echo.New()
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
	// 设置静态文件
	e.Static("/", "front/dist")
	// 设置 BasicAuth 中间件
	e.Use(middleware.BasicAuth(func(usr, pwd string, c echo.Context) (bool, error) {
		return usr == user && pwd == pswd, nil
	}))
	// 必须放前面，提高api匹配优先级
	api := e.Group("/api")
	api.GET("/ping", ping)
	// 组内规则增删改
	api.GET("/rule", getRules)
	api.POST("/rule", appendRule)
	api.PUT("/rule/:id", updateRule)
	api.DELETE("/rule/:id", deleteRule)
	// 规则启用禁用
	api.PUT("/rule/:id/enable", enableRule)
	api.PUT("/rule/:id/disable", disableRule)
	// 组启用禁用
	api.PUT("/group/enable", enableGroup)
	api.PUT("/group/disable", disableGroup)
	// 启动服务器
	go func(addr string) {
		mutex.Lock()
		server = e
		mutex.Unlock()
		if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Web接口服务运行失败: %v", err)
		}
		mutex.Lock()
		server = nil
		mutex.Unlock()
	}(addr)
	return nil
}

func Close() error {
	mutex.Lock()
	defer mutex.Unlock()
	if server == nil {
		return nil
	}
	if err := server.Close(); err != nil {
		return err
	}
	return nil
}

func IsRunning() bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return server != nil
}
