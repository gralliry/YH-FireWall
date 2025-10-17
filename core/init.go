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
	rtable.Cfg = cfg.RuleTable
	if err = rtable.Load(); err != nil {
		return fmt.Errorf("failed to load rule table: %w", err)
	}
	// 初始化队列
	queue.NfQueueNo = cfg.QueueNo
	if err = queue.Start(Context); err != nil {
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
		cmdserver.Cfg = cfg.Cmd
		if err = cmdserver.Start(hder); err != nil {
			return fmt.Errorf("failed to start cmdserver: %w", err)
		}
	}
	// 设置web服务器
	if cfg.Web.Enable {
		webserver.Cfg = cfg.Web
		if err = webserver.Start(hder); err != nil {
			return fmt.Errorf("failed to start webserver: %w", err)
		}
	}
	return nil
}

func Close() error {
	var errs []error
	// 停止 cmdserver 监听
	if cmdserver.IsRunning() {
		errs = append(errs, cmdserver.Close())
	}
	// 停止 webserver 监听
	if webserver.IsRunning() {
		errs = append(errs, webserver.Close())
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
