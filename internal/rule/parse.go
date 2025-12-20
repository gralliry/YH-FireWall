package rule

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"net"
	"sort"
	"strconv"
	"strings"
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
		// 1. 尝试 CIDR
		if _, ipnet, err := net.ParseCIDR(p); err == nil {
			m = append(m, *ipnet)
			continue
		}
		// 2. 尝试单个 IP
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
		// Range port {start}-{end}
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
			// Single port
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
	// 排序
	sort.Slice(tmp, func(i, j int) bool {
		if tmp[i][0] == tmp[j][0] {
			return tmp[i][1] < tmp[j][1]
		}
		return tmp[i][0] < tmp[j][0]
	})
	// 合并区间
	merged := make([][2]uint16, 0, len(tmp))
	for _, r := range tmp {
		if len(merged) == 0 {
			merged = append(merged, r)
			continue
		}
		lastIndex := len(merged) - 1
		// 如果当前区间的起始端口与上一个区间的结束端口相邻或重叠，则合并
		if r[0] <= merged[lastIndex][1]+1 {
			// 更新结束端口为两者中的较大值
			if r[1] > merged[lastIndex][1] {
				merged[lastIndex][1] = r[1]
			}
		} else {
			// 否则添加新区间
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
	for _, p := range split(raw) {
		if ifi, err := net.InterfaceByName(p); err == nil {
			m[uint32(ifi.Index)] = struct{}{}
		}
	}
	return m, nil
}

// 手动建立协议名称到协议类型的映射，因为 layers.IPProtocol 的 String 方法才返回人类可读名称
var protocolName2Protocol = map[string]layers.IPProtocol{
	strings.ToLower(layers.IPProtocolIPv6HopByHop.String()):    layers.IPProtocolIPv6HopByHop,
	strings.ToLower(layers.IPProtocolICMPv4.String()):          layers.IPProtocolICMPv4,
	strings.ToLower(layers.IPProtocolIGMP.String()):            layers.IPProtocolIGMP,
	strings.ToLower(layers.IPProtocolIPv4.String()):            layers.IPProtocolIPv4,
	strings.ToLower(layers.IPProtocolTCP.String()):             layers.IPProtocolTCP,
	strings.ToLower(layers.IPProtocolUDP.String()):             layers.IPProtocolUDP,
	strings.ToLower(layers.IPProtocolRUDP.String()):            layers.IPProtocolRUDP,
	strings.ToLower(layers.IPProtocolIPv6.String()):            layers.IPProtocolIPv6,
	strings.ToLower(layers.IPProtocolIPv6Routing.String()):     layers.IPProtocolIPv6Routing,
	strings.ToLower(layers.IPProtocolIPv6Fragment.String()):    layers.IPProtocolIPv6Fragment,
	strings.ToLower(layers.IPProtocolGRE.String()):             layers.IPProtocolGRE,
	strings.ToLower(layers.IPProtocolESP.String()):             layers.IPProtocolESP,
	strings.ToLower(layers.IPProtocolAH.String()):              layers.IPProtocolAH,
	strings.ToLower(layers.IPProtocolICMPv6.String()):          layers.IPProtocolICMPv6,
	strings.ToLower(layers.IPProtocolNoNextHeader.String()):    layers.IPProtocolNoNextHeader,
	strings.ToLower(layers.IPProtocolIPv6Destination.String()): layers.IPProtocolIPv6Destination,
	strings.ToLower(layers.IPProtocolOSPF.String()):            layers.IPProtocolOSPF,
	strings.ToLower(layers.IPProtocolIPIP.String()):            layers.IPProtocolIPIP,
	strings.ToLower(layers.IPProtocolEtherIP.String()):         layers.IPProtocolEtherIP,
	strings.ToLower(layers.IPProtocolVRRP.String()):            layers.IPProtocolVRRP,
	strings.ToLower(layers.IPProtocolSCTP.String()):            layers.IPProtocolSCTP,
	strings.ToLower(layers.IPProtocolUDPLite.String()):         layers.IPProtocolUDPLite,
	strings.ToLower(layers.IPProtocolMPLSInIP.String()):        layers.IPProtocolMPLSInIP,
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
