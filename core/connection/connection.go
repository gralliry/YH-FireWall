package connection

import (
	"fmt"
	"github.com/google/gopacket/layers"
)

type Direction int

const (
	Inbound Direction = iota
	Outbound
	Forward
)

type Connection interface {
	Update()
	Protocol() layers.IPProtocol
}

func New(param *Parameter) (Connection, error) {
	switch param.Proto {
	case layers.IPProtocolTCP:
		return NewTcp(param), nil
	case layers.IPProtocolUDP:
		return NewUdp(param), nil
	default:
		return nil, fmt.Errorf("unsupported protocol")
	}
}
