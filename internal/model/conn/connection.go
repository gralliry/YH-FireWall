package conn

import (
	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/pkg/sid"
	"fmt"
	"net/netip"
	"sync"
	"sync/atomic"
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

	direction flow.Direction

	lAddrPort netip.AddrPort
	rAddrPort netip.AddrPort

	// 状态信息
	establishTime int64
	activeTime    atomic.Int64
	closed        atomic.Bool
}

func New(f *flow.Flow) (*Conn, bool) {
	// 校验flow是否是连接包
	if !f.IsConnPackage() {
		return nil, false
	}
	// 获取conn
	c := pool.Get().(*Conn)
	// 按需处理
	c.id = sid.New(16)
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
		return nil, false
	}
	//
	cur := time.Now().UnixMilli()
	c.establishTime = cur
	c.activeTime.Store(cur)
	c.closed.Store(false)
	return c, true
}

func Release(c *Conn) {
	if c != nil {
		pool.Put(c)
	}
}

func (c *Conn) ID() string {
	return c.id
}

func (c *Conn) Active() {
	c.activeTime.Store(time.Now().Unix())
}

func (c *Conn) Close() {
	c.closed.Store(true)
}

func (c *Conn) Closed() bool {
	return c.closed.Load()
}

func (c *Conn) Expired() bool {
	const timeout = int64(time.Minute / time.Millisecond)
	return time.Now().UnixMilli()-timeout > c.activeTime.Load()
}

func (c *Conn) Alive() bool {
	return !c.Closed() && !c.Expired()
}

func (c *Conn) LKey() string {
	return fmt.Sprintf("%s-%s-%s", c.protocol, c.lAddrPort, c.rAddrPort)
}

func (c *Conn) RKey() string {
	return fmt.Sprintf("%s-%s-%s", c.protocol, c.rAddrPort, c.lAddrPort)
}
