package rule

import (
	"github.com/google/gopacket/layers"
	"net"
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
	if len(devs) == 0 {
		return true
	}
	if d == nil {
		return false
	}
	_, exists := devs[*d]
	return exists
}

func matchProtocol(ps map[layers.IPProtocol]struct{}, p layers.IPProtocol) bool {
	if len(ps) == 0 {
		return true
	}
	_, exists := ps[p]
	return exists
}
