package group

import (
	"net/netip"
	"slices"
)

type PrefixGroup struct {
	prefixs []netip.Prefix
}

func NewPrefixGroup(prefixs []netip.Prefix) *PrefixGroup {
	return &PrefixGroup{prefixs: slices.Clone(prefixs)}
}

func (p *PrefixGroup) Contains(ip netip.Addr) bool {
	for _, prefix := range p.prefixs {
		if prefix.Contains(ip) {
			return true
		}
	}
	return false
}

func (p *PrefixGroup) Raw() []netip.Prefix {
	return p.prefixs
}
