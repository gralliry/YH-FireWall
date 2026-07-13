package queue

import (
	"YH-FireWall/internal/ctable"
	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/rtable"
	"log"

	"github.com/florianl/go-nfqueue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Handler interface {
	Match(flow *flow.Flow) bool
	Update(flow *flow.Flow) bool
}

var dConfig = gopacket.DecodeOptions{
	Lazy:   true,
	NoCopy: true,
}

func handleFunc(a nfqueue.Attribute) int {
	flow := flow.Flow{}
	payload := *a.Payload
	if len(payload) == 0 {
		nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
		return 0
	}
	var (
		packet gopacket.Packet
	)
	if payload[0]>>4 == 6 {
		packet = gopacket.NewPacket(payload, layers.LayerTypeIPv6, gopacket.Default)
		layer, ok := packet.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
		if !ok {
			nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
			return 0
		}
		flow.SrcIP, flow.DstIP, flow.Protocol, flow.Family = layer.SrcIP, layer.DstIP, layer.NextHeader, 10
	} else {
		packet = gopacket.NewPacket(payload, layers.LayerTypeIPv4, gopacket.Default)
		layer, ok := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
		if !ok {
			nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
			return 0
		}
		flow.SrcIP, flow.DstIP, flow.Protocol, flow.Family = layer.SrcIP, layer.DstIP, layer.Protocol, 2
	}
	// 匹配端口
	if src, dst, ok := extractPort(packet, flow.Protocol); ok {
		flow.SrcPort, flow.DstPort = src, dst
	}
	// 匹配规则
	if !rtable.Match(&flow) {
		nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
		return 0
	}
	// 更新连接表
	ctable.Push(&flow)
	// 打印日志
	log.Print(flow.String())
	// 继续处理
	nfq.SetVerdict(*a.PacketID, nfqueue.NfAccept)
	return 0
}

func errorFunc(err error) int {
	return -1
}

func parseFlow(payload []byte) (flow.Flow, bool) {
	var (
		f      flow.Flow
		packet gopacket.Packet
	)

	if len(payload) == 0 {
		return f, false
	}

	packet = gopacket.NewPacket(payload, layers.LayerTypeEthernet, gopacket.Default)

	switch payload[0] >> 4 {
	case 4:
		a := gopacket.DecodeOptions{
			Lazy:   false,
			NoCopy: false,
		}
		packet = gopacket.NewPacket(payload, layers.LayerTypeEthernet, gopacket.Default)

		ip, ok := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
		if !ok {
			return f, false
		}

		f.Family = 2
		f.Protocol = ip.Protocol
		f.SrcIP = ip.SrcIP
		f.DstIP = ip.DstIP

	case 6:
		packet = gopacket.NewPacket(payload, layers.LayerTypeIPv6, gopacket.Default)

		ip, ok := packet.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
		if !ok {
			return f, false
		}

		f.Family = 10
		f.Protocol = ip.NextHeader
		f.SrcIP = ip.SrcIP
		f.DstIP = ip.DstIP

	default:
		return f, false
	}

	if src, dst, ok := extractPort(packet, f.Protocol); ok {
		f.SrcPort = src
		f.DstPort = dst
	}

	return f, true
}

func extractPort(packet gopacket.Packet, proto layers.IPProtocol) (uint16, uint16, bool) {
	switch proto {
	case layers.IPProtocolTCP:
		if l, ok := packet.Layer(layers.LayerTypeTCP).(*layers.TCP); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true
		}

	case layers.IPProtocolUDP:
		if l, ok := packet.Layer(layers.LayerTypeUDP).(*layers.UDP); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true
		}

	case layers.IPProtocolSCTP:
		if l, ok := packet.Layer(layers.LayerTypeSCTP).(*layers.SCTP); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true
		}

	case layers.IPProtocolUDPLite:
		if l, ok := packet.Layer(layers.LayerTypeUDPLite).(*layers.UDPLite); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true
		}
	}

	return 0, 0, false
}
