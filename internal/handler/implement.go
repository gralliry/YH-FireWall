package handler

import (
	"YH-FireWall/internal/config"
	"YH-FireWall/internal/model/conn"
	"YH-FireWall/internal/model/itf"
	"YH-FireWall/internal/model/rule"
)

// 规则模块
func (h *Handler) CreateRule(ro *rule.Option) (string, error) {
	return h.ruler.Create(ro)
}
func (h *Handler) UpdateRule(id string, ro *rule.Option) error {
	return h.ruler.Update(id, ro)
}
func (h *Handler) DeleteRule(id string) error {
	return h.ruler.Delete(id)
}
func (h *Handler) SearchRule(id string) *rule.Data {
	return h.ruler.Search(id)
}
func (h *Handler) EnableRule(id string, enable bool) error {
	return h.ruler.Enable(id, enable)
}
func (h *Handler) ListRules() []*rule.Data {
	return h.ruler.List()
}

// 连接模块
func (h *Handler) ListConnections() []*conn.Info {
	return h.conner.List()
}
func (h *Handler) CloseConnection(id string) error {
	return h.conner.Remove(id)
}

// 配置模块
func (h *Handler) GetConfig() string {
	return h.configer.Read()
}
func (h *Handler) SetConfig(raw string) error {
	return h.configer.Write(raw)
}

// 接口模块
func (h *Handler) ListInterfaces() []*itf.Itf {
	return h.itfer.List()
}

// 协议模块
func (h *Handler) ListProtocols() []string {
	return []string{}
}

// 版本
func (h *Handler) Version() string {
	return config.Version
}
