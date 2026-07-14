package group

import (
	"slices"
)

type Set[V comparable] struct {
	raw []V
	s   []V
	m   map[V]struct{}
}

func NewSet[V comparable](devs []V) *Set[V] {
	raw := slices.Clone(devs)
	m := make(map[V]struct{}, len(devs))
	for _, dev := range devs {
		m[dev] = struct{}{}
	}
	if len(devs) <= 32 {
		s := make([]V, 0, len(devs))
		for v := range m {
			s = append(s, v)
		}
		return &Set[V]{raw: raw, s: s}
	} else {
		return &Set[V]{raw: raw, m: m}
	}
}

func (s *Set[V]) Contains(v V) bool {
	if s.s != nil {
		return slices.Contains(s.s, v)
	} else {
		_, ok := s.m[v]
		return ok
	}
}

func (s *Set[V]) Raw() []V {
	return slices.Clone(s.raw)
}
