package rule

import (
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

func Parse(cfg *Config) (*Rule, error) {
	if cfg == nil {
		return nil, fmt.Errorf("rule config is nil")
	}
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
	var err error
	if r.srcNets, err = parseIPNet(cfg.SrcNet); err != nil {
		return nil, fmt.Errorf("parse source network failed: %w", err)
	}
	// 源端口
	if r.srcPorts, err = parsePort(cfg.SrcPort); err != nil {
		return nil, fmt.Errorf("parse source port failed: %w", err)
	}
	// 目标网络
	if r.dstNets, err = parseIPNet(cfg.TarNet); err != nil {
		return nil, fmt.Errorf("parse target network failed: %w", err)
	}
	// 目标端口
	if r.dstPorts, err = parsePort(cfg.TarPort); err != nil {
		return nil, fmt.Errorf("parse target port failed: %w", err)
	}
	// 入口
	if r.inDevs, err = parseDev(cfg.InDev); err != nil {
		return nil, fmt.Errorf("parse input device failed: %w", err)
	}
	// 出口
	if r.outDevs, err = parseDev(cfg.OutDev); err != nil {
		return nil, fmt.Errorf("parse output device failed: %w", err)
	}
	// 协议
	if r.protocols, err = parseProtocol(cfg.Protocol); err != nil {
		return nil, fmt.Errorf("parse protocol failed: %w", err)
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

func (r *Rule) Match(
	srcIP net.IP,
	srcPort uint16,
	dstIP net.IP,
	dstPort uint16,
	inDev *uint32,
	outDev *uint32,
	protocol layers.IPProtocol,
) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	// 这里可以通过架构去优化，减少if次数
	if !r.enable {
		return false
	}
	// 匹配 入口网卡 // 为空默认跳过该检查
	if !matchDev(r.inDevs, inDev) {
		return false
	}
	// 匹配 出口网卡 // 为空默认跳过该检查
	if !matchDev(r.outDevs, outDev) {
		return false
	}
	// 匹配 协议 // 为空默认跳过该检查
	if !matchProtocol(r.protocols, protocol) {
		return false
	}
	// 匹配 源 IP
	if !matchIPNet(r.srcNets, srcIP) {
		return false
	}
	// 匹配 目标 IP
	if !matchIPNet(r.dstNets, dstIP) {
		return false
	}
	// 匹配 端口（如果这个协议有端口） // 仅支持tcp udp // sctp dccp 不支持
	if protocol == layers.IPProtocolTCP || protocol == layers.IPProtocolUDP {
		// 匹配 源端口
		if !matchPort(r.srcPorts, srcPort) {
			return false
		}
		// 匹配 目标端口
		if !matchPort(r.dstPorts, dstPort) {
			return false
		}
	}
	return true
}

func (r *Rule) Update(o Option) (err error) {
	// 预先解析所有需要更新的字段，如果有错误则直接返回
	var (
		srcNet    []net.IPNet
		srcPort   [][2]uint16
		tarNet    []net.IPNet
		tarPort   [][2]uint16
		inDevs    map[uint32]struct{}
		outDevs   map[uint32]struct{}
		protocols map[layers.IPProtocol]struct{}
	)
	// 解析源网络
	if o.SrcNet != nil {
		if srcNet, err = parseIPNet(*o.SrcNet); err != nil {
			return fmt.Errorf("parse source network failed: %w", err)
		}
	}
	// 解析源端口
	if o.SrcPort != nil {
		if srcPort, err = parsePort(*o.SrcPort); err != nil {
			return fmt.Errorf("parse source port failed: %w", err)
		}
	}
	// 解析目标网络
	if o.TarNet != nil {
		if tarNet, err = parseIPNet(*o.TarNet); err != nil {
			return fmt.Errorf("parse target network failed: %w", err)
		}
	}
	// 解析目标端口
	if o.TarPort != nil {
		if tarPort, err = parsePort(*o.TarPort); err != nil {
			return fmt.Errorf("parse target port failed: %w", err)
		}
	}
	// 解析入口设备
	if o.InDev != nil {
		if inDevs, err = parseDev(*o.InDev); err != nil {
			return fmt.Errorf("parse input device failed: %w", err)
		}
	}
	// 解析出口设备
	if o.OutDev != nil {
		if outDevs, err = parseDev(*o.OutDev); err != nil {
			return fmt.Errorf("parse output device failed: %w", err)
		}
	}
	// 解析协议
	if o.Protocol != nil {
		if protocols, err = parseProtocol(*o.Protocol); err != nil {
			return fmt.Errorf("parse protocol failed: %w", err)
		}
	}
	// 所有解析成功，加锁更新
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 更新各个字段
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

func (r *Rule) Accept() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.accept
}
