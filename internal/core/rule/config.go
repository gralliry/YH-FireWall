package rule

import (
	"YH-FireWall/internal/core/pkg/sid"
	"fmt"
	"strings"
)

type Config struct {
	Id       string `json:"id"`
	Group    string `json:"group"`
	Comment  string `json:"comment"`
	SrcNet   string `json:"src_net"`
	SrcPort  string `json:"src_port"`
	TarNet   string `json:"tar_net"`
	TarPort  string `json:"tar_port"`
	InDev    string `json:"in_dev"`
	OutDev   string `json:"out_dev"`
	Protocol string `json:"protocol"`
	Accept   bool   `json:"accept"`
	Priority int    `json:"priority"`
	Enable   bool   `json:"enable"`
}

func (c *Config) Refresh() {
	c.Id = sid.Generate()
}

func (c *Config) String() string {
	var sb strings.Builder
	indent := "    "
	sb.WriteString(fmt.Sprintf("Id: %s Group: %s Comment: %s\n", c.Id, c.Group, c.Comment))
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
