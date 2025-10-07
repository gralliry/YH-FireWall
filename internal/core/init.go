package core

import (
	"YH-FireWall/internal/core/manager"
	"YH-FireWall/internal/core/queue"
	"YH-FireWall/internal/core/repo"
	"YH-FireWall/internal/core/rule"
	"context"
	"errors"
	"os"
	"sync"
	"time"
)

const (
	ConfigPath = "/etc/yfw/config.json"
	Version    = "v1.0.0"
)

var (
	cfg    *repo.Config
	ctx    context.Context
	cancel context.CancelFunc

	isRunning bool
	mutex     sync.Mutex
)

func Start(parent context.Context) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	// 检测服务是否已启动
	if isRunning {
		return errors.New("服务已启动")
	}
	// 检测当前用户是否为 root 用户
	if os.Geteuid() != 0 {
		return errors.New("当前用户非 root 用户")
	}
	// 读取配置文件
	cfg, err = repo.Load(ConfigPath)
	if err != nil {
		return err
	}
	// 加载组配置
	var rr *rule.Rule
	for _, rc := range cfg.Rules {
		// 解析规则
		rr, err = rule.Parse(rc)
		if err != nil {
			continue
		}
		// 检测ID是否被使用
		err = manager.RegisterRule(rr)
		if err != nil {
			continue
		}
	}
	//
	ctx, cancel = context.WithCancel(parent)
	// 启动队列
	if err = queue.Start(ctx, cfg.QueueNum, cfg.DefaultAccept); err != nil {
		return err
	}
	isRunning = true
	return nil
}

func Done() <-chan struct{} {
	mutex.Lock()
	defer mutex.Unlock()
	return ctx.Done()
}

func Close() error {
	mutex.Lock()
	defer mutex.Unlock()
	if !isRunning {
		return nil
	}
	cancel()
	var errs []error
	// 关闭队列
	if err := queue.Close(); err != nil {
		errs = append(errs, err)
	}
	// 更新时间
	cfg.UpdateDate = time.Now().Format("2006-01-02 15:04:05")
	// 存储配置文件
	cfg.Rules = manager.GetRules()
	//
	if err := cfg.Store(); err != nil {
		errs = append(errs, err)
	}
	if err := cfg.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
