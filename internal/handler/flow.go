package handler

import (
	"YH-FireWall/internal/model/flow"
	"errors"

	"github.com/florianl/go-nfqueue"
)

func (h *Handler) Handle(a *nfqueue.Attribute) (bool, error) {
	f, ok := flow.New(a)
	if !ok {
		return false, errors.New("s")
	}
	// 回收flow
	defer flow.Release(f)
	// 匹配连接
	if h.conner.Match(f) {
		return true, nil
	}
	// 匹配规则表
	action := h.ruler.Match(f)
	// 回收flow
	flow.Release(f)
	return action, nil
}
