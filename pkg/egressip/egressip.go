package egressip

import (
	"net"
)

// Get learns egress IP address.
func Get() (string, error) {
	c, err := net.Dial("udp", "1.1.1.1:1")
	if err != nil {
		return "", err
	}
	defer c.Close()
	host, _, _ := net.SplitHostPort(c.LocalAddr().String())
	return host, nil
}
