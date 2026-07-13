package rule

import (
	"YH-FireWall/internal/pkg/sid"
	"fmt"
	"strings"
)

type Info struct {
	Id       string `json:"id"`
	Group    string `json:"group"`
	Comment  string `json:"comment"`
	SrcNet   string `json:"srcNet"`
	SrcPort  string `json:"srcPort"`
	TarNet   string `json:"tarNet"`
	TarPort  string `json:"tarPort"`
	InDev    string `json:"inDev"`
	OutDev   string `json:"outDev"`
	Protocol string `json:"protocol"`
	Accept   bool   `json:"accept"`
	Priority int    `json:"priority"`
	Enable   bool   `json:"enable"`
}

func (c *Info) String() string {
	var sb strings.Builder
	indent := "    "
	fmt.Fprintf(&sb, "Index: %s Group: %s Comment: %s\n", c.Id, c.Group, c.Comment)
	fmt.Fprintf(&sb, "%sSrcNet  : %s\n", indent, c.SrcNet)
	fmt.Fprintf(&sb, "%sSrcPort : %s\n", indent, c.SrcPort)
	fmt.Fprintf(&sb, "%sTarNet  : %s\n", indent, c.TarNet)
	fmt.Fprintf(&sb, "%sTarPort : %s\n", indent, c.TarPort)
	fmt.Fprintf(&sb, "%sInDev   : %s\n", indent, c.InDev)
	fmt.Fprintf(&sb, "%sOutDev  : %s\n", indent, c.OutDev)
	fmt.Fprintf(&sb, "%sProtocol: %s\n", indent, c.Protocol)
	fmt.Fprintf(&sb, "%sAccept  : %t\n", indent, c.Accept)
	fmt.Fprintf(&sb, "%sPriority: %d\n", indent, c.Priority)
	fmt.Fprintf(&sb, "%sEnable  : %t\n", indent, c.Enable)
	return sb.String()
}

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
