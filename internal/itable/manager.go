package itable

import (
	"YH-FireWall/internal/model/itf"
	"net"
	"slices"
)

type Manager struct {
	interfaces []*itf.Itf

	name2index map[string]uint32
	index2name map[uint32]string
}

func New() (*Manager, error) {
	// 获取网络接口
	sItfs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var (
		interfaces = make([]*itf.Itf, 0)
		name2index = make(map[string]uint32)
		index2name = make(map[uint32]string)
	)
	for _, sItf := range sItfs {
		interfaces = append(interfaces, itf.New(&sItf))
	}
	for _, i := range interfaces {
		name2index[i.Name] = i.Index
		index2name[i.Index] = i.Name
	}
	return &Manager{
		interfaces: interfaces,

		name2index: name2index,
		index2name: index2name,
	}, nil
}

func (t *Manager) List() []*itf.Itf {
	// todo 这里有深拷贝问题
	return slices.Clone(t.interfaces)
}

func (t *Manager) Name2Index(name string) (uint32, bool) {
	index, ok := t.name2index[name]
	return index, ok
}

func (t *Manager) Index2Name(index uint32) (string, bool) {
	name, exist := t.index2name[index]
	return name, exist
}
