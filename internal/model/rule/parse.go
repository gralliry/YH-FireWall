package rule

import (
	"YH-FireWall/internal/model/rule/codec"
	"YH-FireWall/internal/pkg/container"
	"fmt"
	"net/netip"
	"strconv"
	"strings"
)

func (r *Rule) Update(o *Option, devMap DevName2Index, protoMap Name2Protocol) error {
	// Phase 1: 解析所有 codec 字段，全部通过才进入 Phase 2
	var (
		srcNets    []netip.Prefix
		dstNets    []netip.Prefix
		srcPorts [][2]uint16
		dstPorts [][2]uint16
		err           error
	)
	if o.SrcNets != nil {
		if srcNets, err = codec.ParsePrefix(*o.SrcNets); err != nil {
			return fmt.Errorf("srcNets: %w", err)
		}
	}
	if o.DstNets != nil {
		if dstNets, err = codec.ParsePrefix(*o.DstNets); err != nil {
			return fmt.Errorf("dstNets: %w", err)
		}
	}
	if o.SrcPorts != nil {
		if srcPorts, err = codec.ParsePort(*o.SrcPorts); err != nil {
			return fmt.Errorf("srcPorts: %w", err)
		}
	}
	if o.DstPorts != nil {
		if dstPorts, err = codec.ParsePort(*o.DstPorts); err != nil {
			return fmt.Errorf("dstPorts: %w", err)
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
		protocols := codec.ParseProtocol(*o.Protocols, protoMap)
		r.protocols = container.NewSet(protocols)
	}
	if o.SrcNets != nil {
		r.srcNets = container.NewGroup(srcNets)
	}
	if o.DstNets != nil {
		r.dstNets = container.NewGroup(dstNets)
	}
	if o.SrcPorts != nil {
		r.srcPorts = container.NewRange(srcPorts)
	}
	if o.DstPorts != nil {
		r.dstPorts = container.NewRange(dstPorts)
	}

	return nil
}

func Set(key, value string) (*Option, error) {
	o := &Option{}
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "_", "")
	switch key {
	case "group":
		o.Group = &value
	case "comment":
		o.Comment = &value
	case "accept":
		switch strings.ToLower(value) {
		case "true":
			o.Accept = new(true)
		case "false":
			o.Accept = new(false)
		default:
			return nil, fmt.Errorf("invalid bool value: %s", value)
		}
	case "enable":
		switch strings.ToLower(value) {
		case "true":
			o.Enable = new(true)
		case "false":
			o.Enable = new(false)
		default:
			return nil, fmt.Errorf("invalid bool value: %s", value)
		}
	case "priority":
		n, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid int value: %s", value)
		}
		o.Priority = &n
	case "indevs":
		o.InDevs = &value
	case "outdevs":
		o.OutDevs = &value
	case "protocols":
		o.Protocols = &value
	case "srcnets":
		o.SrcNets = &value
	case "dstnets":
		o.DstNets = &value
	case "srcports":
		o.SrcPorts = &value
	case "dstports":
		o.DstPorts = &value
	default:
		return nil, fmt.Errorf("unknown key: %s", key)
	}
	return o, nil
}
