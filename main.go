package main

import (
	"YH-FireWall/internal"
)

//func pass(mark uint32, packet gopacket.Packet) bool {
//	var (
//		srcIP, dstIP     net.IP
//		srcPort, dstPort uint16
//		protocol         layers.IPProtocol
//	)
//	// 获取 IPv4 或 IPv6 地址
//	if ip4 := packet.Layer(layers.LayerTypeIPv4); ip4 != nil {
//		ip := ip4.(*layers.IPv4)
//		srcIP = ip.SrcIP
//		dstIP = ip.DstIP
//		protocol = ip.Protocol
//	} else if ip6 := packet.Layer(layers.LayerTypeIPv6); ip6 != nil {
//		ip := ip6.(*layers.IPv6)
//		srcIP = ip.SrcIP
//		dstIP = ip.DstIP
//		protocol = ip.NextHeader
//	} else {
//		return false
//	}
//	// TCP/UDP 端口
//	switch protocol {
//	case layers.IPProtocolTCP:
//		if tcp := packet.Layer(layers.LayerTypeTCP); tcp != nil {
//			t := tcp.(*layers.TCP)
//			srcPort = uint16(t.SrcPort)
//			dstPort = uint16(t.DstPort)
//		} else {
//			return false
//		}
//	case layers.IPProtocolUDP:
//		if udp := packet.Layer(layers.LayerTypeUDP); udp != nil {
//			u := udp.(*layers.UDP)
//			srcPort = uint16(u.SrcPort)
//			dstPort = uint16(u.DstPort)
//		} else {
//			return false
//		}
//	}
//	return false
//}

func main() {
	internal.Start()
}
