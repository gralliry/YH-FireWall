package container

import (
	"cmp"
	"slices"
	"sort"
)

type Range[T cmp.Ordered] struct {
	raw    [][2]T
	ranges [][2]T
}

func NewRange[T cmp.Ordered](raw [][2]T) *Range[T] {
	raw = slices.Clone(raw)
	// 排序
	ranges := slices.Clone(raw)
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][0] < ranges[j][0]
	})
	// 合并
	merged := make([][2]T, 0, len(ranges))
	cur := ranges[0]
	for i := 1; i < len(ranges); i++ {
		if ranges[i][0] <= cur[1] {
			if ranges[i][1] > cur[1] {
				cur[1] = ranges[i][1]
			}
		} else {
			merged = append(merged, cur)
			cur = ranges[i]
		}
	}
	merged = append(merged, cur)
	return &Range[T]{raw: raw, ranges: merged}
}

func (r *Range[T]) Contains(port T) bool {
	if n := len(r.ranges); n <= 8 {
		return slices.ContainsFunc(r.ranges, func(rg [2]T) bool {
			return rg[0] <= port && port <= rg[1]
		})
	} else {
		i := sort.Search(n, func(i int) bool {
			return port < r.ranges[i][0]
		})
		return i > 0 && port <= r.ranges[i-1][1]
	}
}

func (r *Range[T]) Raw() [][2]T {
	return slices.Clone(r.raw)
}
