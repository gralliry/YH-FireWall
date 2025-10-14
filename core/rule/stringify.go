package rule

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"net"
	"strings"
)

func stringifyIPNet(nets []net.IPNet) string {
	var parts []string
	for _, n := range nets {
		parts = append(parts, n.String()) // IPNet.String() 返回 "IP/Mask"
	}
	return strings.Join(parts, ",")
}

func stringifyPort(ports [][2]uint16) string {
	var parts []string
	for _, p := range ports {
		if p[0] == p[1] {
			parts = append(parts, fmt.Sprintf("%d", p[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d-%d", p[0], p[1]))
		}
	}
	return strings.Join(parts, ",")
}

func stringifyDev(devs map[uint32]struct{}) string {
	var parts []string
	for k := range devs {
		if ifi, err := net.InterfaceByIndex(int(k)); err == nil {
			parts = append(parts, ifi.Name)
		}
	}
	return strings.Join(parts, ",")
}

func stringifyProtocol(protocols map[layers.IPProtocol]struct{}) string {
	var parts []string
	for p := range protocols {
		parts = append(parts, strings.ToLower(string(p)))
	}
	return strings.Join(parts, ",")
}
