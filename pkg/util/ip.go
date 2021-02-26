package util

import (
	"net"
	"strings"
)

// GetIPSetFromNetworkInterface return a set of IP by given network interface
func GetIPSetFromNetworkInterface(dev *net.Interface) map[string]struct{} {
	ret := make(map[string]struct{})
	addrs, err := dev.Addrs()
	if err != nil {
		return ret
	}
	for _, addr := range addrs {
		ret[strings.Split(addr.String(), "/")[0]] = struct{}{}
	}
	return ret
}

// GetIPSetFromNetworkInterfaces return a set of IP by given network interfaces,
// duplicate IP will be overwrite
func GetIPSetFromNetworkInterfaces(devs []net.Interface) map[string]struct{} {
	ret := make(map[string]struct{})
	for _, dev := range devs {
		addrs, err := dev.Addrs()
		if err != nil {
			return ret
		}
		for _, addr := range addrs {
			ret[strings.Split(addr.String(), "/")[0]] = struct{}{}
		}
	}
	return ret

}
