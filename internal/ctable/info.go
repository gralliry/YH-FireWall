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

	// 进程信息
	Fd   int32  `json:"fd"`
	Pid  int32  `json:"pid"`
	Exe  string `json:"exe"`
	Name string `json:"name"`
	Cmd  string `json:"cmd"`
	User string `json:"user"`

	// 网卡方向信息
	InInterface  string         `json:"inInterface"`
	OutInterface string         `json:"outInterface"`
	Direction    flow.Direction `json:"direction"`

	// 建立时间
	EstablishedTime int64 `json:"establishedTime"`
}
