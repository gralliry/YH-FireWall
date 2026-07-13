package rule

import (
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
