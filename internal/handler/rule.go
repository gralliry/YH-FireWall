package handler

// 规则模块
func (h *Handler) AppendRule(ro *rule.Option) (string, error) {
	return h.ruler.Append(ro)
}
func (h *Handler) UpdateRule(id string, ro *rule.Option) error {
	return h.ruler.Update(id, ro)
}
func (h *Handler) DeleteRule(id string) error {
	return h.ruler.Delete(id)
}
func (h *Handler) SearchRule(id string) *rule.Info {
	return h.ruler.Search(id)
}
func (h *Handler) SelectRules() []rule.Info {
	return h.ruler.SelectAll()
}
func (h *Handler) EnableRule(id string, enable bool) bool {
	return h.ruler.Enable(id, enable)
}
