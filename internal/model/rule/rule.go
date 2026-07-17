package rule

import (
	"YH-FireWall/internal/model/flow"
)



func (r *Rule) ID() string {
	return r.id
}

func (r *Rule) Accept() bool {
	return r.accept
}

func (r *Rule) Compare(other *Rule) int {
	if r.enable != other.enable {
		if r.enable {
			return -1
		} else {
			return 1
		}
	}
	return other.priority - r.priority
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
	if r.srcPrefixs != nil && !r.srcPrefixs.Contains(f.SrcIP) {
		return false
	}
	// 匹配 目标 IP
	if r.dstPrefixs != nil && !r.dstPrefixs.Contains(f.DstIP) {
		return false
	}
	// 匹配 端口（如果这个协议有端口）
	if r.srcPortRanges != nil || r.dstPortRanges != nil {
		if !f.HasPort {
			return false
		}
		// 源端口
		if r.srcPortRanges != nil && !r.srcPortRanges.Contains(f.SrcPort) {
			return false
		}
		// 目标端口
		if r.dstPortRanges != nil && !r.dstPortRanges.Contains(f.DstPort) {
			return false
		}
	}
	return true
}
