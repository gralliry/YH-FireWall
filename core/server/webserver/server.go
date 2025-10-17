package webserver

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"log"
	"net/http"
)

var (
	server    *echo.Echo
	isRunning bool
)

type Config struct {
	Enable     bool   `json:"enable"`
	Address    string `json:"address"`
	Token      string `json:"auth_token"`
	StaticDir  string `json:"static_dir"`
	EnableCORS bool   `json:"enable_cors"`
}

func Start(handler Handler, config Config) error {
	// 初始化 echo 实例
	e := echo.New()
	// 隐藏Banner
	e.HideBanner = true
	// 日志级别设置为OFF，关闭echo官方日志输出
	e.Logger.SetOutput(io.Discard)
	// 设置静态文件
	if config.StaticDir != "" {
		e.Static("/", config.StaticDir)
	}
	// 设置跨域中间件
	if config.EnableCORS {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     []string{"*"}, // 允许所有来源
			AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
			AllowCredentials: true,
		}))
	}
	// 设置 BasicAuth 中间件
	if config.Token != "" {
		e.Use(middleware.KeyAuth(func(auth string, c echo.Context) (bool, error) {
			return auth == config.Token, nil
		}))
	}
	// 必须放前面，提高api匹配优先级
	api := e.Group("/api")
	// 挂载接口
	mount(api, handler)
	// 启动服务器
	go start(e, config.Address)
	//
	server = e
	isRunning = true
	//
	return nil
}

func Close() error {
	isRunning = false
	if err := server.Close(); err != nil {
		return err
	}
	return nil
}

func IsRunning() bool {
	return isRunning
}

func start(e *echo.Echo, addr string) {
	if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println(err)
	}
}
