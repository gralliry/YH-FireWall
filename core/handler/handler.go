package handler

import (
	"YH-FireWall/core/config"
	"YH-FireWall/core/connection"
	"YH-FireWall/core/rtable"
	"YH-FireWall/core/rule"
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

func (h *Handler) AppendRule(ro *rule.Option) error {
	return rtable.Append(ro)
}
func (h *Handler) UpdateRule(id string, ro *rule.Option) error {
	return rtable.Update(id, ro)
}
func (h *Handler) DeleteRule(id string) error {
	return rtable.Delete(id)
}
func (h *Handler) GetRule(id string) *rule.Config {
	return rtable.Get(id)
}
func (h *Handler) GetRules() []rule.Config {
	return rtable.GetAll()
}
func (h *Handler) EnableRule(id string, enable bool) bool {
	return rtable.Enable(id, enable)
}
func (h *Handler) GetConfig() (string, error) {
	buf, err := config.Read()
	return string(buf), err
}

func (h *Handler) SetConfig(raw string) error {
	return config.Store([]byte(raw))
}

func (h *Handler) GetConnections() ([]connection.Connection, error) {
	return connection.GetAll()
}

func (h *Handler) CloseConnection(pid int32, fd uint32) error {
	return connection.Close(pid, fd)
}

func (h *Handler) GetInterfaces() ([]system.Interface, error) {
	return system.GetInterfaces()
}

func (h *Handler) Version() string {
	return config.Version
}
