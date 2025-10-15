package handler

import (
	"YH-FireWall/core/config"
	"YH-FireWall/core/rtable"
	rule2 "YH-FireWall/core/rule"
	"YH-FireWall/core/server/cmdserver"
	"YH-FireWall/core/server/webserver"
	"YH-FireWall/core/system"
	"context"
)

type Handler struct {
	Context context.Context
	Cancel  context.CancelFunc
}

func (h *Handler) Start() error {
	return nil
}

func (h *Handler) Stop() error {
	h.Cancel()
	return nil
}

func (h *Handler) AppendRule(ro *rule2.Option) error {
	return rtable.Append(ro)
}
func (h *Handler) UpdateRule(id string, ro *rule2.Option) error {
	return rtable.Update(id, ro)
}
func (h *Handler) DeleteRule(id string) error {
	return rtable.Delete(id)
}
func (h *Handler) GetRule(id string) *rule2.Config {
	return rtable.GetRule(id)
}
func (h *Handler) GetRules() []rule2.Config {
	return rtable.GetRules()
}
func (h *Handler) EnableRule(id string, enable bool) bool {
	return rtable.SetAbleRule(id, enable)
}
func (h *Handler) GetConfig() (string, error) {
	buf, err := config.Read()
	return string(buf), err
}

func (h *Handler) SetConfig(raw string) error {
	return config.Store([]byte(raw))
}

func (h *Handler) GetConnections() ([]system.Connection, error) {
	return system.GetConnections()
}

func (h *Handler) GetInterfaces() ([]system.Interface, error) {
	return system.GetInterface()
}

func (h *Handler) CmdStart() error {
	return cmdserver.Start(h)
}

func (h *Handler) CmdIsRunning() bool {
	return cmdserver.IsRunning()
}

func (h *Handler) CmdStop() error {
	return cmdserver.Close()
}

func (h *Handler) WebStart() error {
	return webserver.Start(h)
}
func (h *Handler) WebIsRunning() bool {
	return webserver.IsRunning()
}
func (h *Handler) WebStop() error {
	return webserver.Close()
}

func (h *Handler) Version() string {
	return config.Version
}
