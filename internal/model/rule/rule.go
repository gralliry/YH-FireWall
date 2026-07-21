package rule

import (
	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/model/rule/codec"
	"YH-FireWall/internal/pkg/container"
	"fmt"
)

func (r *Rule) ID() string {
	return r.id
}

func (r *Rule) Accept() bool {
	return r.accept
}

func (r *Rule) Compare(other *Rule) int {
	if r.enable == other.enable {
		if other.priority > r.priority {
			return 1
		} else if other.priority < r.priority {
			return -1
		} else {
			return 0
		}
	} else if r.enable {
		return -1
	} else {
		return 1
	}
}

func (r *Rule) Match(f *flow.Flow) bool {
	// 这里可以通过架构去优化，减少if次数
	if !r.enable {
		return false
	}
	// 匹配 入口网卡 // 为空默认跳过该检查
	if r.inDevs != nil && !r.inDevs.Contains(f.InDev) {
		return false
	}
	// 匹配 出口网卡 // 为空默认跳过该检查
	if r.outDevs != nil && !r.outDevs.Contains(f.OutDev) {
		return false
	}
	// 匹配 协议 // 为空默认跳过该检查
	if r.protocols != nil && !r.protocols.Contains(f.Protocol) {
		return false
	}
	// 匹配 源 IP
	if r.srcNets != nil && !r.srcNets.Contains(f.SrcIP) {
		return false
	}
	// 匹配 目标 IP
	if r.dstNets != nil && !r.dstNets.Contains(f.DstIP) {
		return false
	}
	// 匹配 端口（如果这个协议有端口）
	if r.srcPorts != nil || r.dstPorts != nil {
		if !f.HasPort {
			return false
		}
		// 源端口
		if r.srcPorts != nil && !r.srcPorts.Contains(f.SrcPort) {
			return false
		}
		// 目标端口
		if r.dstPorts != nil && !r.dstPorts.Contains(f.DstPort) {
			return false
		}
	}
	return true
}

func setter[V any](p *V, v V) V {
	if p != nil {
		return *p
	} else {
		return v
	}
}

func (r *Rule) Update(o *Option, devMap DevName2Index, protoMap Name2Protocol) (*Rule, error) {
	// Phase 1: 解析所有 codec 字段，全部通过才进入 Phase 2
	nr := &Rule{
		id:      r.id,
		group:   setter(o.Group, r.group),
		comment: setter(o.Comment, r.comment),

		accept:   setter(o.Accept, r.accept),
		priority: setter(o.Priority, r.priority),
		enable:   setter(o.Enable, r.enable),
	}
	// Phase 2: 全部解析成功，统一写入
	if o.SrcNets == nil {
		nr.srcNets = r.srcNets
	} else if srcNets, err := codec.ParsePrefix(*o.SrcNets); err != nil {
		return nil, fmt.Errorf("srcNets: %w", err)
	} else {
		nr.srcNets = container.NewGroup(srcNets)
	}

	if o.DstNets == nil {
		nr.dstNets = r.dstNets
	} else if dstNets, err := codec.ParsePrefix(*o.DstNets); err != nil {
		return nil, fmt.Errorf("dstNets: %w", err)
	} else {
		nr.dstNets = container.NewGroup(dstNets)
	}

	if o.SrcPorts == nil {
		nr.srcPorts = r.srcPorts
	} else if srcPorts, err := codec.ParsePort(*o.SrcPorts); err != nil {
		return nil, fmt.Errorf("srcPorts: %w", err)
	} else {
		nr.srcPorts = container.NewRange(srcPorts)
	}

	if o.DstPorts == nil {
		nr.dstPorts = r.dstPorts
	} else if dstPorts, err := codec.ParsePort(*o.DstPorts); err != nil {
		return nil, fmt.Errorf("dstPorts: %w", err)
	} else {
		nr.dstPorts = container.NewRange(dstPorts)
	}

	if o.InDevs != nil {
		inDevs := codec.ParseDev(*o.InDevs, devMap)
		nr.inDevs = container.NewSet(inDevs)
	} else {
		nr.inDevs = r.inDevs
	}

	// 设备是有可能预存在的
	if o.OutDevs != nil {
		outDevs := codec.ParseDev(*o.OutDevs, devMap)
		nr.outDevs = container.NewSet(outDevs)
	} else {
		nr.outDevs = r.outDevs
	}

	if o.Protocols == nil {
		nr.protocols = r.protocols
	} else if protocols, err := codec.ParseProtocol(*o.Protocols, protoMap); err != nil {
		return nil, fmt.Errorf("protocol: %w", err)
	} else {
		nr.protocols = container.NewSet(protocols)
	}
	return nr, nil
}
