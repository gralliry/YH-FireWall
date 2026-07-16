package handler

// 配置模块
func (h *Handler) GetConfig() string {
	return string(config.Read())
}

func (h *Handler) SetConfig(raw string) error {
	return config.Save([]byte(raw))
}
