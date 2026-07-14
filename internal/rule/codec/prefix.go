package codec

import (
	"fmt"
	"net/netip"
	"strings"
)

func ParsePrefix(raw string) ([]netip.Prefix, error) {
	parts := Split(raw)
	if len(parts) == 0 {
		return nil, nil
	}
	prefixes := make([]netip.Prefix, 0, len(parts))
	for _, p := range parts {
		if prefix, err := netip.ParsePrefix(p); err == nil {
			prefixes = append(prefixes, prefix.Masked())
			continue
		}
		addr, err := netip.ParseAddr(p)
		if err != nil {
			return nil, fmt.Errorf("invalid IP/CIDR %q", p)
		}
		prefixes = append(prefixes, netip.PrefixFrom(addr, addr.BitLen()))
	}
	return prefixes, nil
}

func StringifyPrefix(ns []netip.Prefix) string {
	var parts []string
	for _, n := range ns {
		parts = append(parts, n.String())
	}
	return strings.Join(parts, StandardSplitter)
}
