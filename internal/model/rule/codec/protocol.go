package codec

import (
	"strings"

	"github.com/google/gopacket/layers"
)

type (
	Name2Protocol func(name string) (layers.IPProtocol, bool)
	Protocol2Name func(p layers.IPProtocol) (string, bool)
)

func ParseProtocol(raw string, mapper Name2Protocol) []layers.IPProtocol {
	names := split(raw)
	if len(names) == 0 {
		return nil
	}
	psl := make([]layers.IPProtocol, 0, len(names))
	for _, name := range names {
		if ps, exist := mapper(strings.ToLower(name)); exist {
			psl = append(psl, ps)
		}
	}
	return psl
}

func StringifyProtocol(ps []layers.IPProtocol, mapper Protocol2Name) string {
	psl := make([]string, 0, len(ps))
	for _, p := range ps {
		if name, exist := mapper(p); exist {
			psl = append(psl, name)
		}
	}
	return strings.Join(psl, StandardSplitter)
}
