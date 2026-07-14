package flow

import (
	"net"
	"net/netip"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func setIPs(f *Flow, src, dst net.IP) bool {
	srcIP, ok1 := netip.AddrFromSlice(src)
	dstIP, ok2 := netip.AddrFromSlice(dst)
	if !ok1 || !ok2 {
		return false
	}
	f.SrcIP = srcIP.Unmap()
	f.DstIP = dstIP.Unmap()

	return true
}

func New(payload []byte, inDev, outDev *uint32) (*Flow, bool) {
	if len(payload) == 0 {
		return nil, false
	}

	var f Flow
	if inDev != nil {
		f.InDev = *inDev
	}
	if outDev != nil {
		f.OutDev = *outDev
	}

	// 判断 ip 协议
	var packet gopacket.Packet
	switch payload[0] >> 4 {
	case 4:
		packet = gopacket.NewPacket(payload, layers.LayerTypeIPv4, gopacket.DecodeOptions{
			Lazy:   true,
			NoCopy: true,
		})
		layer := packet.Layer(layers.LayerTypeIPv4).(*layers.IPv4)
		if !setIPs(&f, layer.SrcIP, layer.DstIP) {
			return nil, false
		}
		f.Protocol = layer.Protocol
	case 6:
		packet = gopacket.NewPacket(payload, layers.LayerTypeIPv6, gopacket.DecodeOptions{
			Lazy:   true,
			NoCopy: true,
		})
		layer := packet.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
		if !setIPs(&f, layer.SrcIP, layer.DstIP) {
			return nil, false
		}
		f.Protocol = layer.NextHeader
	default:
		return nil, false
	}

	// 判断传输层协议
	var ok bool
	f.SrcPort, f.DstPort, f.HasPort, ok = extractPort(packet, f.Protocol)
	if !ok {
		return nil, false
	}

	return &f, true
}

func extractPort(packet gopacket.Packet, proto layers.IPProtocol) (uint16, uint16, bool, bool) {
	switch proto {
	case layers.IPProtocolTCP:
		if l, ok := packet.Layer(layers.LayerTypeTCP).(*layers.TCP); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true, true
		}
	case layers.IPProtocolUDP:
		if l, ok := packet.Layer(layers.LayerTypeUDP).(*layers.UDP); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true, true
		}
	case layers.IPProtocolSCTP:
		if l, ok := packet.Layer(layers.LayerTypeSCTP).(*layers.SCTP); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true, true
		}
	case layers.IPProtocolUDPLite:
		if l, ok := packet.Layer(layers.LayerTypeUDPLite).(*layers.UDPLite); ok {
			return uint16(l.SrcPort), uint16(l.DstPort), true, true
		}
	default:
		return 0, 0, false, true
	}
	// 错误
	return 0, 0, true, false
}
