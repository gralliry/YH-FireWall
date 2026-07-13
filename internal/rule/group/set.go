package group

import (
	"slices"
)

type Set[T comparable] struct {
	raw []T
	s   []T
	m   map[T]struct{}
}

// New creates a set containing the provided values (duplicates removed).
// For string the threshold is 8; for all other types it is 32.
func New[T comparable](vals []T) *Set[T] {
	raw := slices.Clone(vals)
	m := make(map[T]struct{}, len(raw))
	for _, v := range raw {
		m[v] = struct{}{}
	}
	if len(m) <= Threshold[T]() {
		out := make([]T, 0, len(m))
		for v := range m {
			out = append(out, v)
		}
		return &Set[T]{s: out}
	} else {
		return &Set[T]{m: m}
	}
}

// Contains reports whether v is in the set.
func (s *Set[T]) Contains(v T) bool {
	if s.s != nil {
		return slices.Contains(s.s, v)
	} else {
		_, ok := s.m[v]
		return ok
	}
}

// String implements fmt.Stringer.
func (s *Set[T]) Raw() []T {
	return s.raw
}

// thresholdFor returns the slice-vs-map cutoff for the element type.
func Threshold[T comparable]() int {
	var zero T
	switch any(zero).(type) {
	case string:
		return 8
	default:
		return 32
	}
}
