package itable

import (
	"YH-FireWall/internal/model/itf"
	"net"
	"net/netip"
	"sync"
)

var (
	mutex sync.RWMutex

	interfaces []*itf.Itf

	name2index = make(map[string]int)
	index2name = make(map[int]string)
	addr2index = make(map[netip.Addr]int)
)

func Refresh() error {
	// 获取网络接口
	itfs, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, i := range itfs {
		interfaces = append(interfaces, itf.New(&i))
	}
	for _, i := range interfaces {
		name2index[i.Name] = i.Index
		index2name[i.Index] = i.Name

		for _, prefix := range i.Addrs {
			addr2index[prefix.Addr()] = i.Index
		}
	}
	return nil
}

func IP2Index(ip netip.Addr) (int, bool) {
	index, ok := addr2index[ip]
	return index, ok
}

func Name2Index(name string) (int, bool) {
	index, ok := name2index[name]
	return index, ok
}

func Index2Name(index int) (string, bool) {
	name, exist := index2name[index]
	return name, exist
}
