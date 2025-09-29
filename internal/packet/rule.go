package packet

import (
	"net"
)

type Rule struct {
	srcNet  []net.IPNet
	srcPort []Port
	tarNet  []net.IPNet
	tarPort []Port

	protocol uint8

	accept bool

	priority int
}

func (r *Rule) Match(p *Packet) bool {
	//  protocol
	if p.protocol != r.protocol {
		return false
	}
	ok := false
	//
	for _, port := range r.srcPort {
		if port.Contains(p.srcPort) {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}
	ok = false
	for _, port := range r.tarPort {
		if port.Contains(p.tarPort) {
			ok = true
			break
		}
	}
	for _, n := range r.srcNet {
		if n.Contains(p.srcIP) {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}
	ok = false
	for _, n := range r.tarNet {
		if n.Contains(p.tarIP) {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}

	return true
}
