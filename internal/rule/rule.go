package rule

import (
	"YH-FireWall/internal/rule/group"
	"fmt"

	"github.com/google/gopacket/layers"
)

type Rule struct {
	id       string
	group    string
	comment  string
	accept   bool
	priority int
	enable   bool

	// 动态解析
	srcPrefixs *group.PrefixGroup
	dstPrefixs *group.PrefixGroup

	srcPorts *group.PortRange
	dstPorts *group.PortRange

	inDevs  *group.Set[uint32]
	outDevs *group.Set[uint32]

	protocols *group.Set[layers.IPProtocol]
}

func New(cfg *Info) (*Rule, error) {
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
	if r.srcPrefixs, err = parsePrefix(cfg.SrcNet); err != nil {
		return nil, fmt.Errorf("parse source network failed: %w", err)
	}
	// 源端口
	if r.srcPorts, err = parsePort(cfg.SrcPort); err != nil {
		return nil, fmt.Errorf("parse source port failed: %w", err)
	}
	// 目标网络
	if r.dstPrefixs, err = parseIPNet(cfg.TarNet); err != nil {
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

func (r *Rule) Info() *Info {
	return &Info{
		Id:       r.id,
		Group:    r.group,
		Comment:  r.comment,
		SrcNet:   stringifyIPNet(r.srcPrefixs.Raw()),
		SrcPort:  stringifyPort(r.srcPorts.Raw()),
		TarNet:   stringifyIPNet(r.dstPrefixs.Raw()),
		TarPort:  stringifyPort(r.dstPorts.Raw()),
		InDev:    stringifyDev(r.inDevs.Raw()),
		OutDev:   stringifyDev(r.outDevs.Raw()),
		Protocol: stringifyProtocol(r.protocols.Raw()),
		Accept:   r.accept,
		Priority: r.priority,
		Enable:   r.enable,
	}
}

func (r *Rule) Update(o Option) error {
	if o.Group != nil {
		r.group = *o.Group
	}
	if o.Comment != nil {
		r.comment = *o.Comment
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

	var err error
	if o.SrcNet != nil {
		if r.srcPrefixs, err = parseIPNet(*o.SrcNet); err != nil {
			return fmt.Errorf("parse source network failed: %w", err)
		}
	}
	if o.SrcPort != nil {
		if r.srcPorts, err = parsePort(*o.SrcPort); err != nil {
			return fmt.Errorf("parse source port failed: %w", err)
		}
	}
	if o.TarNet != nil {
		if r.dstPrefixs, err = parseIPNet(*o.TarNet); err != nil {
			return fmt.Errorf("parse target network failed: %w", err)
		}
	}
	if o.TarPort != nil {
		if r.dstPorts, err = parsePort(*o.TarPort); err != nil {
			return fmt.Errorf("parse target port failed: %w", err)
		}
	}
	if o.InDev != nil {
		if r.inDevs, err = parseDev(*o.InDev); err != nil {
			return fmt.Errorf("parse input device failed: %w", err)
		}
	}
	if o.OutDev != nil {
		if r.outDevs, err = parseDev(*o.OutDev); err != nil {
			return fmt.Errorf("parse output device failed: %w", err)
		}
	}
	if o.Protocol != nil {
		if r.protocols, err = parseProtocol(*o.Protocol); err != nil {
			return fmt.Errorf("parse protocol failed: %w", err)
		}
	}
	return nil
}

func (r *Rule) Id() string {
	// Index 一定不会发生修改，无需加锁
	return r.id
}

func (r *Rule) Group() string {
	return r.group
}

func (r *Rule) Priority() int {
	return r.priority
}

func (r *Rule) SetEnable(enable bool) {
	r.enable = enable
}

func (r *Rule) Accept() bool {
	return r.accept
}
