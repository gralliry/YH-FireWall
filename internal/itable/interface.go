package itable

import (
	"YH-FireWall/internal/pkg/fc"
	"net"
	"strings"
)

var (
	interfaces []*net.Interface

	name2index map[string]int
	index2name map[int]string

	ip2index map[string]int
)

func init() {
	interfaces = make([]*net.Interface, 0)

	name2index = make(map[string]int)
	index2name = make(map[int]string)

	ip2index = make(map[string]int)
}

func Init() error {
	// 获取网络接口
	itfs, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, itf := range itfs {
		interfaces = append(interfaces, &itf)
		name2index[itf.Name] = itf.Index
		index2name[itf.Index] = itf.Name

		addrs, _ := itf.Addrs()
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ipstr := ipnet.IP.String()
			ip2index[ipstr] = itf.Index
		}
	}
	return nil
}

func LookupByIp(ip string) (int, bool) {
	index, ok := ip2index[ip]
	return index, ok
}

func LookupByName(name string) (int, bool) {
	index, ok := name2index[name]
	return index, ok
}

func Index2ItfName(index *int) (string, bool) {
	if index == nil {
		return "", false
	}
	name, exist := index2name[*index]
	return name, exist
}

func Infos() []Info {
	return fc.List2List(interfaces, func(itf *net.Interface) Info {
		addrs, _ := itf.Addrs()
		return Info{
			Index: itf.Index,
			Name:  itf.Name,
			MAC:   itf.HardwareAddr.String(),
			MTU:   itf.MTU,
			Flags: strings.Split(itf.Flags.String(), "|"),
			Addrs: fc.List2List(addrs, func(addr net.Addr) string {
				return addr.String()
			}),
		}
	})
}
