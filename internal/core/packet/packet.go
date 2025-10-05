package packet

import (
	"fmt"
	"github.com/florianl/go-nfqueue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
)

type Packet struct {
	id uint32

	srcIP net.IP
	dstIP net.IP

	srcPort uint16
	dstPort uint16

	inDev    *uint32
	outDev   *uint32
	protocol layers.IPProtocol
}

func Parse(a *nfqueue.Attribute) (*Packet, error) {
	p := &Packet{}
	if a.PacketID == nil {
		return nil, fmt.Errorf("invalid packet")
	}
	p.id = *a.PacketID
	if a.Payload == nil {
		return p, fmt.Errorf("invalid payload")
	}
	// 使用 gopacket 解析 Payload
	rawpacket := gopacket.NewPacket(*a.Payload, layers.LayerTypeEthernet, gopacket.Default)
	// 获取 IPv4 或 IPv6 地址
	if ip4 := rawpacket.Layer(layers.LayerTypeIPv4); ip4 != nil {
		ip := ip4.(*layers.IPv4)
		p.srcIP = ip.SrcIP
		p.dstIP = ip.DstIP
		p.protocol = ip.Protocol
	} else if ip6 := rawpacket.Layer(layers.LayerTypeIPv6); ip6 != nil {
		ip := ip6.(*layers.IPv6)
		p.srcIP = ip.SrcIP
		p.dstIP = ip.DstIP
		p.protocol = ip.NextHeader
	} else {
		return p, fmt.Errorf("invalid IP")
	}
	// TCP/UDP 端口
	switch p.protocol {
	case layers.IPProtocolTCP:
		if tcp := rawpacket.Layer(layers.LayerTypeTCP); tcp != nil {
			t := tcp.(*layers.TCP)
			p.srcPort = uint16(t.SrcPort)
			p.dstPort = uint16(t.DstPort)
		} else {
			return p, fmt.Errorf("invalid Protocal")
		}
	case layers.IPProtocolUDP:
		if udp := rawpacket.Layer(layers.LayerTypeUDP); udp != nil {
			u := udp.(*layers.UDP)
			p.srcPort = uint16(u.SrcPort)
			p.dstPort = uint16(u.DstPort)
		} else {
			return p, fmt.Errorf("invalid Protocal")
		}
	}
	// 网口
	p.inDev = a.InDev
	p.outDev = a.OutDev
	//
	return p, nil
}

func (p *Packet) Id() uint32 {
	return p.id
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

func (p *Packet) UsePort() bool {
	return p.protocol == layers.IPProtocolTCP || p.protocol == layers.IPProtocolUDP
}

func (p *Packet) String() string {
	return fmt.Sprintf("[%s] %s:%d -> %s:%d (%d -> %d)",
		p.protocol, p.srcIP, p.srcPort, p.dstIP, p.dstPort, p.inDev, p.outDev)
}
