package rule

import (
	"YH-FireWall/internal/model/rule/codec"
	"fmt"
	"strings"
)

// Info 以 codec 字符串形式返回当前规则的完整快照。
// devMap 用于将内部数值映射为可读名称。
func (r *Rule) Data(devMap DevIndex2Name, protoMap Protocol2Name) *Data {
	d := &Data{
		ID: r.id,
		Option: Option{
			Group:    new(r.group),
			Comment:  new(r.comment),
			Accept:   new(r.accept),
			Priority: new(r.priority),
			Enable:   new(r.enable),
		},
	}
	if r.srcNets != nil {
		d.Option.SrcNets = new(codec.StringifyPrefix(r.srcNets.Raw()))
	}
	if r.dstNets != nil {
		d.Option.DstNets = new(codec.StringifyPrefix(r.dstNets.Raw()))
	}
	if r.srcPorts != nil {
		d.Option.SrcPorts = new(codec.StringifyPort(r.srcPorts.Raw()))
	}
	if r.dstPorts != nil {
		d.Option.DstPorts = new(codec.StringifyPort(r.dstPorts.Raw()))
	}
	if r.inDevs != nil {
		d.Option.InDevs = new(codec.StringifyDev(r.inDevs.Raw(), devMap))
	}
	if r.outDevs != nil {
		d.Option.OutDevs = new(codec.StringifyDev(r.outDevs.Raw(), devMap))
	}
	if r.protocols != nil {
		d.Option.Protocols = new(codec.StringifyProtocol(r.protocols.Raw(), protoMap))
	}
	return d
}

func val[V any](p *V) (v V) {
	if p != nil {
		return *p
	}
	return
}

func (d *Data) String() string {
	var sb strings.Builder
	symbol := '✗'
	if d.Enable != nil && *d.Enable {
		symbol = '✓'
	}
	action := "DENY"
	if d.Accept != nil && *d.Accept {
		action = "ACCEPT"
	}

	fmt.Fprintf(&sb, "%-8s  %c  %s\n", d.ID, symbol, action)
	fmt.Fprintf(&sb, "  [%d]\n", d.Priority)
	fmt.Fprintf(&sb, "  group:   %s\n", val(d.Group))
	fmt.Fprintf(&sb, "  comment: %s\n", val(d.Comment))
	fmt.Fprintf(&sb, "  proto:   %s\n", val(d.Protocols))
	fmt.Fprintf(&sb, "  src:     %s\n", val(d.SrcNets))
	fmt.Fprintf(&sb, "  dst:     %s\n", val(d.DstNets))
	fmt.Fprintf(&sb, "  sport:   %s\n", val(d.SrcPorts))
	fmt.Fprintf(&sb, "  dport:   %s\n", val(d.DstPorts))
	fmt.Fprintf(&sb, "  in:      %s\n", val(d.InDevs))
	fmt.Fprintf(&sb, "  out:     %s\n", val(d.OutDevs))

	return sb.String()
}
