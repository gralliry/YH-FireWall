package ctable

import (
	"YH-FireWall/internal/model/flow"
	"net"

	"github.com/google/gopacket/layers"
)

type Info struct {
	// 连接id
	Id string `json:"id"`

	// 连接信息
	Family     uint32            `json:"family"`
	Protocol   layers.IPProtocol `json:"protocol"`
	LocalIP    net.IP            `json:"localIP"`
	LocalPort  uint16            `json:"localPort"`
	RemoteIP   net.IP            `json:"remoteIP"`
	RemotePort uint16            `json:"remotePort"`

	// 网卡方向信息
	InInterface  string         `json:"inInterface"`
	OutInterface string         `json:"outInterface"`
	Direction    flow.Direction `json:"direction"`

	// 建立时间
	EstablishedTime int64 `json:"establishedTime"`
}
