package codec

import (
	"fmt"
	"strconv"
	"strings"
)

func ParsePort(raw string) ([][2]uint16, error) {
	parts := split(raw)
	if len(parts) == 0 {
		return nil, nil
	}
	ranges := make([][2]uint16, 0, len(parts))
	for _, part := range parts {
		if strings.Contains(part, "-") {
			se := strings.SplitN(part, "-", 2)
			if len(se) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}
			start, err := strconv.Atoi(strings.TrimSpace(se[0]))
			if err != nil || start < 0 || start > 65535 {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			end, err := strconv.Atoi(strings.TrimSpace(se[1]))
			if err != nil || end < 0 || end > 65535 {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			if start > end {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}
			ranges = append(ranges, [2]uint16{uint16(start), uint16(end)})
		} else {
			val, err := strconv.Atoi(part)
			if err != nil || val < 0 || val > 65535 {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			ranges = append(ranges, [2]uint16{uint16(val), uint16(val)})
		}
	}
	return ranges, nil
}

func StringifyPort(ranges [][2]uint16) string {
	var parts []string
	for _, p := range ranges {
		if p[0] == p[1] {
			parts = append(parts, fmt.Sprintf("%d", p[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d-%d", p[0], p[1]))
		}
	}
	return strings.Join(parts, StandardSplitter)
}

// ====================================================
//
