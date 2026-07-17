package itable

import (
	"YH-FireWall/internal/model/itf"
	"YH-FireWall/internal/pkg/bimap"
	"net"
)

type Manager struct {
	interfaces []*itf.Itf
	pmap       *bimap.Map[string, uint32]
}

func New() (*Manager, error) {
	// 获取网络接口
	sItfs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var (
		interfaces = make([]*itf.Itf, 0)
		pmap       = bimap.New[string, uint32]()
	)
	for _, sItf := range sItfs {
		interfaces = append(interfaces, itf.New(&sItf))
	}
	for _, i := range interfaces {
		pmap.Insert(i.Name, i.Index)
	}
	return &Manager{
		interfaces: interfaces,
		pmap:       pmap,
	}, nil
}

func (t *Manager) List() []*itf.Itf {
	itfs := make([]*itf.Itf, len(t.interfaces))
	for i, itf := range t.interfaces {
		itfs[i] = itf.Clone()
	}
	return itfs
}

func (t *Manager) Name2Index(name string) (uint32, bool) {
	return t.pmap.GetByA(name)
}

func (t *Manager) Index2Name(index uint32) (string, bool) {
	return t.pmap.GetByB(index)
}
