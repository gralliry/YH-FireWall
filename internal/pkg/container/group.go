package container

import (
	"slices"
)

type Container[V any] interface {
	Contains(V) bool
}

type Group[S Container[V], V any] struct {
	elems []S
}

func NewGroup[S Container[V], V any](elems []S) *Group[S, V] {
	containers := slices.Clone(elems)
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

func (g *Group[S, V]) Raw() []S {
	return slices.Clone(g.elems)
}
