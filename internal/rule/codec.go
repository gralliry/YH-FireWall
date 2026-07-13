package rule

import (
	"YH-FireWall/internal/itable"
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	"github.com/google/gopacket/layers"
)

var pName2p map[string]layers.IPProtocol

func init() {
	pName2p = make(map[string]layers.IPProtocol)
	for p := layers.IPProtocol(0); p < 255; p++ {
		name := p.String()
		if strings.HasPrefix(name, "Unknown(") {
			continue
		}
		pName2p[strings.ToLower(name)] = p
	}
}

func split(raw string) []string {
	return strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == ';' || r == '\t' || r == ' '
	})
}

func GetProtocolNames() []string {
	names := make([]string, 0)
	for k := range pName2p {
		names = append(names, k)
	}
	return names
}

// parse =======================================================================

func parsePrefix(raw string) ([]netip.Prefix, error) {
	parts := split(raw)
	if len(parts) == 0 {
		return nil, nil
	}

	prefixes := make([]netip.Prefix, 0, len(parts))
	for _, p := range parts {
		// CIDR
		if prefix, err := netip.ParsePrefix(p); err == nil {
			prefixes = append(prefixes, prefix.Masked())
			continue
		}
		// 单个 IP
		addr, err := netip.ParseAddr(p)
		if err != nil {
			return nil, fmt.Errorf("invalid IP/CIDR %q", p)
		}
		prefixes = append(prefixes, netip.PrefixFrom(addr, addr.BitLen()))
	}
	return prefixes, nil
}

func parsePort(raw string) ([][2]uint16, error) {
	parts := split(raw)
	if len(parts) == 0 {
		return nil, nil
	}
	ranges := make([][2]uint16, 0)
	for _, part := range parts {
		if strings.Contains(part, "-") {
			se := strings.SplitN(part, "-", 2)
			if len(se) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}
			start, err1 := strconv.Atoi(strings.TrimSpace(se[0]))
			if err1 != nil {
				return nil, fmt.Errorf("failed to parse port: %s", part)
			}
			if start < 0 || start > 65535 {
				return nil, fmt.Errorf("port out of range: %s", part)
			}
			end, err2 := strconv.Atoi(strings.TrimSpace(se[1]))
			if err2 != nil {
				return nil, fmt.Errorf("failed to parse port: %s", part)
			}
			if end < 0 || end > 65535 {
				return nil, fmt.Errorf("port out of range: %s", part)
			}
			if start > end {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}
			ranges = append(ranges, [2]uint16{uint16(start), uint16(end)})
		} else {
			val, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			if val < 0 || val > 65535 {
				return nil, fmt.Errorf("port out of range: %s", part)
			}
			ranges = append(ranges, [2]uint16{uint16(val), uint16(val)})
		}
	}
	return ranges, nil
}

func parseDev(raw string) ([]uint32, error) {
	parts := split(raw)
	if len(parts) == 0 {
		return nil, nil
	}
	m := make([]uint32, 0)
	for _, name := range parts {
		if idx, ok := itable.LookupByName(name); ok {
			m = append(m, uint32(idx))
		}
	}
	return m, nil
}

func parseProtocol(raw string) (map[layers.IPProtocol]struct{}, error) {
	parts := split(raw)
	if len(parts) == 0 {
		return nil, nil
	}
	m := make(map[layers.IPProtocol]struct{})
	for _, p := range parts {
		ptcStr := strings.ToLower(p)
		if ptc, ok := pName2p[ptcStr]; ok {
			m[ptc] = struct{}{}
		} else {
			return nil, fmt.Errorf("invalid protocol: %s", p)
		}
	}
	return m, nil
}

// stringify =======================================================================

func stringifyIPNet(nets []netip.Prefix) string {
	var parts []string
	for _, n := range nets {
		parts = append(parts, n.String())
	}
	return strings.Join(parts, ",")
}

func stringifyPort(ports [][2]uint16) string {
	var parts []string
	for _, p := range ports {
		if p[0] == p[1] {
			parts = append(parts, fmt.Sprintf("%d", p[0]))
		} else {
			parts = append(parts, fmt.Sprintf("%d-%d", p[0], p[1]))
		}
	}
	return strings.Join(parts, ",")
}

func stringifyDev(devs []string) string {
	return strings.Join(devs, ",")
}

func stringifyProtocol(protocols []string) string {
	var parts []string
	for p := range protocols {
		parts = append(parts, strings.ToLower(p.String()))
	}
	return strings.Join(parts, ",")
}
