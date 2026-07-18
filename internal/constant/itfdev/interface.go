package itfdev

import (
	"net"
	"sort"
)

var (
	index2name = make(map[uint32]string)
	name2index = make(map[string]uint32)
)

func Load() error {
	itfs, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, i := range itfs {
		index := uint32(i.Index)
		index2name[index] = i.Name
		name2index[i.Name] = index
	}
	return nil
}

func List() []string {
	names := make([]string, 0, len(name2index))
	for name := range name2index {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func Index2Name(index uint32) (string, bool) {
	name, exist := index2name[index]
	return name, exist
}

func Name2Index(name string) (uint32, bool) {
	index, exist := name2index[name]
	return index, exist
}
