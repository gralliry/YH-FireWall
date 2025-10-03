package packet

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"net"
)

func (p *Packet) ID() uint32 {
	return p.id
}

func (p *Packet) UsePort() bool {
	return p.protocol == layers.IPProtocolTCP || p.protocol == layers.IPProtocolUDP
}

func (p *Packet) SrcIP() net.IP {
	return p.srcIP
}

func (p *Packet) DstIP() net.IP {
	return p.dstIP
}

func (p *Packet) SrcPort() uint16 {
	return p.srcPort
}

func (p *Packet) DstPort() uint16 {
	return p.dstPort
}

func (p *Packet) InDev() *uint32 {
	return p.inDev
}

func (p *Packet) OutDev() *uint32 {
	return p.outDev
}

func (p *Packet) Protocol() layers.IPProtocol {
	return p.protocol
}

func (p *Packet) String() string {
	return fmt.Sprintf("[%s] %s:%d -> %s:%d (%d -> %d)", p.protocol, p.srcIP, p.srcPort, p.dstIP, p.dstPort, p.inDev, p.outDev)
}
