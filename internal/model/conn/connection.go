package conn

import (
	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/pkg/sid"
	"fmt"
	"net/netip"
	"time"

	"github.com/google/gopacket/layers"
)

type Connection struct {
	// 连接id
	id string

	// 连接信息
	protocol layers.IPProtocol

	localIP  netip.Addr
    remoteIP netip.Addr

    hasPort  bool
    localPort uint16
    remotePort uint16

	direction flow.Direction

	// 状态信息
	establishedTime time.Time
	lastActiveTime  time.Time
	isClosed        bool
}

func New(f *flow.Flow) *Connection {
	srcAddrPort := netip.AddrPortFrom(f.SrcIP, f.SrcPort)
	dstAddrPort := netip.AddrPortFrom(f.DstIP, f.DstPort)
	var lAddrPort, rAddrPort netip.AddrPort
	direction := f.Direction()
	switch direction {
	case flow.Inbound:
		lAddrPort = dstAddrPort
		rAddrPort = srcAddrPort
	case flow.Outbound:
		lAddrPort = srcAddrPort
		rAddrPort = dstAddrPort
	default:
		return nil
	} // Forward/Unknown 按需处理
	return &Connection{
		id: sid.New(8),

		protocol:  f.Protocol,
		lAddrPort: lAddrPort,
		rAddrPort: rAddrPort,

		direction: direction,

		establishedTime: time.Now(),
		lastActiveTime:  time.Now(),
		isClosed:        false,
	}
}

func (c *Connection) Active() {
	c.lastActiveTime = time.Now()
}

//const killcmd = "src %s and src port %d and dst %s and dst port %d"

func (c *Connection) Close() error {
	c.isClosed = true
	return nil
}

func (c *Connection) Closed() bool {
	return c.isClosed
}

func (c *Connection) Expired() bool {
	return time.Since(c.lastActiveTime) > time.Minute
}

func (c *Connection) Alive() bool {
	return !c.isClosed && time.Since(c.lastActiveTime) < time.Minute
}

func (c *Connection) Info() *Info {
	return &Info{
		Id: c.id,

		// 连接信息
		Protocol:  c.protocol,
		LAddrPort: c.lAddrPort,
		RAddrPort: c.rAddrPort,

		// 网卡方向信息
		Interface:
		Direction:c.direction ,

		// 状态
		EstablishedTime: c.establishedTime.Unix(),
	}
}

func (c *Connection) Id() string {
	return c.id
}

func key(proto layers.IPProtocol, ip1 netip.Addr, port1 uint16, ip2 netip.Addr, port2 uint16) string {
	return fmt.Sprintf("%s-%s-%d-%s-%d", proto, ip1, port1, ip2, port2)
}

func (c *Connection) LKey() string {
	return key(c.protocol, c.localIP, c.localPort, c.remoteIP, c.remotePort)
}

func (c *Connection) RKey() string {
	return key(c.protocol, c.remoteIP, c.remotePort, c.localIP, c.localPort)
}
