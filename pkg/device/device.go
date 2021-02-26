package device

import (
	"net"
	"regexp"
)

// GetNetworkDevices fetch network devices by a regex filtering device name
// todo: regex expression validation
func GetNetworkDevices(regex string) ([]net.Interface, error) {
	var ret []net.Interface
	filter := regexp.MustCompile(regex)
	ifaces, err := net.Interfaces()

	if err != nil {
		return ifaces, err
	}

	for _, dev := range ifaces {
		if filter.MatchString(dev.Name) {
			ret = append(ret, dev)
		}
	}

	return ret, nil
}

// GetAllNetworkDevices fetch all the network devices
func GetAllNetworkDevices() ([]net.Interface, error) {
	return GetNetworkDevices(".*")
}
