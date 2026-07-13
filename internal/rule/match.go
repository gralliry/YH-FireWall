package rule

import (
	"YH-FireWall/internal/flow"
)

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
