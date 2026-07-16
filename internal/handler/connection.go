package handler

import (
	"YH-FireWall/internal/config"
	"YH-FireWall/internal/model/conn"
	"YH-FireWall/internal/model/itf"
)

// 连接模块
func (h *Handler) GetConnections() []conn.Info {
	return ctable.Infos()
}

func (h *Handler) CloseConnection(id string) error {
	return ctable.Remove(id)
}

// 接口模块
func (h *Handler) GetInterfaces() []itf.Itf {
	return itable.Infos()
}

// 版本
func (h *Handler) Version() string {
	return config.Version
}
