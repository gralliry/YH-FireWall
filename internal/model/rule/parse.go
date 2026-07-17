package rule

import (
	"YH-FireWall/internal/model/rule/codec"
	"YH-FireWall/internal/pkg/container"
	"fmt"
	"net/netip"
)

func (r *Rule) Update(o *Option, devMap DevName2Index) error {
	// Phase 1: 解析所有 codec 字段，全部通过才进入 Phase 2
	var (
		srcPrefixs    []netip.Prefix
		dstPrefixs    []netip.Prefix
		srcPortRanges [][2]uint16
		dstPortRanges [][2]uint16
		err           error
	)
	if o.SrcPrefixs != nil {
		srcPrefixs, err = codec.ParsePrefix(*o.SrcPrefixs)
		if err != nil {
			return fmt.Errorf("srcPrefixs: %w", err)
		}
	}
	if o.DstPrefixs != nil {
		dstPrefixs, err = codec.ParsePrefix(*o.DstPrefixs)
		if err != nil {
			return fmt.Errorf("dstPrefixs: %w", err)
		}
	}
	if o.SrcPortRanges != nil {
		srcPortRanges, err = codec.ParsePort(*o.SrcPortRanges)
		if err != nil {
			return fmt.Errorf("srcPortRanges: %w", err)
		}
	}
	if o.DstPortRanges != nil {
		dstPortRanges, err = codec.ParsePort(*o.DstPortRanges)
		if err != nil {
			return fmt.Errorf("dstPortRanges: %w", err)
		}
	}

	// Phase 2: 全部解析成功，统一写入
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
	if o.InDevs != nil {
		inDevs := codec.ParseDev(*o.InDevs, devMap)
		r.inDevs = container.NewSet(inDevs)
	}
	if o.OutDevs != nil {
		outDevs := codec.ParseDev(*o.OutDevs, devMap)
		r.outDevs = container.NewSet(outDevs)
	}
	if o.Protocols != nil {
		protocols := codec.ParseProtocol(*o.Protocols)
		r.protocols = container.NewSet(protocols)
	}
	if o.SrcPrefixs != nil {
		r.srcPrefixs = container.NewGroup(srcPrefixs)
	}
	if o.DstPrefixs != nil {
		r.dstPrefixs = container.NewGroup(dstPrefixs)
	}
	if o.SrcPortRanges != nil {
		r.srcPortRanges = container.NewRange(srcPortRanges)
	}
	if o.DstPortRanges != nil {
		r.dstPortRanges = container.NewRange(dstPortRanges)
	}

	return nil
}
