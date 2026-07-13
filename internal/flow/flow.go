package flow

import (
	"YH-FireWall/internal/itable"
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

type Flow struct {
	// 1: tcp,  2:  udp
	//type_ uint32
	// 2: ipv4, 10: ipv6
	Family uint8
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

func key(proto layers.IPProtocol, srcIP netip.Addr, srcPort uint16, dstIP netip.Addr, dstPort uint16) string {
	return fmt.Sprintf("%s-%s-%d->%s-%d", proto, srcIP, srcPort, dstIP, dstPort)
}

func (f *Flow) LKey() string {
	return key(f.Protocol, f.SrcIP, f.SrcPort, f.DstIP, f.DstPort)
}

func (f *Flow) RKey() string {
	return key(f.Protocol, f.DstIP, f.DstPort, f.SrcIP, f.SrcPort)
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

func (f *Flow) InDevName() string {
	name, _ := itable.Index2ItfName(int(f.InDev))
	return name
}

func (f *Flow) OutDevName() string {
	name, _ := itable.Index2ItfName(int(f.OutDev))
	return name
}

func (f *Flow) String() string {
	var direction string
	switch f.Direction() {
	case Forward:
		direction = fmt.Sprintf("%3d->%-3d", f.InDev, f.OutDev)
	case Outbound:
		direction = fmt.Sprintf("%3d->   ", f.OutDev)
	case Inbound:
		direction = fmt.Sprintf("   ->%-3d", f.InDev)
	default:
		direction = "   ->   "
	}
	return fmt.Sprintf("[%5s] %15s:%5d -> %15s:%5d (%s)",
		f.Protocol, f.SrcIP, f.SrcPort, f.DstIP, f.DstPort, direction)
}
