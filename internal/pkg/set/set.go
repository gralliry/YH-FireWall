package set

import (
	"fmt"
	"slices"
	"strings"
)

type Set[T comparable] struct {
	raw []T
	s   []T
	m   map[T]struct{}
}

func New[T comparable](vals []T) *Set[T] {
	raw := slices.Clone(vals)
	m := make(map[T]struct{}, len(raw))
	for _, v := range raw {
		m[v] = struct{}{}
	}
	if len(m) <= threshold[T]() {
		out := make([]T, 0, len(m))
		for v := range m {
			out = append(out, v)
		}
		return &Set[T]{raw: raw, s: out}
	}
	return &Set[T]{raw: raw, m: m}
}

func (s *Set[T]) Contains(v T) bool {
	if s.s != nil {
		return slices.Contains(s.s, v)
	}
	_, ok := s.m[v]
	return ok
}

func (s *Set[T]) Raw() []T {
	return s.raw
}

func (s *Set[T]) String() string {
	var parts []string
	for _, v := range s.raw {
		parts = append(parts, fmt.Sprintf("%v", v))
	}
	return strings.Join(parts, ",")
}

func threshold[T comparable]() int {
	var zero T
	switch any(zero).(type) {
	case string:
		return 8
	default:
		return 32
	}
}
