package codec

import "strings"

type (
	DevName2Index func(name string) (index uint32, exist bool)
	DevIndex2Name func(index uint32) (name string, exist bool)
)

func ParseDev(raw string, mapper DevName2Index) []uint32 {
	names := split(raw)
	if len(names) == 0 {
		return nil
	}
	devs := make([]uint32, 0, len(names))
	for _, name := range names {
		if dev, exist := mapper(name); exist {
			devs = append(devs, dev)
		}
	}
	return devs
}

func StringifyDev(devs []uint32, mapper DevIndex2Name) string {
	names := make([]string, 0, len(devs))
	for _, dev := range devs {
		if name, exist := mapper(dev); exist {
			names = append(names, name)
		}
	}
	return strings.Join(names, StandardSplitter)
}
