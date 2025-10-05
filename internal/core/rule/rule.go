package rule

import (
	"YH-FireWall/internal/core/packet"
	"fmt"
	"github.com/google/gopacket/layers"
	"net"
)

type Rule struct {
	name     string
	comment  string
	accept   bool
	priority int
	enable   bool

	// 动态解析
	srcNets []net.IPNet
	dstNets []net.IPNet

	srcPorts [][2]uint16
	dstPorts [][2]uint16

	inDevs  map[uint32]struct{}
	outDevs map[uint32]struct{}

	protocols map[layers.IPProtocol]struct{}
}

func Parse(cfg *Config) (*Rule, error) {
	if cfg == nil {
		return nil, fmt.Errorf("rule config is nil")
	}
	if cfg.Name == "" {
		return nil, fmt.Errorf("rule name is empty")
	}
	r := &Rule{
		name:     cfg.Name,
		comment:  cfg.Comment,
		accept:   cfg.Accept,
		priority: cfg.Priority,
		enable:   cfg.Enable,
	}
	// 源网络
	if srcNet, err := parseIPNet(cfg.SrcNet); err == nil {
		r.srcNets = srcNet
	} else {
		return nil, err
	}
	// 目标网络
	if tarNet, err := parseIPNet(cfg.TarNet); err == nil {
		r.dstNets = tarNet
	} else {
		return nil, err
	}
	// 源端口
	if srcPort, err := parsePort(cfg.SrcPort); err == nil {
		r.srcPorts = srcPort
	} else {
		return nil, err
	}
	// 目标端口
	if tarPort, err := parsePort(cfg.TarPort); err == nil {
		r.dstPorts = tarPort
	} else {
		return nil, err
	}
	// 入口
	if inDev, err := parseDev(cfg.InDev); err == nil {
		r.inDevs = inDev
	} else {
		return nil, err
	}
	// 出口
	if outDev, err := parseDev(cfg.OutDev); err == nil {
		r.outDevs = outDev
	} else {
		return nil, err
	}
	// 协议
	if protocols, err := parseProtocol(cfg.Protocol); err == nil {
		r.protocols = protocols
	} else {
		return nil, err
	}
	return r, nil
}

func (r *Rule) Unparse() *Config {
	return &Config{
		Name:     r.name,
		Comment:  r.comment,
		SrcNet:   stringifyIPNet(r.srcNets),
		SrcPort:  stringifyPort(r.srcPorts),
		TarNet:   stringifyIPNet(r.dstNets),
		TarPort:  stringifyPort(r.dstPorts),
		InDev:    stringifyDev(r.inDevs),
		OutDev:   stringifyDev(r.outDevs),
		Protocol: stringifyProtocol(r.protocols),
		Accept:   r.accept,
		Priority: r.priority,
		Enable:   r.enable,
	}
}

func (r *Rule) Match(p *packet.Packet) bool {
	// 这里可以通过架构去优化，减少if次数
	if !r.enable {
		return false
	}
	// 匹配 入口网卡 // 为空默认跳过该检查
	if !matchDev(r.inDevs, p.InDev()) {
		return false
	}
	// 匹配 出口网卡 // 为空默认跳过该检查
	if !matchDev(r.outDevs, p.OutDev()) {
		return false
	}
	// 匹配 协议 // 为空默认跳过该检查
	if !matchProtocol(r.protocols, p.Protocol()) {
		return false
	}
	// 匹配 源 IP
	if !matchIPNet(r.srcNets, p.SrcIP()) {
		return false
	}
	// 匹配 目标 IP
	if !matchIPNet(r.dstNets, p.DstIP()) {
		return false
	}
	// 匹配 端口（如果这个协议有端口） // 仅支持tcp udp
	if p.UsePort() {
		// 匹配 源端口
		if !matchPort(r.srcPorts, p.SrcPort()) {
			return false
		}
		if !matchPort(r.dstPorts, p.DstPort()) {
			return false
		}
	}
	return true
}

func (r *Rule) Name() string {
	return r.name
}

func (r *Rule) Priority() int {
	return r.priority
}

func (r *Rule) Enable() bool {
	return r.enable
}

func (r *Rule) Accept() bool {
	return r.accept
}
