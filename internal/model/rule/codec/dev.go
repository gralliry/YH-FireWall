package codec

import "strings"

func ParseDev(raw string, mapper map[string]uint32) []uint32 {
	names := split(raw)
	if len(names) == 0 {
		return nil
	}
	devs := make([]uint32, 0, len(names))
	for _, name := range names {
		if dev, exist := mapper[name]; exist {
			devs = append(devs, dev)
		}
	}
	return devs
}

func StringifyDev(devs []uint32, mapper map[uint32]string) string {
	names := make([]string, 0, len(devs))
	for _, dev := range devs {
		if name, exist := mapper[dev]; exist {
			names = append(names, name)
		}
	}
	return strings.Join(names, StandardSplitter)
}
