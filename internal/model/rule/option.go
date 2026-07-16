package rule

import (
	"YH-FireWall/internal/model/rule/codec"
	"YH-FireWall/internal/pkg/container"
	"YH-FireWall/internal/pkg/funcs"
	"fmt"
	"net/netip"
	"strings"
)

// Option 是统一的输入/输出 DTO。
// 输入时 *string/*bool/*int 为 nil 表示不更新该字段；
// 输出时所有字段均指向有效值。
type Option struct {
	ID            string  `json:"id"`
	Group         *string `json:"group,omitempty"`
	Comment       *string `json:"comment,omitempty"`
	Accept        *bool   `json:"accept,omitempty"`
	Priority      *int    `json:"priority,omitempty"`
	Enable        *bool   `json:"enable,omitempty"`
	SrcPrefixs    *string `json:"srcPrefixs,omitempty"`
	DstPrefixs    *string `json:"dstPrefixs,omitempty"`
	SrcPortRanges *string `json:"srcPortRanges,omitempty"`
	DstPortRanges *string `json:"dstPortRanges,omitempty"`
	InDevs        *string `json:"inDevs,omitempty"`
	OutDevs       *string `json:"outDevs,omitempty"`
	Protocols     *string `json:"protocols,omitempty"`
}

func (o *Option) String() string {
	var sb strings.Builder

	symbol := '✓'
	if o.Enable != nil && !*o.Enable {
		symbol = '✗'
	}

	action := "DENY"
	if o.Accept != nil && *o.Accept {
		action = "ACCEPT"
	}

	fmt.Fprintf(&sb, "%-8s  %c  %s", o.ID, symbol, action)
	if o.Priority != nil {
		fmt.Fprintf(&sb, "  [%d]", *o.Priority)
	}
	fmt.Fprintln(&sb)

	if o.Group != nil {
		fmt.Fprintf(&sb, "  group:  %s\n", *o.Group)
	}
	if o.Comment != nil {
		fmt.Fprintf(&sb, "  note:   %s\n", *o.Comment)
	}
	if o.SrcPrefixs != nil {
		fmt.Fprintf(&sb, "  src:    %s\n", *o.SrcPrefixs)
	}
	if o.DstPrefixs != nil {
		fmt.Fprintf(&sb, "  dst:    %s\n", *o.DstPrefixs)
	}
	if o.SrcPortRanges != nil {
		fmt.Fprintf(&sb, "  sport:  %s\n", *o.SrcPortRanges)
	}
	if o.DstPortRanges != nil {
		fmt.Fprintf(&sb, "  dport:  %s\n", *o.DstPortRanges)
	}
	if o.Protocols != nil {
		fmt.Fprintf(&sb, "  proto:  %s\n", *o.Protocols)
	}
	if o.InDevs != nil {
		fmt.Fprintf(&sb, "  in:     %s\n", *o.InDevs)
	}
	if o.OutDevs != nil {
		fmt.Fprintf(&sb, "  out:    %s\n", *o.OutDevs)
	}

	return sb.String()
}

// Option 以 codec 字符串形式返回当前规则的完整快照。
// devMap 和 proMap 用于将内部数值映射为可读名称。
func (r *Rule) Option(devMap map[uint32]string) *Option {
	return &Option{
		ID:            r.id,
		Group:         new(r.group),
		Comment:       new(r.comment),
		Accept:        new(r.accept),
		Priority:      new(r.priority),
		Enable:        new(r.enable),
		SrcPrefixs:    new(codec.StringifyPrefix(r.srcPrefixs.Raw())),
		DstPrefixs:    new(codec.StringifyPrefix(r.dstPrefixs.Raw())),
		SrcPortRanges: new(codec.StringifyPort(r.srcPortRanges.Raw())),
		DstPortRanges: new(codec.StringifyPort(r.dstPortRanges.Raw())),
		InDevs:        new(codec.StringifyDev(funcs.Collect(r.inDevs.Raw(), devMap))),
		OutDevs:       new(codec.StringifyDev(funcs.Collect(r.outDevs.Raw(), devMap))),
		Protocols:     new(codec.StringifyProtocol(r.protocols.Raw())),
	}
}

func (r *Rule) Update(o *Option, devMap map[string]uint32) error {
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
		inDevs := codec.ParseDev(*o.InDevs)
		r.inDevs = container.NewSet(inDevs)
	}
	if o.OutDevs != nil {
		outDevs := codec.ParseDev(*o.OutDevs)
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
