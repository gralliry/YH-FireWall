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
	"gopkg.in/yaml.v3"
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
	data, err := config.Read()
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	cfg := config.Default()
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	// 初始化管理器
	rtable.RuleTablePath = cfg.RuleTablePath
	rtable.DefaultAccept = cfg.RuleTableDefaultAccept
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
	if cfg.CmdEnable {
		cmdserver.SocketPath = cfg.CmdSocketPath
		if err = cmdserver.Start(hder); err != nil {
			return fmt.Errorf("failed to start cmdserver: %w", err)
		}
	}
	// 设置web服务器
	if cfg.WebEnable {
		webserver.Address = cfg.WebAddress
		webserver.BasicAuthUser = cfg.WebBasicAuthUser
		webserver.BasicAuthPassword = cfg.WebBasicAuthPassword
		webserver.StaticDir = cfg.WebStaticDir
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
