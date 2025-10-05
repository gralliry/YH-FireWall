package rule

import (
	"strconv"
	"strings"
)

type Config struct {
	Name     string `json:"name"`
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

func (c *Config) StringList() []string {
	return []string{
		"name: " + c.Name,
		"comment: " + c.Comment,
		"src_net: " + c.SrcNet,
		"src_port: " + c.SrcPort,
		"tar_net: " + c.TarNet,
		"tar_port: " + c.TarPort,
		"in_dev: " + c.InDev,
		"out_dev: " + c.OutDev,
		"protocol: " + c.Protocol,
		"accept: " + strconv.FormatBool(c.Accept),
		"priority: " + strconv.Itoa(c.Priority),
		"enable: " + strconv.FormatBool(c.Enable),
	}
}

func (c *Config) String() string {
	return strings.Join(c.StringList(), "\n") + "\n"
}
