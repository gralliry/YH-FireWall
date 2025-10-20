package connection

import (
	"YH-FireWall/core/pkg/sid"
	"github.com/google/gopacket/layers"
	"net"
	"sync"
	"time"
)

type Direction int

const (
	Inbound  Direction = 0
	Outbound Direction = 1
	Forward  Direction = 2
)

type Connection struct {
	// 连接id
	id string
	// 2: ipv4, 10: ipv6
	family uint8

	protocol layers.IPProtocol

	localIP   net.IP
	localPort uint16

	direction Direction

	remoteIP   net.IP
	remotePort uint16

	establishedTime time.Time
	lastSeenTime    time.Time
	isClosed        bool

	//
	mutex sync.RWMutex
}

func New(
	family uint8, proto layers.IPProtocol,
	localIP net.IP, localPort uint16, direction Direction, remoteIP net.IP, remotePort uint16,
) *Connection {
	return &Connection{
		id:              sid.New(8),
		family:          family,
		protocol:        proto,
		localIP:         localIP,
		localPort:       localPort,
		direction:       direction,
		remoteIP:        remoteIP,
		remotePort:      remotePort,
		establishedTime: time.Now(),
		lastSeenTime:    time.Now(),
		isClosed:        false,
	}
}

func (c *Connection) Update() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lastSeenTime = time.Now()
}

func (c *Connection) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isClosed = true
}

func (c *Connection) Closed() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.closed()
}

func (c *Connection) closed() bool {
	return c.isClosed
}

func (c *Connection) Expired() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.expired()
}

func (c *Connection) expired() bool {
	return time.Now().Sub(c.lastSeenTime) > time.Minute
}
func (c *Connection) Status() (isClosed bool, isExpired bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.closed(), c.expired()
}

func (c *Connection) Unparse() *Config {
	return &Config{
		Id:              c.id,
		Family:          c.family,
		Protocol:        c.protocol,
		LocalIP:         c.localIP,
		LocalPort:       c.localPort,
		Direction:       c.direction,
		RemoteIP:        c.remoteIP,
		RemotePort:      c.remotePort,
		EstablishedTime: c.establishedTime.Unix(),
	}
}

func (c *Connection) Id() string {
	return c.id
}
