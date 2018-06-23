package util

import (
	"net"
	"strconv"
)

// lookup portname/service to port - a thin helper over the stdlib
func LookupPort(network, service string) (port int, err error) {
	// first deal with the case that service is convertable to a int
	port, err = strconv.Atoi(service)
	if err == nil {
		return
	}
	// next, perform a service lookup
	return net.LookupPort(network, service)
}
