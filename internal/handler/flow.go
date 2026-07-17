package handler

import (
	"YH-FireWall/internal/model/flow"

	"github.com/florianl/go-nfqueue"
)

func (h *Handler) HandleFlow(a *nfqueue.Attribute) (bool, bool) {
	return h.Handle(a)
}

func (h *Handler) Handle(a *nfqueue.Attribute) (accept bool, ok bool) {
	f, ok := flow.New(a)
	if !ok {
		return false, false
	}
	// 回收flow
	defer flow.Release(f)
	// 匹配 并 更新 连接
	if accept, exist := h.conner.Match(f); exist {
		return accept, true
	}
	// 匹配规则表
	accept = h.ruler.Match(f)
	if accept {
		h.conner.Push(f)
	}
	return accept, true
}
