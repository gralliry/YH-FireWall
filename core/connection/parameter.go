package connection

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"net"
)

type Parameter struct {
	Proto   layers.IPProtocol
	SrcIP   net.IP
	SrcPort uint16
	DstIP   net.IP
	DstPort uint16

	Direction Direction
}

func (p *Parameter) Key() string {
	return buildKey(p.Proto, p.SrcIP, p.SrcPort, p.DstIP, p.DstPort)
}

func (p *Parameter) ReverseKey() string {
	return buildKey(p.Proto, p.DstIP, p.DstPort, p.SrcIP, p.SrcPort)
}

func buildKey(proto layers.IPProtocol, srcIP net.IP, srcPort uint16, dstIP net.IP, dstPort uint16) string {
	return fmt.Sprintf("%s|%s:%d-%s:%d", proto, srcIP, srcPort, dstIP, dstPort)
}
