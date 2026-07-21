package flow

import (
	"net"
	"net/netip"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var pool = sync.Pool{
	New: func() any { return new(Flow) },
}

func setIP(f *Flow, src, dst net.IP) bool {
	srcIP, ok1 := netip.AddrFromSlice(src)
	dstIP, ok2 := netip.AddrFromSlice(dst)
	if !ok1 || !ok2 {
		return false
	}
	f.SrcIP = srcIP.Unmap()
	f.DstIP = dstIP.Unmap()

	return true
}

func detectProto(packet gopacket.Packet, fallback layers.IPProtocol) layers.IPProtocol {
	switch {
	case packet.Layer(layers.LayerTypeTCP) != nil:
		return layers.IPProtocolTCP
	case packet.Layer(layers.LayerTypeUDP) != nil:
		return layers.IPProtocolUDP
	case packet.Layer(layers.LayerTypeSCTP) != nil:
		return layers.IPProtocolSCTP
	case packet.Layer(layers.LayerTypeUDPLite) != nil:
		return layers.IPProtocolUDPLite
	default:
		return fallback
	}
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
