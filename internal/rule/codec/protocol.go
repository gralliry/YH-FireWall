package codec

import (
	"strings"
)

func ParseProtocol(raw string) ([]string, error) {
	parts := Split(raw)
	if len(parts) == 0 {
		return nil, nil
	}
	return parts, nil
}

func StringifyProtocol(ps []string) string {
	return strings.Join(ps, StandardSplitter)
}
