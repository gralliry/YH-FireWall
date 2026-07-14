package flow

import (
	"fmt"
	"net/netip"

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
//
//	不应该调用其他模块库
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

func key(proto layers.IPProtocol, ip1 netip.Addr, port1 uint16, ip2 netip.Addr, port2 uint16) string {
	return fmt.Sprintf("%s-%s-%d-%s-%d", proto, ip1, port1, ip2, port2)
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
