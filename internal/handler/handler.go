package handler

import (
	"YH-FireWall/internal/config"
	"YH-FireWall/internal/connection"
	_const "YH-FireWall/internal/const"
	"YH-FireWall/internal/ctable"
	"YH-FireWall/internal/iface"
	"YH-FireWall/internal/rtable"
	"YH-FireWall/internal/rule"
	"context"
)

type Handler struct {
	Cancel context.CancelFunc
}

func (h *Handler) Start() error {
	return nil
}

func (h *Handler) Stop() error {
	h.Cancel()
	return nil
}

func (h *Handler) AppendRule(ro *rule.Option) (string, error) {
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
func (h *Handler) GetConfig() string {
	return string(config.Read())
}

func (h *Handler) SetConfig(raw string) error {
	return config.Save([]byte(raw))
}

func (h *Handler) GetConnections() []connection.Config {
	return ctable.GetAll()
}

func (h *Handler) CloseConnection(id string) error {
	return ctable.Remove(id)
}

func (h *Handler) GetInterfaces() ([]iface.Config, error) {
	return iface.GetAll()
}

func (h *Handler) GetProtocols() []string {
	return rule.GetAllProtocolNames()
}

func (h *Handler) Version() string {
	return _const.Version
}
