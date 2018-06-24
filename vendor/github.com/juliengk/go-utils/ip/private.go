package ip

import (
	"net"
)

var privateRanges = []IPRange{
	IPRange{
		Start: net.ParseIP("10.0.0.0"),
		End:   net.ParseIP("10.255.255.255"),
	},
	IPRange{
		Start: net.ParseIP("172.16.0.0"),
		End:   net.ParseIP("172.31.255.255"),
	},
	IPRange{
		Start: net.ParseIP("192.168.0.0"),
		End:   net.ParseIP("192.168.255.255"),
	},
}

func IsPrivateSubnet(ipaddr net.IP) bool {
	if ipCheck := ipaddr.To4(); ipCheck != nil {
		for _, r := range privateRanges {
			if IPInRange(r, ipaddr) {
				return true
			}
		}
	}

	return false
}
