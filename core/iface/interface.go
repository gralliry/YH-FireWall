package iface

import (
	"YH-FireWall/core/pkg/fc"
	"net"
	"strings"
	"sync"
)

type Interface struct {
	Index int        `json:"index"`
	Name  string     `json:"name"`
	MAC   string     `json:"mac"`
	MTU   int        `json:"mtu"`
	Flags []string   `json:"flags"`
	Addrs []net.Addr `json:"addrs"`
}

var (
	interfaces []Interface
	mutex      sync.RWMutex
)

func Get(reflush bool) ([]Interface, error) {
	mutex.Lock()
	defer mutex.Unlock()
	// 读取缓存
	if !reflush && interfaces != nil {
		return interfaces, nil
	}
	// 获取网络接口
	nitf, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	// 转换为配置
	itfs := fc.List2List(nitf, func(ifs net.Interface) Interface {
		oaddrs, _ := ifs.Addrs()
		return Interface{
			Index: ifs.Index,
			Name:  ifs.Name,
			MAC:   ifs.HardwareAddr.String(),
			MTU:   ifs.MTU,
			Flags: strings.Split(ifs.Flags.String(), "|"),
			Addrs: oaddrs,
		}
	})

	interfaces = itfs

	return itfs, nil
}

func GetAll() ([]Config, error) {
	itfs, err := Get(false)
	if err != nil {
		return nil, err
	}
	return fc.List2List(itfs, func(ifs Interface) Config {
		return Config{
			Index: ifs.Index,
			Name:  ifs.Name,
			MAC:   ifs.MAC,
			MTU:   ifs.MTU,
			Flags: ifs.Flags,
			Addrs: fc.List2List(ifs.Addrs, func(addr net.Addr) string {
				return addr.String()
			}),
		}
	}), nil
}

func FindNameByIp(ip net.IP) string {
	mutex.RLock()
	defer mutex.RUnlock()
	for _, itf := range interfaces {
		for _, addr := range itf.Addrs {
			if addr.(*net.IPNet).Contains(ip) {
				return itf.Name
			}
		}
	}
	return ""
}
