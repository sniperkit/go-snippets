package driver

import (
	"github.com/kassisol/metadata/api/types"
)

type Storager interface {
	AddData(name, dtype, value, description string) error
	UpdateData(name, value, description string) error
	ListData(map[string]string) []types.Data
	RemoveData(name string) error
	AddDataToProfile(profile, data string)
	RemoveDataFromProfile(profile, data string)
	CountData() int

	AddProfile(name string) error
	ListProfile(filter map[string]string) map[string][]string
	RemoveProfile(name string) error
	CountProfile() int

	AddIP(ipaddr, netmask, gateway string) error
	ListIP(filter map[string]string) []types.IP
	RemoveIP(ipaddr string) error
	CountIP() int

	AddInterface(index int, mac, ip, floatingIP string) error
	UpdateInterface(mac, itype, value string) error
	ListInterface(filter map[string]string) []types.Interface
	RemoveInterface(mac string) error
	CountInterface() int

	AddHost(enable bool, name, fqdn, profile string, interfaces []string) error
	ListHost(filter map[string]string) []types.Host
	RemoveHost(name string) error
	EnableHost(name string) error
	DisableHost(name string) error
	CountHost() int

	GetIDFromIP(ip string) int

	GetID(srvid int) int
	GetHostname(srvid int) string
	GetFQDN(srvid int) string
	GetUserData(srvid int) string
	GetVendorData(srvid int) string
	GetPublicKeys(srvid int) []string
	GetRegion(srvid int) string
	GetInterfaces(srvid int) []string
	GetInterfacesType(srvid int, itype string) []int
	GetEnumeratedInterface(srvid int, itype string, index int) []string
	GetInterfaceMACAddress(srvid int, itype string, index int) string
	GetInterfaceType(srvid int, itype string, index int) string
	GetInterfaceIPv4Address(srvid int, itype string, index int) string
	GetInterfaceIPv4Netmask(srvid int, itype string, index int) string
	GetInterfaceIPv4Gateway(srvid int, itype string, index int) string
	FloatingIPExists(srvid int, itype string, index int) bool
	GetInterfaceFloatingIPAddress(srvid int, itype string, index int) string
	GetInterfaceFloatingIPNetmask(srvid int, itype string, index int) string
	GetInterfaceFloatingIPGateway(srvid int, itype string, index int) string
	GetDNSIndex(srvid int) []string
	GetDNSNameservers(srvid int) []string
	GetDNSSearchDomains(srvid int) []string
	GetDNSOptions(srvid int) []string
	GetTags(srvid int) []string
	GetKeys(srvid int) []string
	GetKey(srvid int, key string) string

	End()
}
