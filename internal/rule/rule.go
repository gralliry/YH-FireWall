package rule

import (
	"YH-FireWall/internal/pkg/sid"
	"YH-FireWall/internal/rule/group"
	"net/netip"

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
	srcPrefixs, dstPrefixs       *group.Group[netip.Prefix, netip.Addr]
	srcPortRanges, dstPortRanges *group.Range[uint16]
	inDevs, outDevs              *group.Set[uint32]
	protocols                    *group.Set[layers.IPProtocol]
}

func New() *Rule {
	return &Rule{id: sid.New(8)}
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
