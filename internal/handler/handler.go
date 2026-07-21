package handler

import (
	"YH-FireWall/internal/config"
	"YH-FireWall/internal/constant/itfdev"
	"YH-FireWall/internal/ctable"
	"YH-FireWall/internal/queue"
	"YH-FireWall/internal/rtable"
	"YH-FireWall/internal/server/cmd"
	"YH-FireWall/internal/server/http"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sync"
)

type Handler struct {
	configer  *config.Manager
	ruler     *rtable.Manager
	connecter *ctable.Manager
	queuer    *queue.NFQ
	//
	cmder *cmdserver.Server
	weber *webserver.Server

	closeOnce sync.Once
}

func New(configPath string) (h *Handler, err error) {

	// 回调清理
	var cleanups []func()
	defer func() {
		if err == nil {
			return
		}
		for _, cleanup := range slices.Backward(cleanups) {
			cleanup()
		}
	}()

	// 启动日志
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.TimeKey {
				return a
			}
			t := a.Value.Time()
			a.Value = slog.StringValue(t.Format("2006-01-02 15:04:05.000"))
			return a
		},
	}))

	// 加载接口
	if err = itfdev.Load(); err != nil {
		return nil, fmt.Errorf("failed to load interfaces: %w", err)
	}

	// 读取配置管理器
	configer, err := config.New(configPath, logger.With("module", "config"))
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	cleanups = append(cleanups, func() { configer.Close() })

	// 加载配置
	cfg := h.configer.Load()

	// 初始化规则管理器
	ruler, err := rtable.New(cfg.Rule, logger.With("module", "rule"))
	if err != nil {
		return nil, fmt.Errorf("failed to load rule table: %w", err)
	}
	cleanups = append(cleanups, func() { ruler.Close() })

	// 初始化连接表
	connecter := ctable.New()

	// 初始化队列
	queuer, err := queue.New(cfg.Queue, h)
	if err != nil {
		return nil, fmt.Errorf("failed to start queue: %w", err)
	}
	cleanups = append(cleanups, func() { queuer.Close() })

	// 设置cmd服务器
	cmder, err := cmdserver.New(cfg.Cmd, h)
	if err != nil {
		return nil, fmt.Errorf("failed to start cmdserver: %w", err)
	}
	cleanups = append(cleanups, func() { cmder.Close() })

	// 设置web服务器
	weber, err := webserver.New(cfg.Web, h)
	if err != nil {
		return nil, fmt.Errorf("failed to start webserver: %w", err)
	}
	cleanups = append(cleanups, func() { weber.Close() })

	return &Handler{
		configer:  configer,
		ruler:     ruler,
		connecter: connecter,
		queuer:    queuer,

		cmder: cmder,
		weber: weber,
	}, nil
}

func (h *Handler) Close() error {
	var errs []error
	h.closeOnce.Do(func() {
		// 停止 cmdserver 监听
		if err := h.cmder.Close(); err != nil {
			errs = append(errs, err)
		}
		// 停止 webserver 监听
		if err := h.weber.Close(); err != nil {
			errs = append(errs, err)
		}
		// 关闭队列
		if err := h.queuer.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := h.connecter.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := h.ruler.Close(); err != nil {
			errs = append(errs, err)
		}
		if err := h.configer.Close(); err != nil {
			errs = append(errs, err)
		}
	})
	return errors.Join(errs...)
}
