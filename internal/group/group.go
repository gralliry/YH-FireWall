package group

import (
	"YH-FireWall/internal/packet"
	"YH-FireWall/internal/rule"
	"sort"
)

type Group struct {
	id    uint16
	rules []*rule.Rule
}

func (g *Group) Match(p *packet.Packet) (bool, bool) {
	for _, r := range g.rules {
		if r.Match(p) {
			return true, r.Accept()
		}
	}
	return false, false
}

func (g *Group) Register(r *rule.Rule) {
	g.rules = append(g.rules, r)
}

func (g *Group) Fresh() {
	sort.Slice(g.rules, func(i, j int) bool {
		return g.rules[i].Priority() > g.rules[j].Priority()
	})
}
