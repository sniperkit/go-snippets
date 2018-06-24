package ip

import (
	"bytes"
	"net"
)

type IPRange struct {
	Start net.IP
	End   net.IP
}

func IPInRange(r IPRange, ipaddr net.IP) bool {
	if bytes.Compare(ipaddr, r.Start) >= 0 && bytes.Compare(ipaddr, r.End) < 0 {
		return true
	}

	return false
}
