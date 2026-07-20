package rule

import (
	"YH-FireWall/internal/model/rule/codec"
	"YH-FireWall/internal/pkg/container"
	"YH-FireWall/internal/pkg/sid"
	"net/netip"

	"github.com/google/gopacket/layers"
)

// Option 是统一的输入。
// 输入时 *string/*bool/*int 为 nil 表示不更新该字段；
// 输出时所有字段均指向有效值。
type Option struct {
	Group     *string `json:"group,omitempty"`
	Comment   *string `json:"comment,omitempty"`
	Accept    *bool   `json:"accept,omitempty"`
	Priority  *int    `json:"priority,omitempty"`
	Enable    *bool   `json:"enable,omitempty"`
	SrcNets   *string `json:"srcNets,omitempty"`
	DstNets   *string `json:"dstNets,omitempty"`
	SrcPorts  *string `json:"srcPorts,omitempty"`
	DstPorts  *string `json:"dstPorts,omitempty"`
	InDevs    *string `json:"inDevs,omitempty"`
	OutDevs   *string `json:"outDevs,omitempty"`
	Protocols *string `json:"protocols,omitempty"`
}

type Data struct {
	ID string `json:"id"`
	Option
}



type Rule struct {
	id       string
	group    string
	comment  string
	accept   bool
	priority int
	enable   bool

	// 动态解析
	srcNets, dstNets   *container.Group[netip.Prefix, netip.Addr]
	srcPorts, dstPorts *container.Range[uint16]
	inDevs, outDevs    *container.Set[uint32]
	protocols          *container.Set[layers.IPProtocol]
}

type (
	DevName2Index = codec.DevName2Index
	DevIndex2Name = codec.DevIndex2Name
	//
	Name2Protocol = codec.Name2Protocol
	Protocol2Name = codec.Protocol2Name
)

func New(o *Option, devMap DevName2Index, protoMap Name2Protocol) (*Rule, error) {
	r := &Rule{id: sid.New(12)}
	if err := r.Update(o, devMap, protoMap); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func Parse(d *Data, devMap DevName2Index, protoMap Name2Protocol) (*Rule, error) {
	r := &Rule{id: d.ID}
	if err := r.Update(&d.Option, devMap, protoMap); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}
