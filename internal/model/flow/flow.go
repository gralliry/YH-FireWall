package flow

import (
	"fmt"
	"net/netip"

	"github.com/florianl/go-nfqueue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Direction uint8

const (
	Inbound  Direction = 0
	Outbound Direction = 1
	Forward  Direction = 2
	Unknown  Direction = 3
)

// flow 不应该出现任何修改
type Flow struct {
	// IP上协议
	// 6: tcp,  17:  udp
	Protocol layers.IPProtocol

	// IPv4 IPv6
	SrcIP netip.Addr
	DstIP netip.Addr

	// TCP UDP SCTP DCCP UDPLite
	HasPort bool
	SrcPort uint16
	DstPort uint16

	// 网卡设备
	// 0 代表不存在
	InDev  uint32
	OutDev uint32
}

func New(a *nfqueue.Attribute) (*Flow, bool) {
	// a 一定不为空
	if a.Payload == nil {
		return nil, false
	}
	payload := *a.Payload
	if len(payload) == 0 {
		return nil, false
	}

	// 解析到flow
	f := pool.Get().(*Flow)
	f.InDev = 0
	f.OutDev = 0
	if a.InDev != nil {
		f.InDev = *a.InDev
	}
	if a.OutDev != nil {
		f.OutDev = *a.OutDev
	}

	// 判断 ip 协议
	var packet gopacket.Packet
	detectProto := func(fallback layers.IPProtocol) layers.IPProtocol {
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
	switch payload[0] >> 4 {
	case 4:
		packet = gopacket.NewPacket(payload, layers.LayerTypeIPv4, gopacket.DecodeOptions{
			Lazy:   true,
			NoCopy: true,
		})
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			Release(f)
			return nil, false
		}
		layer := ipLayer.(*layers.IPv4)
		if !setIP(f, layer.SrcIP, layer.DstIP) {
			Release(f)
			return nil, false
		}
		f.Protocol = detectProto(layer.Protocol)
	case 6:
		packet = gopacket.NewPacket(payload, layers.LayerTypeIPv6, gopacket.DecodeOptions{
			Lazy:   true,
			NoCopy: true,
		})
		ipLayer := packet.Layer(layers.LayerTypeIPv6)
		if ipLayer == nil {
			Release(f)
			return nil, false
		}
		layer := ipLayer.(*layers.IPv6)
		if !setIP(f, layer.SrcIP, layer.DstIP) {
			Release(f)
			return nil, false
		}
		f.Protocol = detectProto(layer.NextHeader)
	default:
		Release(f)
		return nil, false
	}

	// 判断传输层协议
	var ok bool
	f.SrcPort, f.DstPort, f.HasPort, ok = extractPort(packet, f.Protocol)
	if !ok {
		Release(f)
		return nil, false
	}

	return f, true
}

func Release(f *Flow) {
	if f != nil {
		pool.Put(f)
	} else {
		panic("Something try to release a nil pointer of flow")
	}
}

func (f *Flow) Key() string {
	return fmt.Sprintf("%s-%s-%s", f.Protocol, f.SrcAddrPort(), f.DstAddrPort())
}

func (f *Flow) IsConnection() bool {
	return f.Protocol == layers.IPProtocolTCP || f.Protocol == layers.IPProtocolUDP
}

func (f *Flow) Direction() Direction {
	switch {
	case f.InDev != 0 && f.OutDev != 0:
		// 转发模式
		return Forward
	case f.InDev != 0:
		// 入口模式 // 源是外部连接
		return Inbound
	case f.OutDev != 0:
		// 出口模式 // 源是内部连接
		return Outbound
	default:
		// 未知数据，直接停止
		return Unknown
	}
}

func (f *Flow) SrcAddrPort() netip.AddrPort {
	return netip.AddrPortFrom(f.SrcIP, f.SrcPort)
}

func (f *Flow) DstAddrPort() netip.AddrPort {
	return netip.AddrPortFrom(f.DstIP, f.DstPort)
}
