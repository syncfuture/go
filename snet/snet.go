package snet

import (
	"net"

	"github.com/syncfuture/go/serr"
)

func GetLocalV4IP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, serr.WithStack(err)
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, serr.WithStack(err)
		}
		for _, addr := range addrs {
			ip := getv4IPFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, serr.New("connected to the network?")
}

func getv4IPFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
