package itf

import (
	"net"
	"net/netip"
)

type Itf struct {
	Index uint32   `json:"index"`
	Name  string   `json:"name"`
	MAC   string   `json:"mac"`
	MTU   int      `json:"mtu"`
	Flags []string `json:"flags"`
	Addrs []string `json:"addrs"`
}

func List() ([]Itf, error) {
	itfs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	out := make([]Itf, 0, len(itfs))
	for _, i := range itfs {
		out = append(out, newItf(&i))
	}
	return out, nil
}

func newItf(i *net.Interface) Itf {
	addrs, _ := i.Addrs()
	prefixes := make([]string, 0, len(addrs))
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
		prefixes = append(prefixes, netip.PrefixFrom(ip.Unmap(), ones).String())
	}
	return Itf{
		Index: uint32(i.Index),
		Name:  i.Name,
		MAC:   i.HardwareAddr.String(),
		MTU:   i.MTU,
		Flags: flagsToStrings(i.Flags),
		Addrs: prefixes,
	}
}

func flagsToStrings(f net.Flags) []string {
	var out []string
	if f&net.FlagUp != 0 {
		out = append(out, "up")
	}
	if f&net.FlagBroadcast != 0 {
		out = append(out, "broadcast")
	}
	if f&net.FlagLoopback != 0 {
		out = append(out, "loopback")
	}
	if f&net.FlagPointToPoint != 0 {
		out = append(out, "pointtopoint")
	}
	if f&net.FlagMulticast != 0 {
		out = append(out, "multicast")
	}
	if f&net.FlagRunning != 0 {
		out = append(out, "running")
	}
	return out
}
