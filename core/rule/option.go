package rule

import (
	"math/rand"
)

type Option struct {
	// 不应该通过结构体中的id定位规则
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

// 自定义字符集
const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func id(n int) string {
	runes := []rune(alphabet) // 支持 Unicode
	length := len(runes)

	if n > length {
		n = length // 防止 n 太大
	}

	// Fisher–Yates 洗牌
	for i := length - 1; i > length-1-n; i-- {
		j := rand.Intn(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes[length-n:])
}

func (o *Option) Default() *Config {
	// 生成 8 位长度的 ID
	c := &Config{
		Id: id(8),
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
