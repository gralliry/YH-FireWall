package connection

import (
	"github.com/google/gopacket/layers"
	"net"
	"time"
)

const (
	UdpTimeoutUnreplied = 30 * time.Second
	UdpTimeoutReplied   = 180 * time.Second
)

type UdpConnection struct {
	SrcIP   net.IP
	SrcPort uint16
	DstIP   net.IP
	DstPort uint16

	Direction Direction

	LastSeen time.Time
	TTL      time.Duration
}

func NewUdp(param *Parameter) *UdpConnection {
	return &UdpConnection{
		SrcIP:     param.SrcIP,
		SrcPort:   param.SrcPort,
		DstIP:     param.DstIP,
		DstPort:   param.DstPort,
		Direction: param.Direction,
		LastSeen:  time.Now(),
		TTL:       UdpTimeoutUnreplied,
	}
}

func (u *UdpConnection) Update() {
	u.LastSeen = time.Now()
	u.TTL = UdpTimeoutReplied
}

func (u *UdpConnection) Protocol() layers.IPProtocol {
	return layers.IPProtocolUDP
}
