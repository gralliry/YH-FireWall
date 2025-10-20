package system

import (
	"YH-FireWall/core/pkg/funcall"
	"net"
	"strings"
)

type Interface struct {
	Index int      `json:"index"`
	Name  string   `json:"name"`
	MAC   string   `json:"mac"`
	MTU   int      `json:"mtu"`
	Flags []string `json:"flags"`
	Addrs []string `json:"addrs"`
}

func GetInterfaces() ([]Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	return funcall.Convert(
		interfaces,
		func(ifs net.Interface) Interface {
			oaddrs, _ := ifs.Addrs()
			return Interface{
				Index: ifs.Index,
				Name:  ifs.Name,
				MAC:   ifs.HardwareAddr.String(),
				MTU:   ifs.MTU,
				Flags: strings.Split(ifs.Flags.String(), "|"),
				Addrs: funcall.Convert(
					oaddrs,
					func(a net.Addr) string {
						return a.String()
					},
				),
			}
		},
	), nil
}
