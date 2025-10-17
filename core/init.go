package core

import (
	"YH-FireWall/core/config"
	"YH-FireWall/core/handler"
	"YH-FireWall/core/queue"
	"YH-FireWall/core/rtable"
	"YH-FireWall/core/server/cmdserver"
	"YH-FireWall/core/server/webserver"
	"context"
	"errors"
	"fmt"
	"os"
)

var (
	Context context.Context
	Cancel  context.CancelFunc

	hder *handler.Handler
)

func Start() (err error) {
	// 检测当前用户是否为 root 用户
	if os.Geteuid() != 0 {
		return errors.New("current user is not root")
	}
	// 读取配置
	if err = config.Init(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	// 初始化管理器
	if err = rtable.Load(cfg.RuleTable); err != nil {
		return fmt.Errorf("failed to load rule table: %w", err)
	}
	// 初始化队列
	if err = queue.Start(Context, cfg.QueueNo); err != nil {
		return fmt.Errorf("failed to start queue: %w", err)
	}
	// 初始化接口
	hder = &handler.Handler{
		Context: Context,
		Cancel:  Cancel,
	}
	// 启动服务
	if err = hder.Start(); err != nil {
		return fmt.Errorf("failed to start handler: %w", err)
	}
	// 设置cmd服务器
	if cfg.Cmd.Enable {
		if err = cmdserver.Start(hder, cfg.Cmd); err != nil {
			return fmt.Errorf("failed to start cmdserver: %w", err)
		}
	}
	// 设置web服务器
	if cfg.Web.Enable {
		if err = webserver.Start(hder, cfg.Web); err != nil {
			return fmt.Errorf("failed to start webserver: %w", err)
		}
	}
	return nil
}

func Close() error {
	var errs []error
	// 停止 cmdserver 监听
	if cmdserver.IsRunning() {
		if err := cmdserver.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	// 停止 webserver 监听
	if webserver.IsRunning() {
		if err := webserver.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	// 停止进程
	if err := hder.Stop(); err != nil {
		errs = append(errs, err)
	}
	// 关闭队列
	if err := queue.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := rtable.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := config.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
