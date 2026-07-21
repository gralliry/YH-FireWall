package rule

import (
	"fmt"
	"strconv"
	"strings"
)

func Set(key, value string) (*Option, error) {
	o := &Option{}
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "_", "")
	switch key {
	case "group":
		o.Group = &value
	case "comment":
		o.Comment = &value
	case "accept":
		switch strings.ToLower(value) {
		case "true":
			o.Accept = new(true)
		case "false":
			o.Accept = new(false)
		default:
			return nil, fmt.Errorf("invalid bool value: %s", value)
		}
	case "enable":
		switch strings.ToLower(value) {
		case "true":
			o.Enable = new(true)
		case "false":
			o.Enable = new(false)
		default:
			return nil, fmt.Errorf("invalid bool value: %s", value)
		}
	case "priority":
		n, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid int value: %s", value)
		}
		o.Priority = &n
	case "indevs":
		o.InDevs = &value
	case "outdevs":
		o.OutDevs = &value
	case "protocols":
		o.Protocols = &value
	case "srcnets":
		o.SrcNets = &value
	case "dstnets":
		o.DstNets = &value
	case "srcports":
		o.SrcPorts = &value
	case "dstports":
		o.DstPorts = &value
	default:
		return nil, fmt.Errorf("unknown key: %s", key)
	}
	return o, nil
}
