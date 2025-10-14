package connection

import (
	"github.com/google/gopacket/layers"
	"net"
	"time"
)

type TcpConnection struct {
	SrcIP   net.IP
	SrcPort uint16
	DstIP   net.IP
	DstPort uint16

	Direction Direction

	LastSeen time.Time
	TTL      time.Duration
}

func NewTcp(param *Parameter) *TcpConnection {
	return &TcpConnection{
		SrcIP:     param.SrcIP,
		SrcPort:   param.SrcPort,
		DstIP:     param.DstIP,
		DstPort:   param.DstPort,
		Direction: param.Direction,
		LastSeen:  time.Now(),
		TTL:       time.Minute * 5,
	}
}

func (t *TcpConnection) Update() {
	t.LastSeen = time.Now()
}

func (t *TcpConnection) Protocol() layers.IPProtocol {
	return layers.IPProtocolTCP
}
