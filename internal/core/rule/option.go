package rule

import "YH-FireWall/internal/core/pkg/sid"

type Option struct {
	// 不应该通过结构体中的id定位规则
	// Id       *string `json:"id"`
	Group    *string `json:"group"`
	Comment  *string `json:"comment"`
	SrcNet   *string `json:"src_net"`
	SrcPort  *string `json:"src_port"`
	TarNet   *string `json:"tar_net"`
	TarPort  *string `json:"tar_port"`
	InDev    *string `json:"in_dev"`
	OutDev   *string `json:"out_dev"`
	Protocol *string `json:"protocol"`
	Accept   *bool   `json:"accept"`
	Priority *int    `json:"priority"`
	Enable   *bool   `json:"enable"`
}

func (o *Option) Default() *Config {
	c := &Config{
		Id: sid.Generate(),
	}
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
