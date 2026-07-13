package flow

import (
	"net/netip"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func New(payload []byte, inDev, outDev *uint32) (*Flow, bool) {
	var (
		f      Flow
		packet gopacket.Packet
	)
	if len(payload) == 0 {
		return nil, false
	}

	// 注册网卡
	if inDev != nil {
		f.InDev = *inDev
	}
	if outDev != nil {
		f.OutDev = *outDev
	}

	switch payload[0] >> 4 {
	case 4:
		// IPv4
		packet = gopacket.NewPacket(payload, layers.LayerTypeEthernet, gopacket.Default)

		ip, ok := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
		if !ok {
			return nil, false
		}

		f.Family = 2
		f.Protocol = ip.Protocol

		srcIP, ok1 := netip.AddrFromSlice(ip.SrcIP)
		dstIP, ok2 := netip.AddrFromSlice(ip.SrcIP)

		if !ok1 || !ok2 {
			return nil, false
		}

		f.SrcIP = srcIP
		f.DstIP = dstIP

	case 6:
		// IPv6
		packet = gopacket.NewPacket(payload, layers.LayerTypeIPv6, gopacket.Default)

		ip, ok := packet.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
		if !ok {
			return nil, false
		}

		f.Family = 10
		f.Protocol = ip.NextHeader

		srcIP, ok1 := netip.AddrFromSlice(ip.SrcIP)
		dstIP, ok2 := netip.AddrFromSlice(ip.SrcIP)

		if !ok1 || !ok2 {
			return nil, false
		}

		f.SrcIP = srcIP
		f.DstIP = dstIP

	default:
		// 其他协议不支持
		return nil, false
	}

	// 是否使用端口
	if src, dst, hasPort, ok := extractPort(packet, f.Protocol); !ok {
		return nil, false
	} else if hasPort {
		f.HasPort = true
		f.SrcPort = src
		f.DstPort = dst
	}

	return &f, true
}

func extractPort(packet gopacket.Packet, proto layers.IPProtocol) (uint16, uint16, bool, bool) {
	switch proto {
	case layers.IPProtocolTCP:
		if l, ok := packet.Layer(layers.LayerTypeTCP).(*layers.TCP); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true, true
		}
		return 0, 0, true, false

	case layers.IPProtocolUDP:
		if l, ok := packet.Layer(layers.LayerTypeUDP).(*layers.UDP); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true, true
		}
		return 0, 0, true, false

	case layers.IPProtocolSCTP:
		if l, ok := packet.Layer(layers.LayerTypeSCTP).(*layers.SCTP); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true, true
		}
		return 0, 0, true, false

	case layers.IPProtocolUDPLite:
		if l, ok := packet.Layer(layers.LayerTypeUDPLite).(*layers.UDPLite); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true, true
		}
		return 0, 0, true, false

	default:
		// 没有端口的协议
		return 0, 0, false, true
	}
}
