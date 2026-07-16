package flow

import (
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
	if a.InDev != nil {
		f.InDev = *a.InDev
	}
	if a.OutDev != nil {
		f.OutDev = *a.OutDev
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
		if !setIP(f, layer.SrcIP, layer.DstIP) {
			return nil, false
		}
		f.Protocol = layer.Protocol
	case 6:
		packet = gopacket.NewPacket(payload, layers.LayerTypeIPv6, gopacket.DecodeOptions{
			Lazy:   true,
			NoCopy: true,
		})
		layer := packet.Layer(layers.LayerTypeIPv6).(*layers.IPv6)
		if !setIP(f, layer.SrcIP, layer.DstIP) {
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

	return f, true
}

func Release(f *Flow) {
	pool.Put(f)
}

func (f *Flow) Key() string {
	return key(f.Protocol, f.SrcIP, f.SrcPort, f.DstIP, f.DstPort)
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
