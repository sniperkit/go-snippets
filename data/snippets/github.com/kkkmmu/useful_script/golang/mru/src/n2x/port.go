package n2x

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

var ErrorNoOSPFSession = errors.New("No OSPF Session exist")
var ErrorNoBGPSession = errors.New("No BGP Session exist")
var ErrorNoTrafficSession = errors.New("No Traffic session exist")
var ErrorNoTrafficByName = errors.New("No Traffic session exist")

type PortMediaType int

type Port struct {
	Name    string
	Handler string
	*NSession
	MediaTypes              map[PortMediaType]string
	MediaType               PortMediaType
	LegacyLinkSutMac        string
	LegacyLinkSutMac6       string
	LegacyLinkSutIP         map[string]string
	AddressPools            map[string]*AddressPool
	LegacyLinkSutIP6        map[string]string
	AddressPools6           map[string]*AddressPool
	LegacyLinkDefaultSutIP6 string
	LegacyLinkDefaultSutIP  string
	OSPFs                   map[string]*OSPF
	BGPs                    map[string]*BGP
	Traffics                map[string]*Traffic
}

const (
	RJ45 = iota
	SFP
)

func (pmt PortMediaType) String() string {
	if pmt == RJ45 {
		return "AGT_MEDIA_RJ45"
	} else if pmt == SFP {
		return "AGT_MEDIA_SFP"
	} else {
		return "UNKNOWN"
	}
}

func (p *Port) LegacyLinkGetAllAddressPools() ([]*AddressPool, error) {
	cmd := fmt.Sprintf("AgtEthernetAddresses ListAddressPools %s", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get  ip address pools of port: %s : %s", p.Name, err.Error())
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.TrimSpace(res)
	fields := strings.Split(res, " ")

	/*Direct drop existed entries */
	p.AddressPools = make(map[string]*AddressPool, 2)

	pools := make([]*AddressPool, 0, 2)

	for i, field := range fields {
		if strings.TrimSpace(field) != "" {
			np := &AddressPool{
				Handler: field,
				Family:  ADDRESS_FAMILY_IPV4,
				Type:    ADDRESS_POOL_LEGACY,
				Port:    p,
			}
			if i == 0 {
				np.Default = true
			} else {
				np.Default = false
			}

			err := np.GetTesterIPAddresses()
			if err != nil {
				return nil, fmt.Errorf("Cannot restore test IP address %s", err)
			}
			pools = append(pools, np)
			p.AddressPools[field] = np
		}
	}

	return pools, nil
}

func (p *Port) LegacyLinkGetAllIP6AddressPools() ([]*AddressPool, error) {
	cmd := fmt.Sprintf("AgtEthernetIpv6Addresses2 ListAddressPools %s", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get ipv6 address pools of port: %s : %s", p.Name, err.Error())
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.TrimSpace(res)
	fields := strings.Split(res, " ")

	/*Direct drop existed entries */
	p.AddressPools6 = make(map[string]*AddressPool, 2)

	pools := make([]*AddressPool, 0, 2)

	for i, field := range fields {
		if strings.TrimSpace(field) != "" {
			np := &AddressPool{
				Handler: field,
				Family:  ADDRESS_FAMILY_IPV6,
				Type:    ADDRESS_POOL_LEGACY,
				Port:    p,
			}
			if i == 0 {
				np.Default = true
			} else {
				np.Default = false
			}

			err := np.GetTesterIP6Addresses()
			if err != nil {
				return nil, fmt.Errorf("Cannot restore test IP6 address %s", err)
			}
			pools = append(pools, np)
			p.AddressPools6[field] = np
		}
	}

	return pools, nil
}

func (p *Port) LegacyLinkRemoveAllAddressPools() error {
	pools, err := p.LegacyLinkGetAllAddressPools()
	if err != nil {
		return fmt.Errorf("Cannot delete all address pools on port: %s: %s", p.Name, err)
	}

	for _, pool := range pools {
		if !pool.Default {
			err := p.LegacyLinkRemoveAddressPool(pool.Handler)
			if err != nil {
				return fmt.Errorf("Cannot delete address pool: %s:%s with: %s", pool.Handler, pool.First, err)
			}

			delete(p.AddressPools, pool.Handler)
		} else {
			err := pool.SetTesterAddress("0", pool.Handler+".255.255.255", "30", "00:ff:ff:ff:ff:ff")
			if err != nil {
				return fmt.Errorf("Cannot reset tester address pool: %s:%s with: %s", pool.Handler, pool.First, err)
			}
		}
	}

	return nil
}

func (p *Port) LegacyLinkRemoveAllIP6AddressPools() error {
	pools, err := p.LegacyLinkGetAllIP6AddressPools()
	if err != nil {
		return fmt.Errorf("Cannot delete all ipv6 address pools on port: %s: %s", p.Name, err)
	}

	for _, pool := range pools {
		if !pool.Default {
			err := p.LegacyLinkRemoveIP6AddressPool(pool.Handler)
			if err != nil {
				return fmt.Errorf("Cannot delete ipv6 address pool: %s:%s with: %s", pool.Handler, pool.First, err)
			}

			delete(p.AddressPools6, pool.Handler)
		}
	}

	return nil
}

func (p *Port) LegacyLinkSet(vid, testerip, tplen, testmac, sutip string) error {
	err := p.LegacyLinkReset()
	if err != nil {
		return fmt.Errorf("Cannot clear the previous configuration on port: %s: %s", p.Name, err)
	}

	//Use the default pool, no need to create a new one.
	/*
		handler, err := p.LegacyLinkAddAddressPool(testerip, tplen, "1", "1")
		if err != nil {
			return fmt.Errorf("Cannot set test ip address for port %s with: %s", p.Name, err)
		}
	*/
	err = p.LegacyLinkAddSutIPAddress(sutip)
	if err != nil {
		return fmt.Errorf("Cannot set sut ip address for port %s with: %s", p.Name, err)
	}

	pool, err := p.LegacyLinkGetDefaultAddressPool()
	if err != nil {
		return fmt.Errorf("Cannot set sut ip address for port %s with: %s", p.Name, err)
	}

	err = pool.SetTesterAddress(vid, testerip, tplen, testmac)
	if err != nil {
		return fmt.Errorf("Cannot set tester ip address for port %s with: %s", p.Name, err)
	}

	return nil
}

func (p *Port) LegacyLinkSet6(vid, testerip, tplen, testmac, sutip string) error {
	err := p.LegacyLinkReset6()
	if err != nil {
		return fmt.Errorf("Cannot clear the previous configuration on port: %s: %s", p.Name, err)
	}

	//Use the default pool, no need to create a new one.
	/*
		handler, err := p.LegacyLinkAddAddressPool(testerip, tplen, "1", "1")
		if err != nil {
			return fmt.Errorf("Cannot set test ip address for port %s with: %s", p.Name, err)
		}
	*/
	err = p.LegacyLinkAddSutIP6Address(sutip)
	if err != nil {
		return fmt.Errorf("Cannot set sut ipv6 address for port %s with: %s", p.Name, err)
	}

	pool, err := p.LegacyLinkGetDefaultIP6AddressPool()
	if err != nil {
		return fmt.Errorf("Cannot set sut ipv6 address for port %s with: %s", p.Name, err)
	}

	err = pool.SetTesterAddress6(vid, testerip, tplen, testmac)
	if err != nil {
		return fmt.Errorf("Cannot set tester ip address for port %s with: %s", p.Name, err)
	}

	return nil
}

func (p *Port) LegacyLinkGetDefaultAddressPool() (*AddressPool, error) {
	if len(p.AddressPools) == 0 {
		return nil, fmt.Errorf("There is no pool found on port: %s", p.Name)
	}

	for _, pool := range p.AddressPools {
		if pool.Default {
			return pool, nil
		}
	}

	return nil, fmt.Errorf("There is no default pool on port: %s", p.Name)
}

func (p *Port) LegacyLinkGetDefaultIP6AddressPool() (*AddressPool, error) {
	if len(p.AddressPools6) == 0 {
		return nil, fmt.Errorf("There is no ipv6 pool found on port: %s", p.Name)
	}

	for _, pool := range p.AddressPools6 {
		if pool.Default {
			return pool, nil
		}
	}

	return nil, fmt.Errorf("There is no default ipv6 pool on port: %s", p.Name)
}

func (p *Port) LegacyLinkRemoveAllSutIPAddresses() error {
	ips, err := p.LegacyLinkGetSutIPAddresses()
	if err != nil {
		return fmt.Errorf("Cannot reset sut ip on port: %s %s", p.Name, err)
	}

	for _, ip := range ips {
		err = p.LegacyLinkRemoveSutIPAddress(ip)
		if err != nil {
			return fmt.Errorf("Cannot delet sut ip %s on port: %s %s", ip, p.Name, err)
		}

		delete(p.LegacyLinkSutIP, ip)
	}

	return nil
}

func (p *Port) LegacyLinkRemoveAllSutIP6Addresses() error {
	ips, err := p.LegacyLinkGetSutIP6Addresses()
	if err != nil {
		return fmt.Errorf("Cannot reset sut ipv6 on port: %s %s", p.Name, err)
	}

	for _, ip := range ips {
		err = p.LegacyLinkRemoveSutIP6Address(ip)
		if err != nil {
			return fmt.Errorf("Cannot delet sut ipv6 %s on port: %s %s", ip, p.Name, err)
		}

		delete(p.LegacyLinkSutIP6, ip)
	}

	return nil
}

func (p *Port) LegacyLinkReset() error {
	err := p.LegacyLinkRemoveAllSutIPAddresses()
	if err != nil {
		return err
	}
	return p.LegacyLinkRemoveAllAddressPools()
}

func (p *Port) LegacyLinkReset6() error {
	err := p.LegacyLinkRemoveAllSutIP6Addresses()
	if err != nil {
		return err
	}
	return p.LegacyLinkRemoveAllIP6AddressPools()
}

func (p *Port) LegacyLinkAddAddressPool(ip, plen, count, step string) (string, error) {
	cmd := fmt.Sprintf("AgtEthernetAddresses AddAddressPool %s %s %s %s %s", p.Handler, ip, plen, count, step)
	res, err := p.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot add ipv6 address pools of port: %s : %s", p.Name, err.Error())
	}

	np := &AddressPool{
		Handler: strings.TrimSpace(res),
		Family:  ADDRESS_FAMILY_IPV4,
		Type:    ADDRESS_POOL_LEGACY,
		Port:    p,
		First:   ip,
		Plen:    plen,
		Count:   count,
		Step:    step,
		Default: false,
	}

	if p.AddressPools == nil {
		p.AddressPools = make(map[string]*AddressPool, 2)
	}

	p.AddressPools[np.Handler] = np

	return np.Handler, nil
}

//AgtEthernetIpv6Addresses2 AddAddressPool
func (p *Port) LegacyLinkAddIP6AddressPool(ip, plen, count, step string) (string, error) {
	cmd := fmt.Sprintf("AgtEthernetIpv6Addresses2 AddAddressPool %s %s %s %s %s", p.Handler, ip, plen, count, step)
	res, err := p.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot add ipv6 address pools of port: %s : %s", p.Name, err.Error())
	}

	np := &AddressPool{
		Handler: strings.TrimSpace(res),
		Family:  ADDRESS_FAMILY_IPV6,
		Type:    ADDRESS_POOL_LEGACY,
		Port:    p,
		First:   ip,
		Plen:    plen,
		Count:   count,
		Step:    step,
		Default: false,
	}

	if p.AddressPools6 == nil {
		p.AddressPools6 = make(map[string]*AddressPool, 2)
	}

	p.AddressPools6[np.Handler] = np

	return np.Handler, nil
}

func (p *Port) LegacyLinkRemoveAddressPool(handler string) error {
	cmd := fmt.Sprintf("AgtEthernetAddresses RemoveAddressPool %s %s", p.Handler, handler)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot remove ip address pools of port: %s : %s", p.Name, err.Error())
	}

	return nil
}

//AgtEthernetIpv6Addresses2 RemoveAddressPool
func (p *Port) LegacyLinkRemoveIP6AddressPool(handler string) error {
	cmd := fmt.Sprintf("AgtEthernetIpv6Addresses2 RemoveAddressPool %s %s", p.Handler, handler)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot remove ipv6 address pools of port: %s : %s", p.Name, err.Error())
	}

	return nil
}

func (p *Port) LegacyLinkRemoveAddressPoolByIP(ip string) error {
	var handler string
	for h, pool := range p.AddressPools {
		if pool.First == strings.TrimSpace(ip) {
			handler = h
			break
		}
	}

	if handler == "" {
		return fmt.Errorf("Cannot find address pool with: %s on port: %s", ip, p.Name)
	}

	return p.LegacyLinkRemoveAddressPool(handler)
}

func (p *Port) LegacyLinkRemoveIP6AddressPoolByIP(ip string) error {
	var handler string
	for h, pool := range p.AddressPools6 {
		if pool.First == strings.TrimSpace(ip) {
			handler = h
			break
		}
	}

	if handler == "" {
		return fmt.Errorf("Cannot find ipv6 address pool with: %s on port: %s", ip, p.Name)
	}

	return p.LegacyLinkRemoveIP6AddressPool(handler)
}

func (p *Port) LegacyLinkRemoveAddressPoolBySutIP(ip string) error {
	var handler string
	for h, pool := range p.AddressPools {
		if pool.Sut == strings.TrimSpace(ip) {
			handler = h
			break
		}
	}

	if handler == "" {
		return fmt.Errorf("Cannot find address pool with sut: %s on port: %s", ip, p.Name)
	}

	return p.LegacyLinkRemoveAddressPool(handler)
}

func (p *Port) LegacyLinkRemoveIP6AddressPoolBySutIP(ip string) error {
	var handler string
	for h, pool := range p.AddressPools6 {
		if pool.Sut == strings.TrimSpace(ip) {
			handler = h
			break
		}
	}

	if handler == "" {
		return fmt.Errorf("Cannot find ipv6 address pool with sut: %s on port: %s", ip, p.Name)
	}

	return p.LegacyLinkRemoveIP6AddressPool(handler)
}

func (p *Port) LegacyLinkGetSutIPAddresses() ([]string, error) {
	cmd := fmt.Sprintf("AgtEthernetAddresses ListSutIpAddresses %s", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get sup ip address : %s", err.Error())
	}

	/*If previous entries exist, drop it */
	p.LegacyLinkSutIP = make(map[string]string, 1)

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	fields := strings.Split(res, " ")

	ips := make([]string, 0, 10)
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			continue
		}
		ips = append(ips, strings.TrimSpace(field))
		p.LegacyLinkSutIP[strings.TrimSpace(field)] = strings.TrimSpace(field)
	}

	return ips, nil
}

// AgtEthernetIpv6Addresses2 ListSutIpv6AddressesByVlan
func (p *Port) LegacyLinkGetSutIP6Addresses() ([]string, error) {
	cmd := fmt.Sprintf("AgtEthernetIpv6Addresses2 ListSutIpv6AddressesByVlan %s 0", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get sup ip6 address : %s", err.Error())
	}

	/*If previous entries exist, drop it */
	p.LegacyLinkSutIP6 = make(map[string]string, 1)

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	fields := strings.Split(res, " ")

	ips := make([]string, 0, 10)
	for _, field := range fields {
		ips = append(ips, strings.TrimSpace(field))
		p.LegacyLinkSutIP6[strings.TrimSpace(field)] = strings.TrimSpace(field)
	}

	return ips, nil
}

func (p *Port) LegacyLinkAddSutIPAddress(ip string) error {
	//AgtEthernetAddresses AddSutIpAddress
	cmd := fmt.Sprintf("AgtEthernetAddresses AddSutIpAddress %s %s", p.Handler, ip)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set Add SUT ip %s to port %s : %s", ip, p.Name, err.Error())
	}

	if p.LegacyLinkSutIP == nil {
		p.LegacyLinkSutIP = make(map[string]string, 1)
	}

	p.LegacyLinkSutIP[strings.TrimSpace(ip)] = strings.TrimSpace(ip)

	return nil
}

func (p *Port) LegacyLinkAddSutIP6Address(ip string) error {
	//AgtEthernetAddresses AddSutIpAddress
	cmd := fmt.Sprintf("AgtEthernetIpv6Addresses2 AddSutIpv6Address %s 0 %s", p.Handler, ip)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set Add SUT ipv6 %s to port %s : %s", ip, p.Name, err.Error())
	}

	if p.LegacyLinkSutIP6 == nil {
		p.LegacyLinkSutIP6 = make(map[string]string, 1)
	}

	p.LegacyLinkSutIP6[strings.TrimSpace(ip)] = strings.TrimSpace(ip)

	return nil
}

func (p *Port) LegacyLinkRemoveSutIPAddress(ip string) error {
	//AgtEthernetAddresses AddSutIpAddress
	cmd := fmt.Sprintf("AgtEthernetAddresses RemoveSutIpAddress %s %s", p.Handler, ip)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot remove SUT ip %s from port %s : %s", ip, p.Name, err.Error())
	}

	if p.LegacyLinkSutIP == nil {
		p.LegacyLinkSutIP = make(map[string]string, 1)
	}

	if _, ok := p.LegacyLinkSutIP[strings.TrimSpace(ip)]; ok {
		delete(p.LegacyLinkSutIP, strings.TrimSpace(ip))
	}

	return nil
}

func (p *Port) LegacyLinkRemoveSutIP6Address(ip string) error {
	//AgtEthernetAddresses AddSutIpAddress
	cmd := fmt.Sprintf("AgtEthernetIpv6Addresses2 RemoveSutIpv6Address %s 0 %s", p.Handler, ip)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot remove SUT ipv6 %s from port %s : %s", ip, p.Name, err.Error())
	}

	if p.LegacyLinkSutIP6 == nil {
		p.LegacyLinkSutIP6 = make(map[string]string, 1)
	}

	if _, ok := p.LegacyLinkSutIP6[strings.TrimSpace(ip)]; ok {
		delete(p.LegacyLinkSutIP6, strings.TrimSpace(ip))
	}

	return nil
}

//In normal case we use the address pool to configure the test/sut ip addresses.
//we can create multiple address pool under one port
func (p *Port) LegacyLinkAddSutIPAddresses(start, plen, count, step string) error {
	//AgtEthernetAddresses AddSutIpAddress
	cmd := fmt.Sprintf("AgtEthernetAddresses AddSutIpAddresses %s %s %s %s %s", p.Handler, count, start, plen, step)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot Add SUT ips to port %s : %s", p.Name, err.Error())
	}

	_, err = p.LegacyLinkGetSutIPAddresses()
	if err != nil {
		return fmt.Errorf("Cannot get SUT ips on port %s : %s", p.Name, err.Error())
	}

	return nil
}

func (p *Port) LegacyLinkAddSutIP6Addresses(start, plen, count, step string) error {
	//AgtEthernetAddresses AddSutIpAddress
	cmd := fmt.Sprintf("AgtEthernetIpv6Addresses2 AddSutIpv6Addresses %s %s %s %s %s %s", p.Handler, start, count, step, "0", "0")
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot Add SUT ips to port %s : %s", p.Name, err.Error())
	}

	_, err = p.LegacyLinkGetSutIP6Addresses()
	if err != nil {
		return fmt.Errorf("Cannot get SUT ips on port %s : %s", p.Name, err.Error())
	}

	return nil
}

//SUT's mac address is resoved by arp. Anyway we can also configure
func (p *Port) LegacyLinkGetSutMacAddress() (string, error) {
	cmd := fmt.Sprintf("AgtEthernetAddresses GetSutMacAddress %s", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get SUT mac for port %s : %s", p.Name, err.Error())
	}

	p.LegacyLinkSutMac = strings.TrimSpace(res)

	return strings.TrimSpace(res), nil
}

func (p *Port) LegacyLinkGetSutMacAddress6() (string, error) {
	cmd := fmt.Sprintf("AgtEthernetIpv6Addresses2 GetSutMacAddress %s 0 %s", p.Handler, p.LegacyLinkDefaultSutIP6)
	res, err := p.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get SUT mac for port %s address: %s: %s", p.Name, p.LegacyLinkDefaultSutIP6, err.Error())
	}

	p.LegacyLinkSutMac6 = strings.TrimSpace(res)

	return strings.TrimSpace(res), nil
}

func (p *Port) LegacyLinkSetSutMacAddress(ip, mac string) (string, error) {
	cmd := fmt.Sprintf("AgtEthernetAddresses SetSutMacAddress %s %s %s", p.Handler, ip, mac)
	_, err := p.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot set SUT mac %s for port %s : %s", mac, p.Name, err.Error())
	}

	p.LegacyLinkSutMac = strings.TrimSpace(mac)

	return strings.TrimSpace(mac), nil
}

//AgtEthernetIpv6Addresses2 SetSutMacAddress
func (p *Port) LegacyLinkSetSutMacAddress6(ip, mac string) (string, error) {
	cmd := fmt.Sprintf("AgtEthernetIpv6Addresses2 SetSutMacAddress %s 0 %s %s", p.Handler, ip, mac)
	_, err := p.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot set SUT mac %s for port %s address: %s: %s", mac, p.Name, ip, err.Error())
	}

	p.LegacyLinkSutMac6 = strings.TrimSpace(mac)

	return strings.TrimSpace(mac), nil
}

func (p *Port) SetMediaType(media PortMediaType) error {
	if _, ok := p.MediaTypes[media]; !ok {
		return fmt.Errorf("Unsupported mediay type: %d", media)
	}
	cmd := fmt.Sprintf("AgtOpticalInterface SetMediaType %s %s", p.Handler, p.MediaTypes[media])
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set port %s media type : %s", p.Name, err.Error())
	}
	p.MediaType = media
	return nil
}

func (p *Port) GetMediaType() (PortMediaType, error) {
	cmd := fmt.Sprintf("AgtOpticalInterface GetMediaType %s", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return 0, fmt.Errorf("Cannot set port %s media type : %s", p.Name, err.Error())
	}

	for iType, sType := range p.MediaTypes {
		if sType == strings.TrimSpace(res) {
			p.MediaType = iType
			return iType, nil
		}
	}
	return 0, fmt.Errorf("Unknown port media type :%s for port %s", res, p.Name)
}

func (p *Port) GetAvailableMediaTypes() ([]string, error) {
	cmd := fmt.Sprintf("AgtOpticalInterface ListAvailableMediaTypes %s", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get available media types for port %s : %s", p.Name, err.Error())
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	fields := strings.Split(res, " ")

	if p.MediaTypes == nil {
		p.MediaTypes = make(map[PortMediaType]string, 2)
	}

	types := make([]string, 0, 2)

	for _, field := range fields {
		if strings.TrimSpace(field) == "AGT_MEDIA_SFP" {
			p.MediaTypes[SFP] = strings.TrimSpace(field)
		} else if strings.TrimSpace(field) == "AGT_MEDIA_RJ45" {
			p.MediaTypes[RJ45] = strings.TrimSpace(field)
		} else {
			fmt.Printf("Unknown media type: %s for port: %s", field, p.Name)
			continue
		}
		types = append(types, field)
	}
	return types, nil
}

func (p *Port) LaserOn() error {
	cmd := fmt.Sprintf("AgtOpticalInterface LaserOn %s", p.Handler)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set laseron for port %s : %s", p.Name, err.Error())
	}

	return nil
}

func (p *Port) LaserOff() error {
	cmd := fmt.Sprintf("AgtOpticalInterface LaserOff %s", p.Handler)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set laseroff for port %s : %s", p.Name, err.Error())
	}

	return nil
}

func (p *Port) SendAllArpRequests() error {
	cmd := fmt.Sprintf("AgtArpEmulation SendAllArpRequests %s", p.Handler)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot send arp on port %s : %s", p.Name, err.Error())
	}

	return nil
}

func (p *Port) SendAllNeighborSolicitations() error {
	cmd := fmt.Sprintf("AgtArpEmulation SendAllNeighborSolicitations %s", p.Handler)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot send ns on port %s : %s", p.Name, err.Error())
	}

	return nil
}

func (p *Port) AddOSPF(area, rid, srid, name string) (*OSPF, error) {
	cmd := fmt.Sprintf("AgtTestTopology AddSession %s AGT_SESSION_OSPF", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot add ospf session on port %s : %s", p.Name, err.Error())
	}

	o := &OSPF{
		Handler: strings.TrimSpace(res),
		Version: "2",
		AreaID:  area,
		RID:     rid,
		SutRID:  srid,
		Port:    p,
	}

	err = o.Sync()
	if err != nil {
		return nil, fmt.Errorf("Cannot create ospf with: %s", err.Error())
	}

	err = o.SetName(name)
	if err != nil {
		return nil, fmt.Errorf("Cannot create ospf with: %s", err.Error())
	}

	_, err = o.GetRouterHandler()
	if err != nil {
		return nil, fmt.Errorf("Cannot create ospf with: %s", err.Error())
	}

	err = o.SetAreaID(area)
	if err != nil {
		return nil, fmt.Errorf("Cannot create ospf with: %s", err.Error())
	}
	err = o.SetRouterID(rid)
	if err != nil {
		return nil, fmt.Errorf("Cannot create ospf with: %s", err.Error())
	}

	err = o.SetSutRouterID(srid)
	if err != nil {
		return nil, fmt.Errorf("Cannot create ospf with: %s", err.Error())
	}

	return o, nil
}

func (p *Port) GetOSPFByName(name string) (*OSPF, error) {
	if len(p.OSPFs) == 0 {
		_, err := p.GetAllOSPFs()
		if err != nil {
			if err == ErrorNoOSPFSession {
				return nil, err
			}
			return nil, fmt.Errorf("Cannot find ospf by name: %s with: %s", name, err)
		}
	}

	name = strings.TrimSpace(name)
	for _, ospf := range p.OSPFs {
		if ospf.Name == name {
			return ospf, nil
		}
	}

	return nil, fmt.Errorf("Cannot find ospf session by name: %s", name)
}

func (p *Port) RemoveOSPF(o *OSPF) error {
	cmd := fmt.Sprintf("AgtTestTopology RemoveSession %s", o.Handler)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot remove ospf session on port %s : %s", p.Name, err.Error())
	}

	return nil
}

//AgtTestTopology ListSessions
func (p *Port) GetAllOSPFs() ([]*OSPF, error) {
	cmd := fmt.Sprintf("AgtTestTopology ListSessions %s", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot remove ospf session on port %s : %s", p.Name, err.Error())
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.TrimSpace(res)

	p.OSPFs = make(map[string]*OSPF, 1)
	ospfs := make([]*OSPF, 0, 1)

	fields := strings.Split(res, " ")
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			continue
		}

		typ, err := p.GetSessionType(field)
		if typ != "AGT_SESSION_OSPF" {
			continue
		}

		nospf := &OSPF{
			Handler: strings.TrimSpace(field),
			Port:    p,
		}

		err = nospf.Sync()
		if err != nil {
			return nil, fmt.Errorf("Cannot get all ospf with: %s", err)
		}

		ospfs = append(ospfs, nospf)
		p.OSPFs[field] = nospf
	}

	if len(ospfs) == 0 {
		return nil, ErrorNoOSPFSession
	}

	return ospfs, nil
}

func (p *Port) DeleteAllOSPFs() error {
	ospfs, err := p.GetAllOSPFs()
	if err != nil {
		if err == ErrorNoOSPFSession {
			return nil
		}
		return fmt.Errorf("Cannot Delete all ospfs with: %s", err)
	}

	for _, o := range ospfs {
		err := p.RemoveOSPF(o)
		if err != nil {
			return fmt.Errorf("Cannot Delete all ospfs with: %s", err)
		}
	}

	return nil
}

func (p *Port) GetEmulationState(handler string) (string, error) {
	cmd := fmt.Sprintf("AgtTestTopology GetEmulationState %s", handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get emulation state session on port %s : %s", p.Name, err.Error())
	}

	return strings.TrimSpace(res), nil
}

func (p *Port) ResetSession(handler string) error {
	cmd := fmt.Sprintf("AgtTestTopology ResetSession %s", handler)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot reset sesssion on port %s : %s", p.Name, err.Error())
	}

	return nil
}

func (p *Port) RemoveSession(handler string) error {
	cmd := fmt.Sprintf("AgtTestTopology RemoveSession %s", handler)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot remove sesssion on port %s : %s", p.Name, err.Error())
	}

	return nil
}

func (p *Port) GetSessionType(handler string) (string, error) {
	cmd := fmt.Sprintf("AgtTestTopology GetSessionType %s", handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get session type on port %s : %s", p.Name, err.Error())
	}

	res = strings.Replace(res, "\"", "", -1)

	return strings.TrimSpace(res), nil
}

func (p *Port) AddBGP(typ int, as, ras, name string) (*BGP, error) {
	tname, ok := BGPTypeNameMap[typ]
	if !ok {
		return nil, fmt.Errorf("Unknown BGP type: %d", typ)
	}

	cmd := fmt.Sprintf("AgtTestTopology AddSession %s %s", p.Handler, tname)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot Add BGP type on port %s : %s", p.Name, err.Error())
	}

	session := BGPTypeSessionMap[typ]

	res = strings.TrimSpace(res)
	if res == "" {
		return nil, fmt.Errorf("Cannot Add BGP with empty return value %s", res)
	}

	bgp := &BGP{
		Name:            name,
		Handler:         res,
		Object:          session,
		ASN:             as,
		RASN:            ras,
		Port:            p,
		PeerPoolObject:  BGPTypePeerPoolMap[typ],
		PeerPoolHandler: res,
		Type:            typ,
	}

	err = bgp.Sync()
	if err != nil {
		return nil, fmt.Errorf("Cannot get bgp instance infor with: %s", err)
	}

	err = bgp.SetSutAsNumber(ras)
	if err != nil {
		return nil, fmt.Errorf("Cannot add bgp with %s", err)
	}

	err = bgp.SetTesterAsNumber(as)
	if err != nil {
		return nil, fmt.Errorf("Cannot add bgp with %s", err)
	}

	return bgp, nil
}

func (p *Port) RemoveBGP(b *BGP) error {
	cmd := fmt.Sprintf("AgtTestTopology RemoveSession %s", b.Handler)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot remove bgp session on port %s : %s", p.Name, err.Error())
	}

	return nil
}

func (p *Port) GetAllBGPs() ([]*BGP, error) {
	cmd := fmt.Sprintf("AgtTestTopology ListSessions %s", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get all session on port %s : %s", p.Name, err.Error())
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.TrimSpace(res)

	p.BGPs = make(map[string]*BGP, 1)
	bgps := make([]*BGP, 0, 1)

	fields := strings.Split(res, " ")
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			continue
		}

		typ, err := p.GetSessionType(field)
		if typ != "AGT_SESSION_BGP4_E_BGP" && typ != "AGT_SESSION_BGP4_I_BGP" &&
			typ != "AGT_SESSION_BGP4_I_BGP_IPV6" && typ != "AGT_SESSION_BGP4_E_BGP_IPV6" {
			continue
		}
		ityp := BGPSessionTypeTypeMap[typ]

		nbgp := &BGP{
			Handler:         strings.TrimSpace(field),
			Object:          BGPTypeSessionMap[ityp],
			PeerPoolObject:  BGPTypePeerPoolMap[ityp],
			PeerPoolHandler: strings.TrimSpace(field),
			Port:            p,
		}

		err = nbgp.Sync()
		if err != nil {
			return nil, fmt.Errorf("Cannot get all bgp with: %s", err)
		}

		bgps = append(bgps, nbgp)
		p.BGPs[field] = nbgp
	}

	if len(bgps) == 0 {
		return nil, ErrorNoBGPSession
	}

	return bgps, nil
}

func (p *Port) GetAllBGPsByType(typ int) ([]*BGP, error) {
	if typ != IPV4_EBGP && typ != IPV4_IBGP &&
		typ != IPV6_EBGP && typ != IPV6_IBGP {
		return nil, fmt.Errorf("Unknown bgp type: %d", typ)
	}

	cmd := fmt.Sprintf("AgtTestTopology ListSessions %s", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get all session on port %s : %s", p.Name, err.Error())
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.TrimSpace(res)

	if p.BGPs == nil {
		p.BGPs = make(map[string]*BGP, 1)
	}
	bgps := make([]*BGP, 0, 1)

	fields := strings.Split(res, " ")
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			continue
		}

		typs, err := p.GetSessionType(field)
		if typs != BGPTypeNameMap[typ] {
			continue
		}

		nbgp := &BGP{
			Object:          BGPTypeNameMap[typ],
			Handler:         strings.TrimSpace(field),
			PeerPoolObject:  BGPTypePeerPoolMap[typ],
			PeerPoolHandler: strings.TrimSpace(field),
			Port:            p,
		}

		err = nbgp.Sync()
		if err != nil {
			return nil, fmt.Errorf("Cannot get all bgp with: %s", err)
		}

		bgps = append(bgps, nbgp)
		p.BGPs[field] = nbgp
	}

	if len(bgps) == 0 {
		return nil, ErrorNoBGPSession
	}

	return bgps, nil
}

func (p *Port) RemoveAllBGPs() error {
	bgps, err := p.GetAllBGPs()
	if err != nil {
		if err == ErrorNoBGPSession {
			return nil
		}
		return fmt.Errorf("Cannot Delete all bgps with: %s", err)
	}

	for _, o := range bgps {
		err := p.RemoveBGP(o)
		if err != nil {
			return fmt.Errorf("Cannot Delete all bgps with: %s", err)
		}
	}

	return nil
}

func (p *Port) RemoveAllIPv4EBGPs() error {
	bgps, err := p.GetAllBGPsByType(IPV4_EBGP)
	if err != nil {
		if err == ErrorNoBGPSession {
			return nil
		}
		return fmt.Errorf("Cannot remove all ipv4 ebgp with: %s", err)
	}

	for _, bgp := range bgps {
		err = p.RemoveBGP(bgp)
		if err != nil {
			return fmt.Errorf("Cannot remove bgp %s with %s", bgp.Handler, err)
		}
	}

	return nil
}

func (p *Port) RemoveAllIPv4IBGPs() error {
	bgps, err := p.GetAllBGPsByType(IPV4_IBGP)
	if err != nil {
		if err == ErrorNoBGPSession {
			return nil
		}
		return fmt.Errorf("Cannot remove all ipv4 ibgp with: %s", err)
	}

	for _, bgp := range bgps {
		err = p.RemoveBGP(bgp)
		if err != nil {
			return fmt.Errorf("Cannot remove bgp %s with %s", bgp.Handler, err)
		}
	}

	return nil
}

func (p *Port) RemoveAllIPv6EBGPs() error {
	bgps, err := p.GetAllBGPsByType(IPV6_EBGP)
	if err != nil {
		if err == ErrorNoBGPSession {
			return nil
		}
		return fmt.Errorf("Cannot remove all ipv6 ebgp with: %s", err)
	}

	for _, bgp := range bgps {
		err = p.RemoveBGP(bgp)
		if err != nil {
			return fmt.Errorf("Cannot remove bgp %s with %s", bgp.Handler, err)
		}
	}

	return nil
}

func (p *Port) RemoveAllIPv6IBGPs() error {
	bgps, err := p.GetAllBGPsByType(IPV6_IBGP)
	if err != nil {
		if err == ErrorNoBGPSession {
			return nil
		}
		return fmt.Errorf("Cannot remove all ipv6 ibgp with: %s", err)
	}

	for _, bgp := range bgps {
		err = p.RemoveBGP(bgp)
		if err != nil {
			return fmt.Errorf("Cannot remove bgp %s with %s", bgp.Handler, err)
		}
	}

	return nil
}

//By default use CONSTANT_PROFILE, it's easy to use.
func (p *Port) AddTraffic(scount, name string) (*Traffic, error) {
	cmd := fmt.Sprintf("AgtProfileList AddProfile %s AGT_CONSTANT_PROFILE", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot add traffic with: %s", err)
	}

	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)
	if res == "" {
		return nil, fmt.Errorf("Invalid return value for add traffic: %s", res)
	}

	if p.Traffics == nil {
		p.Traffics = make(map[string]*Traffic, 1)
	}

	nt := &Traffic{
		Handler:     res,
		StreamCount: scount,
		Name:        name,
		Port:        p,
		Object:      "AgtConstantProfile",
	}

	err = nt.Init()
	if err != nil {
		return nt, fmt.Errorf("Cannot add traffic with: %s", err)
	}

	p.Traffics[name] = nt

	return nt, nil
}

func (p *Port) RemoveTraffic(tr *Traffic) error {
	if tr.Handler == "" {
		return fmt.Errorf("Invalid traffic %+v", tr)
	}

	_, ok := p.Traffics[tr.Name]
	if !ok {
		return fmt.Errorf("Traffic %s does not exist")
	}

	cmd := fmt.Sprintf("AgtProfileList Remove %s", tr.Handler)
	_, err := p.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot remove traffic %s from %s with: %s", tr.Handler, p.Handler, err)
	}

	delete(p.Traffics, tr.Name)

	return nil
}

func (p *Port) RemoveTrafficByName(name string) error {
	t, ok := p.Traffics[name]
	if !ok {
		return ErrorNoTrafficByName
	}

	return p.RemoveTraffic(t)
}

func (p *Port) Sync() error {
	_, err := p.GetAllOSPFs()
	if err != nil && err != ErrorNoOSPFSession {
		return fmt.Errorf("Cannot sync port info with: %s", err)
	}

	_, err = p.GetAllBGPs()
	if err != nil && err != ErrorNoBGPSession {
		return fmt.Errorf("Cannot sync port info with: %s", err)
	}

	_, err = p.GetAllTraffics()
	if err != nil && err != ErrorNoTrafficSession {
		return fmt.Errorf("Cannot sync port info with: %s", err)
	}

	return nil
}

func (p *Port) GetAllTraffics() ([]*Traffic, error) {
	cmd := fmt.Sprintf("AgtProfileList ListProfilesOnPort %s", p.Handler)
	res, err := p.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get traffic profile on port %s with: %s", p.Handler, err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.TrimSpace(res)

	fields := strings.Split(res, " ")
	trs := make([]*Traffic, 0, len(fields))

	if p.Traffics == nil {
		p.Traffics = make(map[string]*Traffic, len(fields))
	}

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		tr := &Traffic{
			Handler: field,
			Object:  "AgtConstantProfile",
			Port:    p,
			Unit:    LoadUnitsNameMap[DEFAULT_TRAFFIC_UNIT],
		}

		err = tr.Sync()
		if err != nil {
			return nil, fmt.Errorf("Cannot sync traffic %s with %s", field, err)
		}

		p.Traffics[tr.Name] = tr

		trs = append(trs, tr)
	}

	return trs, nil
}

func (p *Port) RemoveAllTraffics() error {
	for _, t := range p.Traffics {
		err := p.RemoveTraffic(t)
		if err != nil {
			return fmt.Errorf("Cannot remove all traffic with %s", err)
		}
	}

	return nil
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
