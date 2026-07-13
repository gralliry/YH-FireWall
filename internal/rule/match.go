package rule

import (
	"net"

	"github.com/google/gopacket/layers"
)

func matchIPNet(nets []net.IPNet, n net.IP) bool {
	if len(nets) == 0 {
		return true
	}
	for _, ipnet := range nets {
		if ipnet.Contains(n) {
			return true
		}
	}
	return false
}

func matchPort(ports [][2]uint16, port uint16) bool {
	if len(ports) == 0 {
		return true
	}
	for _, pair := range ports {
		if pair[0] <= port && port <= pair[1] {
			return true
		}
	}
	return false
}

func matchDev(devs map[uint32]struct{}, d *uint32) bool {
	if devs == nil || d == nil {
		return true
	}
	_, exists := devs[*d]
	return exists
}

func matchProtocol(ps map[layers.IPProtocol]struct{}, p layers.IPProtocol) bool {
	if ps == nil {
		return true
	}
	_, exists := ps[p]
	return exists
}
