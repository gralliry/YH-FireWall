package http

import (
	"YH-FireWall/core/config"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"log"
	"net/http"
	"sync"
)

var (
	server    *echo.Echo
	isRunning bool
	mutex     sync.RWMutex
)

func Start(cfg config.Web, handler Handler) error {
	mutex.Lock()
	defer mutex.Unlock()
	//
	if isRunning {
		return errors.New("web service already be started")
	}
	// 初始化 echo 实例
	e := echo.New()
	// 隐藏Banner
	e.HideBanner = true
	// 日志级别设置为OFF，关闭echo官方日志输出
	e.Logger.SetOutput(io.Discard)
	// 设置静态文件
	if cfg.StaticDir != "" {
		e.Static("/", cfg.StaticDir)
	}
	// 设置 BasicAuth 中间件
	if cfg.BasicAuthPassword != "" {
		e.Use(middleware.BasicAuth(func(usr, pwd string, c echo.Context) (bool, error) {
			return usr == cfg.BasicAuthUser && pwd == cfg.BasicAuthPassword, nil
		}))
	}
	// 必须放前面，提高api匹配优先级
	api := e.Group("/api")
	// 挂载接口
	mount(api, handler)
	// 启动服务器
	go start(e, cfg.Address)
	//
	server = e
	isRunning = true
	//
	return nil
}

func Close() error {
	mutex.Lock()
	defer mutex.Unlock()
	if !isRunning {
		return errors.New("web service not be stared")
	}
	isRunning = false
	if err := server.Close(); err != nil {
		return err
	}
	return nil
}

func IsRunning() bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return isRunning
}

func start(e *echo.Echo, addr string) {
	if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println(err)
	}
}
