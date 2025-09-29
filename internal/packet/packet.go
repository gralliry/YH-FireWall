package packet

import (
	"net"
)

type Packet struct {
	srcIP   net.IP
	srcPort uint16
	tarIP   net.IP
	tarPort uint16

	protocol uint8
}
