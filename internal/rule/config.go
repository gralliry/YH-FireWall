package rule

import (
	"fmt"
	"strings"
)

type Config struct {
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

func (c *Config) String() string {
	var sb strings.Builder
	indent := "    "
	sb.WriteString(fmt.Sprintf("Index: %s Group: %s Comment: %s\n", c.Id, c.Group, c.Comment))
	sb.WriteString(fmt.Sprintf("%sSrcNet  : %s\n", indent, c.SrcNet))
	sb.WriteString(fmt.Sprintf("%sSrcPort : %s\n", indent, c.SrcPort))
	sb.WriteString(fmt.Sprintf("%sTarNet  : %s\n", indent, c.TarNet))
	sb.WriteString(fmt.Sprintf("%sTarPort : %s\n", indent, c.TarPort))
	sb.WriteString(fmt.Sprintf("%sInDev   : %s\n", indent, c.InDev))
	sb.WriteString(fmt.Sprintf("%sOutDev  : %s\n", indent, c.OutDev))
	sb.WriteString(fmt.Sprintf("%sProtocol: %s\n", indent, c.Protocol))
	sb.WriteString(fmt.Sprintf("%sAccept  : %t\n", indent, c.Accept))
	sb.WriteString(fmt.Sprintf("%sPriority: %d\n", indent, c.Priority))
	sb.WriteString(fmt.Sprintf("%sEnable  : %t\n", indent, c.Enable))
	return sb.String()
}
