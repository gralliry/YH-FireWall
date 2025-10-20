package connection

import (
	"github.com/google/gopacket/layers"
	"net"
	"time"
)

type Config struct {
	Id string `json:"id"`

	// 2: ipv4, 10: ipv6
	Family uint8 `json:"family"`
	// 6: tcp,  17:  udp
	Protocol layers.IPProtocol `json:"protocol"`

	LocalIP   net.IP `json:"localIP"`
	LocalPort uint16 `json:"localPort"`

	// 0: inbound, 1: outbound, 2: forward
	Direction Direction `json:"direction"`

	RemoteIP   net.IP `json:"remoteIP"`
	RemotePort uint16 `json:"remotePort"`

	// 建立时间
	EstablishedTime time.Time `json:"establishedTime"`
}
