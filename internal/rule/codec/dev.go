package codec

import "strings"

func ParseDev(raw string) ([]string, error) {
	parts := Split(raw)
	if len(parts) == 0 {
		return nil, nil
	}
	return parts, nil
}
func StringifyDev(devs []string) string {
	return strings.Join(devs, StandardSplitter)
}
