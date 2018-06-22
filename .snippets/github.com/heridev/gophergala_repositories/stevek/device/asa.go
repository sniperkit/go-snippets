package device

import (
	"fmt"
	"regexp"

	"github.com/gophergala/stevek/util"
)

// Add this device definition to the device.Registry
// This 'registry pattern' is something I likely want to replace with types and interfaces
func init() {
	Registry["Cisco ASA"] = Device{
		Name:      "Cisco ASA",
		filterFn:  asaFilter,
		Transform: Identity,
	}
}

// TODO - crying out for more generalisation
//  - lots of ASA/iptables duplication. Make a firewall interface?
//  - need to work on composiation and idomatic go

var asaRules = []rule{
	//{apply: testRegex(asaRe, config), action: permitAction},
	{apply: testAsaRegex, action: permitAction},
}

var asaRe = []*regexp.Regexp{
	regexp.MustCompile(`^\s*access-list 100 extended permit (tcp|udp) any host ([^\s]+) eq ([^\s]+)\s*$`),
}

// Note need to treat permit/drop differently
// Need a lot more structure here
func asaFilter(msg message, config config) (verdict bool) {
	for _, rule := range asaRules {
		verdict = rule.apply(msg, config)
		// need to map verdict to action here
		if verdict {
			return
		}
	}
	// need to implement a default policy
	// (and also perhaps abstract a bit with a firewall interface?)
	return
}

// Needs to be generalised - I am moving to configurable actions
func testAsaRegex(msg message, config config) bool {
	for _, re := range asaRe {
		if asaPortFilter(msg, re, config) {
			return true
		}
	}
	return false
}

func asaPortFilter(msg message, re *regexp.Regexp, config config) (permit bool) {
	info := re.FindStringSubmatch(string(config))
	port, err := util.LookupPort(info[1], info[3])
	if err != nil {
		// need to fail in a sensible way - TBC!
		fmt.Println(err)
	}
	if msg.Port == port {
		permit = true
	}
	return
}
