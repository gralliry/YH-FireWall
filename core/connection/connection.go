package connection

import (
	"YH-FireWall/core/pkg/sid"
	"errors"
	"fmt"
	"github.com/google/gopacket/layers"
	"net"
	"sync"
	"time"
)

type Connection struct {
	// 连接id
	id string
	// 2: ipv4, 10: ipv6
	family uint32
	// 1: tcp,  2:  udp
	//type_ uint32
	// 6: "TCP", 17: "UDP",
	protocol layers.IPProtocol
	// 接口
	interface_ string

	// 进程信息
	isProcessInfoEmpty bool
	fd                 uint32
	pid                int32
	exe                string
	name               string
	cmd                string
	username           string

	localIP   net.IP
	localPort uint16

	direction Direction

	remoteIP   net.IP
	remotePort uint16

	status string

	establishedTime time.Time
	lastActiveTime  time.Time
	isClosed        bool

	//
	mutex sync.RWMutex
}

func NewByPush(
	family uint32, proto layers.IPProtocol,
	localIP net.IP, localPort uint16, direction Direction, remoteIP net.IP, remotePort uint16,
	interface_ string,
) *Connection {
	return &Connection{
		id:       sid.New(8),
		family:   family,
		protocol: proto,

		isProcessInfoEmpty: true,

		localIP:    localIP,
		localPort:  localPort,
		direction:  direction,
		remoteIP:   remoteIP,
		remotePort: remotePort,

		interface_: interface_,
		status:     "ESTABLISHED",

		establishedTime: time.Now(),
		lastActiveTime:  time.Now(),
		isClosed:        false,
	}
}

func (c *Connection) UpdateByProcess(
	fd uint32, pid int32, exe string, name string, cmd string, username string, status string,
) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.fd = fd
	c.pid = pid
	c.exe = exe
	c.name = name
	c.cmd = cmd
	c.username = username

	c.status = status

	c.lastActiveTime = time.Now()
}

func MakeKey(proto layers.IPProtocol, srcIP net.IP, srcPort uint16, dstIP net.IP, dstPort uint16) string {
	return fmt.Sprintf("%s-%s-%d-%s-%d", proto, srcIP, srcPort, dstIP, dstPort)
}

func (c *Connection) UpdateByPush() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lastActiveTime = time.Now()
}

//const killcmd = "src %s and src port %d and dst %s and dst port %d"

func (c *Connection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	//
	if c.status != "ESTABLISHED" {
		return errors.New("not ESTABLISHED")
	}
	c.isClosed = true
	return nil
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
	return time.Now().Sub(c.lastActiveTime) > time.Minute
}
func (c *Connection) Status() (isClosed bool, isExpired bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.closed(), c.expired()
}

func (c *Connection) Unparse() *Config {
	return &Config{
		Id: c.id,

		Family:   c.family,
		Protocol: c.protocol,

		Fd:       c.fd,
		Pid:      c.pid,
		Exe:      c.exe,
		Name:     c.name,
		Cmd:      c.cmd,
		Username: c.username,

		Interface: c.interface_,

		LocalIP:   c.localIP,
		LocalPort: c.localPort,

		Direction: c.direction,

		RemoteIP:   c.remoteIP,
		RemotePort: c.remotePort,

		Status:          c.status,
		EstablishedTime: c.establishedTime.Unix(),
	}
}

func (c *Connection) Id() string {
	return c.id
}

func (c *Connection) LKey() string {
	return MakeKey(c.protocol, c.localIP, c.localPort, c.remoteIP, c.remotePort)
}

func (c *Connection) RKey() string {
	return MakeKey(c.protocol, c.remoteIP, c.remotePort, c.localIP, c.localPort)
}

func (c *Connection) IsProcessInfoEmpty() bool {
	return c.isProcessInfoEmpty
}
