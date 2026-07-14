package rule

import (
	"YH-FireWall/internal/rule/group"
	"net/netip"

	"github.com/google/gopacket/layers"
	"github.com/samber/lo"
)

type Option struct {
	Group         *string
	Comment       *string
	SrcPrefixs    []netip.Prefix
	DstPrefixs    []netip.Prefix
	SrcPortRanges [][2]uint16
	DstPortRanges [][2]uint16
	InDevs        []string
	OutDevs       []string
	Protocols     []layers.IPProtocol
	Accept        *bool
	Priority      *int
	Enable        *bool
}

func (r *Rule) Update(o *Option, devMap map[string]uint32) error {
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
	if o.SrcPrefixs != nil {
		r.srcPrefixs = group.NewGroup[netip.Prefix, netip.Addr](o.SrcPrefixs)
	}
	if o.SrcPortRanges != nil {
		r.srcPortRanges = group.NewRange(o.SrcPortRanges)
	}
	if o.DstPrefixs != nil {
		r.dstPrefixs = group.NewGroup[netip.Prefix, netip.Addr](o.DstPrefixs)
	}
	if o.DstPortRanges != nil {
		r.dstPortRanges = group.NewRange(o.DstPortRanges)
	}
	if o.InDevs != nil {
		list := lo.Map(o.InDevs, func(s string, _ int) uint32 { return devMap[s] })
		r.inDevs = group.NewSet(list)
	}
	if o.OutDevs != nil {
		list := lo.Map(o.OutDevs, func(s string, _ int) uint32 { return devMap[s] })
		r.outDevs = group.NewSet(list)
	}
	if o.Protocols != nil {
		r.protocols = group.NewSet(o.Protocols)
	}
	return nil
}
