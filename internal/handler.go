package internal

import (
	"YH-FireWall/internal/config"
	"YH-FireWall/internal/ctable"
	"YH-FireWall/internal/itable"
	"YH-FireWall/internal/rtable"
	"YH-FireWall/internal/rule"
)

type Handler struct {
}

// 规则模块
func (h *Handler) AppendRule(ro *rule.Option) (string, error) {
	return rtable.Append(ro)
}
func (h *Handler) UpdateRule(id string, ro *rule.Option) error {
	return rtable.Update(id, ro)
}
func (h *Handler) DeleteRule(id string) error {
	return rtable.Delete(id)
}
func (h *Handler) SearchRule(id string) *rule.Info {
	return rtable.Search(id)
}
func (h *Handler) SearchRules() []rule.Info {
	return rtable.SearchAll()
}
func (h *Handler) EnableRule(id string, enable bool) bool {
	return rtable.Enable(id, enable)
}

// 配置模块
func (h *Handler) GetConfig() string {
	return string(config.Read())
}

func (h *Handler) SetConfig(raw string) error {
	return config.Save([]byte(raw))
}

// 连接模块
func (h *Handler) GetConnections() []ctable.Info {
	return ctable.Infos()
}

func (h *Handler) CloseConnection(id string) error {
	return ctable.Remove(id)
}

// 接口模块
func (h *Handler) GetInterfaces() []itable.Info {
	return itable.Infos()
}

// 协议常量
func (h *Handler) GetProtocols() []string {
	return rule.GetProtocolNames()
}

// 版本
func (h *Handler) Version() string {
	return config.Version
}
