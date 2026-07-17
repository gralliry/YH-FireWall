package conn

import (
	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/pkg/sid"
	"fmt"
	"net/netip"
	"sync"
	"time"

	"github.com/google/gopacket/layers"
)

var pool = sync.Pool{
	New: func() any { return new(Conn) },
}

type Conn struct {
	// 连接id
	id string

	// 连接信息
	protocol layers.IPProtocol

	lAddrPort netip.AddrPort
	rAddrPort netip.AddrPort

	direction flow.Direction

	// 状态信息
	establishTime  time.Time
	lastActiveTime time.Time
	isClosed       bool
}

func New(f *flow.Flow) *Conn {
	// 校验flow是否是连接包

	c := pool.Get().(*Conn)
	// 按需处理
	c.id = sid.New(8)
	//
	c.protocol = f.Protocol
	//
	c.direction = f.Direction()
	switch c.direction {
	case flow.Inbound:
		c.lAddrPort = netip.AddrPortFrom(f.DstIP, f.DstPort)
		c.rAddrPort = netip.AddrPortFrom(f.SrcIP, f.SrcPort)
	case flow.Outbound:
		c.lAddrPort = netip.AddrPortFrom(f.SrcIP, f.SrcPort)
		c.rAddrPort = netip.AddrPortFrom(f.DstIP, f.DstPort)
	default:
		Release(c)
		return nil
	}
	//
	c.establishTime = time.Now()
	c.lastActiveTime = time.Now()
	c.isClosed = false
	return c
}

func Release(c *Conn) {
	if c != nil {
		pool.Put(c)
	}
}

func (c *Conn) Active() {
	c.lastActiveTime = time.Now()
}

func (c *Conn) Close() error {
	c.lastActiveTime = time.Now()
	c.isClosed = true
	return nil
}

func (c *Conn) Closed() bool {
	return c.isClosed
}

func (c *Conn) Expired() bool {
	return time.Since(c.lastActiveTime) > time.Minute
}

func (c *Conn) ID() string {
	return c.id
}

func (c *Conn) LKey() string {
	return fmt.Sprintf("%s-%s-%s", c.protocol, c.lAddrPort, c.rAddrPort)
}

func (c *Conn) RKey() string {
	return fmt.Sprintf("%s-%s-%s", c.protocol, c.rAddrPort, c.lAddrPort)
}
