package rule

import (
	"YH-FireWall/internal/itable"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/google/gopacket/layers"
)

func split(raw string) []string {
	return strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == ';' || r == '\t' || r == ' '
	})
}

func parseIPNet(raw string) ([]net.IPNet, error) {
	parts := split(raw)
	if len(parts) == 0 {
		return nil, nil
	}
	m := make([]net.IPNet, 0)
	for _, p := range parts {
		if _, ipnet, err := net.ParseCIDR(p); err == nil {
			m = append(m, *ipnet)
			continue
		}
		ip := net.ParseIP(p)
		if ip == nil {
			return nil, fmt.Errorf("invalid ip/net: %s", p)
		}
		var mask net.IPMask
		if ip.To4() != nil {
			mask = net.CIDRMask(32, 32)
		} else {
			mask = net.CIDRMask(128, 128)
		}
		ipnet := net.IPNet{IP: ip, Mask: mask}
		m = append(m, ipnet)
	}
	return m, nil
}

func parsePort(raw string) ([][2]uint16, error) {
	parts := split(raw)
	if len(parts) == 0 {
		return nil, nil
	}
	tmp := make([][2]uint16, 0)
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
			tmp = append(tmp, [2]uint16{uint16(start), uint16(end)})
		} else {
			val, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			if val < 0 || val > 65535 {
				return nil, fmt.Errorf("port out of range: %s", part)
			}
			tmp = append(tmp, [2]uint16{uint16(val), uint16(val)})
		}
	}
	sort.Slice(tmp, func(i, j int) bool {
		if tmp[i][0] == tmp[j][0] {
			return tmp[i][1] < tmp[j][1]
		}
		return tmp[i][0] < tmp[j][0]
	})
	merged := make([][2]uint16, 0, len(tmp))
	for _, r := range tmp {
		if len(merged) == 0 {
			merged = append(merged, r)
			continue
		}
		lastIndex := len(merged) - 1
		if r[0] <= merged[lastIndex][1]+1 {
			if r[1] > merged[lastIndex][1] {
				merged[lastIndex][1] = r[1]
			}
		} else {
			merged = append(merged, r)
		}
	}
	return merged, nil
}

func parseDev(raw string) (map[uint32]struct{}, error) {
	parts := split(raw)
	if len(parts) == 0 {
		return nil, nil
	}
	m := make(map[uint32]struct{})
	for _, name := range parts {
		if idx, ok := itable.LookupByName(name); ok {
			m[uint32(idx)] = struct{}{}
		}
	}
	return m, nil
}

var protocolName2Protocol map[string]layers.IPProtocol

func init() {
	protocolName2Protocol = make(map[string]layers.IPProtocol)
	for p := layers.IPProtocol(0); p < 255; p++ {
		name := p.String()
		if strings.HasPrefix(name, "Unknown(") {
			continue
		}
		protocolName2Protocol[strings.ToLower(name)] = p
	}
}

func GetAllProtocolNames() []string {
	protocols := make([]string, 0, len(protocolName2Protocol))
	for k := range protocolName2Protocol {
		protocols = append(protocols, k)
	}
	return protocols
}

func parseProtocol(raw string) (map[layers.IPProtocol]struct{}, error) {
	parts := split(raw)
	if len(parts) == 0 {
		return nil, nil
	}
	m := make(map[layers.IPProtocol]struct{})
	for _, p := range parts {
		ptcStr := strings.ToLower(p)
		if ptc, ok := protocolName2Protocol[ptcStr]; ok {
			m[ptc] = struct{}{}
		} else {
			return nil, fmt.Errorf("invalid protocol: %s", p)
		}
	}
	return m, nil
}

func stringifyIPNet(nets []net.IPNet) string {
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

func stringifyDev(devs map[uint32]struct{}) string {
	var parts []string
	for k := range devs {
		if ifi, err := net.InterfaceByIndex(int(k)); err == nil {
			parts = append(parts, ifi.Name)
		}
	}
	return strings.Join(parts, ",")
}

func stringifyProtocol(protocols map[layers.IPProtocol]struct{}) string {
	var parts []string
	for p := range protocols {
		parts = append(parts, strings.ToLower(p.String()))
	}
	return strings.Join(parts, ",")
}
