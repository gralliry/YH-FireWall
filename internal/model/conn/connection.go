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

	localIP  netip.Addr
	remoteIP netip.Addr

	hasPort    bool
	localPort  uint16
	remotePort uint16

	direction flow.Direction

	// 状态信息
	establishedTime time.Time
	lastActiveTime  time.Time
	isClosed        bool
}

func New(f *flow.Flow) *Conn {
	c := pool.Get().(*Conn)

	c.direction = f.Direction()
	switch c.direction {
	case flow.Inbound:
		c.localIP, c.localPort = f.DstIP, f.DstPort
		c.remoteIP, c.remotePort = f.SrcIP, f.SrcPort
	case flow.Outbound:
		c.localIP, c.localPort = f.SrcIP, f.SrcPort
		c.remoteIP, c.remotePort = f.DstIP, f.DstPort
	default:
		pool.Put(c)
		return nil
	}
	// 按需处理
	c.id = sid.New(8)
	c.hasPort = f.HasPort
	//
	c.establishedTime = time.Now()
	c.lastActiveTime = time.Now()
	c.isClosed = false
	return c
}

func (c *Conn) Active() {
	c.lastActiveTime = time.Now()
}

//const killcmd = "src %s and src port %d and dst %s and dst port %d"

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

func (c *Conn) Info() *Info {
	return &Info{
		Id: c.id,

		// 连接信息
		Protocol: c.protocol,

		// 网卡方向信息

		// 状态
		EstablishedTime: c.establishedTime.Unix(),
	}
}

func (c *Conn) ID() string {
	return c.id
}

func key(proto layers.IPProtocol, ip1 netip.Addr, port1 uint16, ip2 netip.Addr, port2 uint16) string {
	return fmt.Sprintf("%s-%s-%d-%s-%d", proto, ip1, port1, ip2, port2)
}

func (c *Conn) LKey() string {
	return key(c.protocol, c.localIP, c.localPort, c.remoteIP, c.remotePort)
}

func (c *Conn) RKey() string {
	return key(c.protocol, c.remoteIP, c.remotePort, c.localIP, c.localPort)
}
