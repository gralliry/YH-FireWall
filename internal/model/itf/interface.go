package itf

import (
	"net"
	"net/netip"
)

type Itf struct {
	Index int            `json:"index"`
	Name  string         `json:"name"`
	MAC   string         `json:"mac"`
	MTU   int            `json:"mtu"`
	Flags net.Flags      `json:"flags"`
	Addrs []netip.Prefix `json:"addrs"`
}

func New(i *net.Interface) *Itf {
	addrs, _ := i.Addrs()
	prefixes := make([]netip.Prefix, 0, len(addrs))
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ip, ok := netip.AddrFromSlice(ipnet.IP)
		if !ok {
			continue
		}
		ones, _ := ipnet.Mask.Size()
		prefixes = append(prefixes, netip.PrefixFrom(ip.Unmap(), ones))
	}
	return &Itf{
		Index: i.Index,
		Name:  i.Name,
		MAC:   i.HardwareAddr.String(),
		MTU:   i.MTU,
		Flags: i.Flags,
		Addrs: prefixes,
	}
}
