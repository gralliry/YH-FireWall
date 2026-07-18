package handler

import (
	"YH-FireWall/internal/config"
	"YH-FireWall/internal/constant/itfdev"
	"YH-FireWall/internal/ctable"
	"YH-FireWall/internal/queue"
	"YH-FireWall/internal/rtable"
	"YH-FireWall/internal/server/cmdserver"
	"YH-FireWall/internal/server/webserver"
	"errors"
	"fmt"
	"os"
)

type Handler struct {
	configer *config.Manager
	ruler    *rtable.Manager
	conner   *ctable.Manager
	queuer   *queue.NFQ
	//
	cmder *cmdserver.Server
	weber *webserver.Server
}

func New(configPath string) (*Handler, error) {
	h := &Handler{}
	// 检测当前用户是否为 root 用户
	if os.Geteuid() != 0 {
		return nil, errors.New("current user is not root")
	}

	// 读取配置管理器
	var err error
	h.configer, err = config.New(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	// 加载配置
	cfg := h.configer.Load()
	// 加载接口
	err = itfdev.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load interfaces: %w", err)
	}
	// 初始化规则管理器
	h.ruler, err = rtable.New(cfg.Rule)
	if err != nil {
		return nil, fmt.Errorf("failed to load rule table: %w", err)
	}
	// 初始化连接表
	h.conner = ctable.New()
	// 初始化队列
	h.queuer, err = queue.New(cfg.Queue, h)
	if err != nil {
		return nil, fmt.Errorf("failed to start queue: %w", err)
	}
	// 设置cmd服务器
	h.cmder, err = cmdserver.New(cfg.Cmd, h)
	if err != nil {
		return nil, fmt.Errorf("failed to start cmdserver: %w", err)
	}
	// 设置web服务器
	h.weber, err = webserver.New(cfg.Web, h)
	if err != nil {
		return nil, fmt.Errorf("failed to start webserver: %w", err)
	}
	return h, nil
}

func (h *Handler) Close() error {
	var errs []error
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
	if err := h.conner.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := h.ruler.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := h.configer.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
