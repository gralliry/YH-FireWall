package protocol

import (
	"strings"

	"github.com/google/gopacket/layers"
)

var (
	name2protocol = make(map[string]layers.IPProtocol)
	protocol2name = make(map[layers.IPProtocol]string)
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
		name2protocol[name] = p
		protocol2name[p] = name
	}
}

func Name2Protocol(name string) (layers.IPProtocol, bool) {
	name = strings.ToLower(name)
	p, exist := name2protocol[name]
	return p, exist
}

func Protocol2Name(p layers.IPProtocol) (string, bool) {
	n, exist := protocol2name[p]
	return n, exist
}
