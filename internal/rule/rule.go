package rule

import (
	"YH-FireWall/internal/packet"
	"github.com/google/gopacket/layers"
	"net"
)

type Rule struct {
	id          uint32
	description string
	accept      bool
	enable      bool
	priority    int

	srcNet  []net.IPNet
	srcPort [][2]uint16
	dstNet  []net.IPNet
	dstPort [][2]uint16

	inDev    map[uint32]struct{}
	outDev   map[uint32]struct{}
	protocol map[layers.IPProtocol]struct{}
}

func Parse(cfg Config) (*Rule, error) {
	r := &Rule{
		id:          cfg.Id,
		description: cfg.Description,
		accept:      cfg.Accept,
		enable:      cfg.Enable,
		priority:    int(cfg.Priority),
	}
	// 源网络
	if srcNet, err := parseIPNet(cfg.SrcNet); err == nil {
		r.srcNet = srcNet
	} else {
		return nil, err
	}
	// 目标网络
	if tarNet, err := parseIPNet(cfg.TarNet); err == nil {
		r.dstNet = tarNet
	} else {
		return nil, err
	}
	// 源端口
	if srcPort, err := parsePort(cfg.SrcPort); err == nil {
		r.srcPort = srcPort
	} else {
		return nil, err
	}
	// 目标端口
	if tarPort, err := parsePort(cfg.TarPort); err == nil {
		r.dstPort = tarPort
	} else {
		return nil, err
	}
	// 入口
	if inDev, err := parseDev(cfg.InDev); err == nil {
		r.inDev = inDev
	} else {
		return nil, err
	}
	// 出口
	if outDev, err := parseDev(cfg.OutDev); err == nil {
		r.outDev = outDev
	} else {
		return nil, err
	}
	// 协议
	if protocols, err := parseProtocol(cfg.Protocol); err == nil {
		r.protocol = protocols
	} else {
		return nil, err
	}
	return r, nil
}

func (r *Rule) Match(p *packet.Packet) bool {
	// 这里可以通过架构去优化，减少if次数
	if !r.enable {
		return false
	}
	// 匹配 入口网卡 // 为空默认跳过该检查
	if !matchDev(r.inDev, p.InDev()) {
		return false
	}
	// 匹配 出口网卡 // 为空默认跳过该检查
	if !matchDev(r.outDev, p.OutDev()) {
		return false
	}
	// 匹配 协议 // 为空默认跳过该检查
	if !matchProtocol(r.protocol, p.Protocol()) {
		return false
	}
	// 匹配 源 IP
	if !matchIPNet(r.srcNet, p.SrcIP()) {
		return false
	}
	// 匹配 目标 IP
	if !matchIPNet(r.dstNet, p.DstIP()) {
		return false
	}
	// 匹配 端口（如果这个协议有端口） // 仅支持tcp udp
	if p.UsePort() {
		// 匹配 源端口
		if !matchPort(r.srcPort, p.SrcPort()) {
			return false
		}
		if !matchPort(r.dstPort, p.DstPort()) {
			return false
		}
	}
	return true
}

func (r *Rule) Accept() bool {
	return r.accept
}

func (r *Rule) Priority() int {
	return r.priority
}
