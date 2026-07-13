package rule

import (
	"YH-FireWall/internal/pkg/sid"
)

type Option struct {
	Group    *string `json:"group"`
	Comment  *string `json:"comment"`
	SrcNet   *string `json:"srcNet"`
	SrcPort  *string `json:"srcPort"`
	TarNet   *string `json:"tarNet"`
	TarPort  *string `json:"tarPort"`
	InDev    *string `json:"inDev"`
	OutDev   *string `json:"outDev"`
	Protocol *string `json:"protocol"`
	Accept   *bool   `json:"accept"`
	Priority *int    `json:"priority"`
	Enable   *bool   `json:"enable"`
}

func (o *Option) Build() *Info {
	c := &Info{Id: sid.New(8)}
	if o.Group != nil {
		c.Group = *o.Group
	}
	if o.Comment != nil {
		c.Comment = *o.Comment
	}
	if o.SrcNet != nil {
		c.SrcNet = *o.SrcNet
	}
	if o.SrcPort != nil {
		c.SrcPort = *o.SrcPort
	}
	if o.TarNet != nil {
		c.TarNet = *o.TarNet
	}
	if o.TarPort != nil {
		c.TarPort = *o.TarPort
	}
	if o.InDev != nil {
		c.InDev = *o.InDev
	}
	if o.OutDev != nil {
		c.OutDev = *o.OutDev
	}
	if o.Protocol != nil {
		c.Protocol = *o.Protocol
	}
	if o.Accept != nil {
		c.Accept = *o.Accept
	}
	if o.Enable != nil {
		c.Enable = *o.Enable
	}
	return c
}
