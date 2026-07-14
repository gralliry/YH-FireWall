package group

import (
	"slices"
)

type Container[V any] interface {
	Contains(V) bool
}

type Group[S Container[V], V any] struct {
	elems []Container[V]
}

func NewGroup[S Container[V], V any](elems []S) *Group[S, V] {
	containers := make([]Container[V], len(elems))
	for i := range elems {
		containers[i] = elems[i]
	}
	return &Group[S, V]{elems: containers}
}

func (g *Group[S, V]) Contains(v V) bool {
	for _, elem := range g.elems {
		if elem.Contains(v) {
			return true
		}
	}
	return false
}

func (g *Group[S, V]) Raw(v V) []Container[V] {
	return slices.Clone(g.elems)
}
