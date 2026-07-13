package ctable

import (
	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/pkg/sid"
	"time"
)

type Connection struct {
	// 连接id
	id string

	// 连接流信息
	flow *flow.Flow
	// 状态信息
	establishedTime time.Time
	lastActiveTime  time.Time
	isClosed        bool
}

func New(flow *flow.Flow) *Connection {
	return &Connection{
		id: sid.New(8),

		flow: flow,

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

func (c *Connection) Unparse() *Info {
	return &Info{
		Id: c.id,
		// 连接信息
		Family:     uint32(c.flow.Family),
		Protocol:   c.flow.Protocol,
		LocalIP:    c.flow.SrcIP,
		LocalPort:  c.flow.SrcPort,
		RemoteIP:   c.flow.SrcIP,
		RemotePort: c.flow.SrcPort,

		// 网卡信息
		InInterface:  c.flow.InDevName(),
		OutInterface: c.flow.OutDevName(),
		Direction:    c.flow.Direction(),

		// 状态
		EstablishedTime: c.establishedTime.Unix(),
	}
}

func (c *Connection) Id() string {
	return c.id
}

func (c *Connection) LKey() string {
	return c.flow.LKey()
}

func (c *Connection) RKey() string {
	return c.flow.RKey()
}
