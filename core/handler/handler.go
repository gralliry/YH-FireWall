package handler

import (
	"YH-FireWall/core/config"
	"YH-FireWall/core/manager"
	"YH-FireWall/core/rule"
	"YH-FireWall/core/server/http"
	"YH-FireWall/core/server/unix"
	"context"
	"errors"
	"github.com/jinzhu/copier"
)

type Handler struct {
	Config  *config.Config
	Context context.Context
	Cancel  context.CancelFunc
}

func (h *Handler) Start() error {
	// 默认启动 unix 监听
	if err := h.UnixStart(); err != nil {
		return err
	}
	// 默认启动 http 监听
	if err := h.WebStart(); err != nil {
		return err
	}
	return nil
}

func (h *Handler) Stop() error {
	var errs []error
	// 停止 unix 监听
	if unix.IsRunning() {
		errs = append(errs, unix.Close())
	}
	// 停止 http 监听
	if http.IsRunning() {
		errs = append(errs, http.Close())
	}
	// 停止进程
	if h.Cancel != nil {
		h.Cancel()
	}
	return errors.Join(errs...)
}

func (h *Handler) AppendRule(ro *rule.Option) error {
	return manager.AppendRule(ro)
}
func (h *Handler) UpdateRule(id string, ro *rule.Option) error {
	return manager.UpdateRule(id, ro)
}
func (h *Handler) DeleteRule(id string) error {
	return manager.DeleteRule(id)
}
func (h *Handler) GetRule(id string) *rule.Config {
	return manager.GetRule(id)
}
func (h *Handler) GetRules() []rule.Config {
	return manager.GetRules()
}
func (h *Handler) EnableRule(id string, enable bool) bool {
	return manager.EnableRule(id, enable)
}
func (h *Handler) EnableGroup(group string, enable bool) bool {
	return manager.EnableGroup(group, enable)
}
func (h *Handler) GetConfig() *config.Config {
	d := &config.Config{}
	_ = copier.Copy(d, h.Config)
	return d
}

func (h *Handler) UnixStart() error {
	// 启动 unix 监听 // 里面
	return unix.Start(h.Config.Unix, h)
}

func (h *Handler) UnixIsRunning() bool {
	return unix.IsRunning()
}

func (h *Handler) UnixStop() error {
	return unix.Close()
}

func (h *Handler) WebStart() error {
	return http.Start(h.Config.Web, h)
}
func (h *Handler) WebIsRunning() bool {
	return http.IsRunning()
}
func (h *Handler) WebStop() error {
	return http.Close()
}

func (h *Handler) Version() string {
	return config.Version
}
