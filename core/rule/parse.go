package rule

import (
	"fmt"
	"github.com/google/gopacket/layers"
	"net"
	"sort"
	"strconv"
	"strings"
)

var protocolName2Protocol = map[string]layers.IPProtocol{}

func init() {
	// 可以扩展其他协议
	for _, p := range []layers.IPProtocol{
		layers.IPProtocolTCP,
		layers.IPProtocolUDP,
		layers.IPProtocolICMPv4,
		layers.IPProtocolIGMP,
		layers.IPProtocolSCTP,
		layers.IPProtocolGRE,
		layers.IPProtocolESP,
		layers.IPProtocolAH,
		layers.IPProtocolVRRP,
	} {
		protocolName2Protocol[strings.ToLower(string(p))] = p
	}
}

func split(raw string) []string {
	// 先统一用换行分割，再按逗号拆
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == ';' || r == '\t'
	})
	out := make([]string, 0)
	for _, p := range parts {
		line := strings.TrimSpace(p)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}

func parseIPNet(raw string) ([]net.IPNet, error) {
	m := make([]net.IPNet, 0)
	for _, p := range split(raw) {
		// 1. 尝试 CIDR
		if _, ipnet, err := net.ParseCIDR(p); err == nil {
			m = append(m, *ipnet)
			continue
		}
		// 2. 尝试单个 IP
		ip := net.ParseIP(p)
		if ip == nil {
			return nil, fmt.Errorf("无效输入: %s", p)
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
	tmp := make([][2]uint16, 0)
	for _, part := range split(raw) {
		part = strings.TrimSpace(part)
		// 范围端口 {start}-{end}
		if strings.Contains(part, "-") {
			se := strings.SplitN(part, "-", 2)
			if len(se) != 2 {
				return nil, fmt.Errorf("无效的端口范围: %s", part)
			}
			start, err1 := strconv.Atoi(strings.TrimSpace(se[0]))
			if start < 0 || start > 65535 {
				return nil, fmt.Errorf("端口范围超出范围: %s", part)
			}
			end, err2 := strconv.Atoi(strings.TrimSpace(se[1]))
			if end < 0 || end > 65535 {
				return nil, fmt.Errorf("端口范围超出范围: %s", part)
			}
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("端口范围解析失败: %s", part)
			}
			if start > end {
				start, end = end, start
			}
			tmp = append(tmp, [2]uint16{uint16(start), uint16(end)})
		} else {
			// 单个端口
			val, err := strconv.Atoi(part)
			if val < 0 || val > 65535 {
				return nil, fmt.Errorf("端口超出范围: %s", part)
			}
			if err != nil {
				return nil, fmt.Errorf("无效端口: %s", part)
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
		if len(merged) == 0 || r[0] > merged[len(merged)-1][1]+1 {
			merged = append(merged, r)
		} else {
			if r[1] > merged[len(merged)-1][1] {
				merged[len(merged)-1][1] = r[1]
			}
		}
	}
	return merged, nil
}

func parseDev(raw string) (map[uint32]struct{}, error) {
	m := make(map[uint32]struct{})
	for _, p := range split(raw) {
		line := strings.TrimSpace(p)
		if line == "" {
			continue
		}
		ifi, err := net.InterfaceByName(line)
		if err != nil {
			continue
		}
		m[uint32(ifi.Index)] = struct{}{}
	}
	return m, nil
}

func parseProtocol(raw string) (map[layers.IPProtocol]struct{}, error) {
	m := make(map[layers.IPProtocol]struct{})
	for _, p := range split(raw) {
		line := strings.TrimSpace(p)
		if line == "" {
			continue
		}
		if ptc, ok := protocolName2Protocol[strings.ToLower(line)]; ok {
			m[ptc] = struct{}{}
			continue
		}
	}
	return m, nil
}
