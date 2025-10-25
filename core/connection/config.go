package connection

import (
	"github.com/google/gopacket/layers"
	"net"
)

type Config struct {
	Id string `json:"id"`

	// 文件描述符
	Fd uint32 `json:"fd"`
	// 2: ipv4, 10: ipv6
	Family uint32 `json:"family"`
	// 6: tcp,  17:  udp
	Protocol layers.IPProtocol `json:"protocol"`

	Pid      int32  `json:"pid"`
	Exe      string `json:"exe"`
	Name     string `json:"name"`
	Cmd      string `json:"cmd"`
	Username string `json:"username"`
	// 看方向
	Interface string `json:"interface"`

	LocalIP   net.IP `json:"localIP"`
	LocalPort uint16 `json:"localPort"`

	// 0: inbound, 1: outbound, 2: forward
	Direction Direction `json:"direction"`

	RemoteIP   net.IP `json:"remoteIP"`
	RemotePort uint16 `json:"remotePort"`

	Status string `json:"status"`
	// 建立时间
	EstablishedTime int64 `json:"establishedTime"`
}
