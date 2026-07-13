package flow

import (
	"fmt"
	"net"

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
	// 2: ipv4, 10: ipv6
	Family uint32
	// 1: tcp,  2:  udp
	//type_ uint32
	// 6: tcp,  17:  udp
	Protocol layers.IPProtocol

	SrcIP   net.IP
	SrcPort uint16
	DstIP   net.IP
	DstPort uint16

	// 网卡设备
	InDev  *uint32
	OutDev *uint32
}

func key(proto layers.IPProtocol, srcIP net.IP, srcPort uint16, dstIP net.IP, dstPort uint16) string {
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
	case f.InDev != nil && f.OutDev != nil:
		// 转发模式
		return Forward
	case f.InDev != nil:
		// 入口模式 // 源是外部连接
		return Inbound
	case f.OutDev != nil:
		// 出口模式 // 源是内部连接
		return Outbound
	default:
		// 未知数据，直接停止
		return Unknown
	}
}

func (f *Flow) InDevName() string {
	if f.inDevName == nil {
		name := index2interface(f.InDev)
		f.inDevName = &name
	}
	return *f.inDevName
}

func (f *Flow) OutDevName() string {
	if f.outDevName == nil {
		name := index2interface(f.OutDev)
		f.outDevName = &name
	}
	return *f.outDevName
}

func (f *Flow) String() string {
	var direction string
	switch {
	case f.InDev == nil && f.OutDev == nil:
		direction = "   ->   "
	case f.InDev == nil:
		direction = fmt.Sprintf("%3d->   ", *f.OutDev)
	case f.OutDev == nil:
		direction = fmt.Sprintf("   ->%-3d", *f.InDev)
	default:
		direction = fmt.Sprintf("%3d->%-3d", *f.InDev, *f.OutDev)
	}
	return fmt.Sprintf("[%5s] %15s:%5d -> %15s:%5d (%s)",
		f.Protocol, f.SrcIP, f.SrcPort, f.DstIP, f.DstPort, direction)
}
