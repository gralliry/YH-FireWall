package core

import (
	"YH-FireWall/core/config"
	"YH-FireWall/core/handler"
	"YH-FireWall/core/manager"
	"YH-FireWall/core/queue"
	"YH-FireWall/core/repo"
	"context"
	"errors"
	"os"
)

var (
	cfg *config.Config

	Context context.Context
	Cancel  context.CancelFunc

	hder *handler.Handler
)

func Start() (err error) {
	// 检测当前用户是否为 root 用户
	if os.Geteuid() != 0 {
		return errors.New("当前用户非 root 用户")
	}
	// 读取配置文件
	cfg = config.DefaultConfig()
	// 加载配置
	if err = repo.Start(); err != nil {
		return err
	}
	// 读取配置
	if err = repo.Load(cfg); err != nil {
		return err
	}
	// 初始化管理器
	if err = manager.Init(cfg); err != nil {
		return err
	}
	// 启动队列
	if err = queue.Start(Context, cfg.Queue); err != nil {
		return err
	}
	// 初始化接口
	hder = &handler.Handler{
		Context: Context,
		Cancel:  Cancel,
	}
	// 启动服务
	if err = hder.Start(); err != nil {
		return err
	}
	return nil
}

func Close() error {
	var errs []error
	// 关闭队列
	if err := queue.Close(); err != nil {
		errs = append(errs, err)
	}
	// 同步规则
	if err := manager.Sync(cfg); err != nil {
		errs = append(errs, err)
	}
	// 存储配置
	if err := repo.Store(cfg); err != nil {
		errs = append(errs, err)
	}
	// 关闭存储
	if err := repo.Close(); err != nil {
		errs = append(errs, err)
	}
	// 停止进程
	if err := hder.Stop(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
