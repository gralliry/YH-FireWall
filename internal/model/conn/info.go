package conn

import (
	"YH-FireWall/internal/model/flow"
	"net/netip"

	"github.com/google/gopacket/layers"
)

type Info struct {
	// 连接id
	Id string `json:"id"`

	// 连接信息
	Protocol  layers.IPProtocol `json:"protocol"`
	LAddrPort netip.AddrPort    `json:"localIP"`
	RAddrPort netip.AddrPort    `json:"remoteIP"`

	// 网卡方向信息
	Interface string         `json:"interface"`
	Direction flow.Direction `json:"direction"`

	// 建立时间
	EstablishedTime int64 `json:"establishedTime"`
}
