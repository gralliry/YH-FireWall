package codec

import (
	"strings"

	"github.com/google/gopacket/layers"
)

var (
	name2Protocol = make(map[string]layers.IPProtocol)
)

func init() {
	for _, p := range []layers.IPProtocol{
		layers.IPProtocolIPv6HopByHop,
		layers.IPProtocolICMPv4,
		layers.IPProtocolIGMP,
		layers.IPProtocolIPv4,
		layers.IPProtocolTCP,
		layers.IPProtocolUDP,
		layers.IPProtocolRUDP,
		layers.IPProtocolIPv6,
		layers.IPProtocolIPv6Routing,
		layers.IPProtocolIPv6Fragment,
		layers.IPProtocolGRE,
		layers.IPProtocolESP,
		layers.IPProtocolAH,
		layers.IPProtocolICMPv6,
		layers.IPProtocolNoNextHeader,
		layers.IPProtocolIPv6Destination,
		layers.IPProtocolOSPF,
		layers.IPProtocolIPIP,
		layers.IPProtocolEtherIP,
		layers.IPProtocolVRRP,
		layers.IPProtocolSCTP,
		layers.IPProtocolUDPLite,
		layers.IPProtocolMPLSInIP,
	} {
		name := strings.ToLower(p.String())
		name2Protocol[name] = p
	}
}

func ParseProtocol(raw string) []layers.IPProtocol {
	names := split(raw)
	if len(names) == 0 {
		return nil
	}
	psl := make([]layers.IPProtocol, 0, len(names))
	for _, name := range names {
		if ps, exist := name2Protocol[strings.ToLower(name)]; exist {
			psl = append(psl, ps)
		}
	}
	return psl
}

func StringifyProtocol(ps []layers.IPProtocol) string {
	psl := make([]string, len(ps))
	for i, p := range ps {
		psl[i] = strings.ToLower(p.String())
	}
	return strings.Join(psl, StandardSplitter)
}
