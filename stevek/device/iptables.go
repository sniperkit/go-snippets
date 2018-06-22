package device

import (
	"fmt"
	"regexp"

	"github.com/gophergala/stevek/util"
)

// Add this device definition to the device.Registry
// This 'registry pattern' is something I likely want to replace with types and interfaces
func init() {
	Registry["iptables"] = Device{
		Name:      "iptables host",
		filterFn:  iptablesFilter,
		Transform: Identity,
	}
}

// TODO - crying out for more generalisation
//  - lots of ASA/iptables duplication. Make a firewall interface?
//  - need to work on composiation and idomatic go

var iptablesRules = []rule{
	//{apply: testRegex(asaRe, config), action: permitAction},
	{apply: testIptablesRegex, action: permitAction},
}

var iptablesRe = []*regexp.Regexp{
	regexp.MustCompile(`-A INPUT -s ([^\s]+) -p (tcp|udp) -m (?:tcp|udp) --dport ([^\s]+) -j ACCEPT`),
}

// Note need to treat permit/drop differently
// Need a lot more structure here
func iptablesFilter(msg message, config config) (verdict bool) {
	for _, rule := range iptablesRules {
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
func testIptablesRegex(msg message, config config) bool {
	for _, re := range iptablesRe {
		if iptablesPortFilter(msg, re, config) {
			return true
		}
	}
	return false
}

func iptablesPortFilter(msg message, re *regexp.Regexp, config config) (permit bool) {
	info := re.FindStringSubmatch(string(config))
	port, err := util.LookupPort(info[2], info[3])
	if err != nil {
		// need to fail in a sensible way - TBC!
		fmt.Println(err)
	}
	if msg.Port == port {
		permit = true
	}
	return
}
