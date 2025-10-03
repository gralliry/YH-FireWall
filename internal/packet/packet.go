package packet

import (
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

func Parse(a *nfqueue.Attribute) *Packet {
	p := &Packet{}
	id, payload := a.PacketID, a.Payload
	if id == nil || payload == nil {
		return nil
	}
	p.id = *id
	// 使用 gopacket 解析 Payload
	rawpacket := gopacket.NewPacket(*payload, layers.LayerTypeEthernet, gopacket.Default)
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
		return nil
	}
	// TCP/UDP 端口
	switch p.protocol {
	case layers.IPProtocolTCP:
		if tcp := rawpacket.Layer(layers.LayerTypeTCP); tcp != nil {
			t := tcp.(*layers.TCP)
			p.srcPort = uint16(t.SrcPort)
			p.dstPort = uint16(t.DstPort)
		} else {
			return nil
		}
	case layers.IPProtocolUDP:
		if udp := rawpacket.Layer(layers.LayerTypeUDP); udp != nil {
			u := udp.(*layers.UDP)
			p.srcPort = uint16(u.SrcPort)
			p.dstPort = uint16(u.DstPort)
		} else {
			return nil
		}
	}
	// 网口
	p.inDev = a.InDev
	p.outDev = a.OutDev
	//
	return p
}
