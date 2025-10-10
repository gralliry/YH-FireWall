package rule

import (
	"YH-FireWall/core/packet"
	"fmt"
	"github.com/google/gopacket/layers"
	"net"
	"sync"
)

type Rule struct {
	id       string
	group    string
	comment  string
	accept   bool
	priority int
	enable   bool

	// 动态解析
	srcNets  []net.IPNet
	srcPorts [][2]uint16

	dstNets  []net.IPNet
	dstPorts [][2]uint16

	inDevs  map[uint32]struct{}
	outDevs map[uint32]struct{}

	protocols map[layers.IPProtocol]struct{}

	//
	mutex sync.RWMutex
}

func Parse(cfg Config) (*Rule, error) {
	if cfg.Id == "" {
		return nil, fmt.Errorf("rule id is empty")
	}
	r := &Rule{
		id:       cfg.Id,
		group:    cfg.Group,
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
	// 源端口
	if srcPort, err := parsePort(cfg.SrcPort); err == nil {
		r.srcPorts = srcPort
	} else {
		return nil, err
	}
	// 目标网络
	if tarNet, err := parseIPNet(cfg.TarNet); err == nil {
		r.dstNets = tarNet
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
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return &Config{
		Id:       r.id,
		Group:    r.group,
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
	r.mutex.RLock()
	defer r.mutex.RUnlock()
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

func (r *Rule) Update(o Option) (err error) {
	var srcNet []net.IPNet
	if o.SrcNet != nil {
		if srcNet, err = parseIPNet(*o.SrcNet); err != nil {
			return err
		}
	}
	var srcPort [][2]uint16
	if o.SrcPort != nil {
		if srcPort, err = parsePort(*o.SrcPort); err != nil {
			return err
		}
	}
	var tarNet []net.IPNet
	if o.TarNet != nil {
		if tarNet, err = parseIPNet(*o.TarNet); err != nil {
			return err
		}
	}
	var tarPort [][2]uint16
	if o.TarPort != nil {
		if tarPort, err = parsePort(*o.TarPort); err != nil {
			return err
		}
	}
	var inDevs map[uint32]struct{}
	if o.InDev != nil {
		if inDevs, err = parseDev(*o.InDev); err != nil {
			return err
		}
	}
	var outDevs map[uint32]struct{}
	if o.OutDev != nil {
		if outDevs, err = parseDev(*o.OutDev); err != nil {
			return err
		}
	}
	var protocols map[layers.IPProtocol]struct{}
	if o.Protocol != nil {
		if protocols, err = parseProtocol(*o.Protocol); err != nil {
			return err
		}
	}
	// 锁
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// 到此为止，没有错误，开始更新
	if o.Group != nil {
		r.group = *o.Group
	}
	if o.Comment != nil {
		r.comment = *o.Comment
	}
	if o.SrcNet != nil {
		r.srcNets = srcNet
	}
	if o.SrcPort != nil {
		r.srcPorts = srcPort
	}
	if o.TarNet != nil {
		r.dstNets = tarNet
	}
	if o.TarPort != nil {
		r.dstPorts = tarPort
	}
	if o.InDev != nil {
		r.inDevs = inDevs
	}
	if o.OutDev != nil {
		r.outDevs = outDevs
	}
	if o.Protocol != nil {
		r.protocols = protocols
	}
	if o.Accept != nil {
		r.accept = *o.Accept
	}
	if o.Priority != nil {
		r.priority = *o.Priority
	}
	if o.Enable != nil {
		r.enable = *o.Enable
	}
	return nil
}

func (r *Rule) Id() string {
	// Id 一定不会发生修改，无需加锁
	return r.id
}

func (r *Rule) Group() string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.group
}

func (r *Rule) Priority() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.priority
}

func (r *Rule) SetEnable(enable bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.enable = enable
}

func (r *Rule) Enable() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.enable
}

func (r *Rule) Accept() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.accept
}
