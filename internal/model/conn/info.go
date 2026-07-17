package conn

import (
	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/model/proc"
	"net/netip"

	"github.com/google/gopacket/layers"
)

type Info struct {
	// 连接id
	ID string `json:"id"`

	// 连接信息
	Protocol  layers.IPProtocol `json:"protocol"`
	LAddrPort netip.AddrPort    `json:"localIP"`
	RAddrPort netip.AddrPort    `json:"remoteIP"`

	// 网卡方向信息
	Direction flow.Direction `json:"direction"`

	// 建立时间
	EstablishTime int64 `json:"establishTime"`

	// 进程信息
	Process proc.Info `json:"process"`
}

func (c *Conn) Info(pid int32) *Info {
	info := &Info{
		ID:        c.id,
		Protocol:  c.protocol,
		LAddrPort: c.lAddrPort,
		RAddrPort: c.rAddrPort,

		Direction: c.direction,

		EstablishTime: c.establishTime.Unix(),
	}
	if pinfo, err := proc.NewByPID(pid); err == nil {
		info.Process = *pinfo
	}
	return info
}
