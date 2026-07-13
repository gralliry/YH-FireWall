package group

import (
	"slices"
	"sort"
)

type PortRange struct {
	raw    [][2]uint16
	ranges [][2]uint16
}

func NewPortRange(raw [][2]uint16) *PortRange {
	// [2]uint16 是值类型，可以复制的
	ranges := slices.Clone(raw)
	// 排序范围
	sorted := slices.Clone(ranges)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i][0] < sorted[j][0]
	})
	// 合并重复范围
	merged := make([][2]uint16, 0, len(sorted))
	cur := sorted[0]
	for i := 1; i < len(sorted); i++ {
		if sorted[i][0] <= cur[1] {
			if sorted[i][1] > cur[1] {
				cur[1] = sorted[i][1]
			}
		} else {
			merged = append(merged, cur)
			cur = sorted[i]
		}
	}
	merged = append(merged, cur)
	//
	return &PortRange{
		raw:    ranges,
		ranges: merged,
	}
}

func (p *PortRange) Contains(port uint16) bool {
	n := len(p.ranges)
	if n <= 8 {
		return slices.ContainsFunc(p.ranges, func(rg [2]uint16) bool {
			return rg[0] <= port && port <= rg[1]
		})
	} else {
		i := sort.Search(n, func(i int) bool {
			return port < p.ranges[i][0]
		})
		return i > 0 && port <= p.ranges[i-1][1]
	}
}

func (p *PortRange) Raw() [][2]uint16 {
	return p.raw
}
