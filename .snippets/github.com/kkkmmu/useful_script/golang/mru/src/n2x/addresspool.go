package n2x

import (
	"fmt"
	"strings"
)

const (
	ADDRESS_POOL_LEGACY = iota
	ADDRESS_POOL_NORMAL

	ADDRESS_FAMILY_IPV4 = iota
	ADDRESS_FAMILY_IPV6
)

type AddressPool struct {
	Type    int
	Family  int
	Handler string
	Sut     string
	First   string
	Plen    string
	Count   string
	Step    string
	Default bool
	Mac     string
	Vlan    string
	*Port
}

func (ap *AddressPool) SetSutIPAddress(ip string) error {
	cmd := fmt.Sprintf("AgtEthernetAddressPool SetTesterAndSutIpAddresses %s %s %s %s %s %s %s", ap.Handler, ap.First, ap.Plen, ap.Count, ap.Step, ip, ip)
	_, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set sut ip address for pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}
	return nil
}

func (ap *AddressPool) SetTesterAddress(vid, ip, plen, mac string) error {
	err := ap.SetTesterIPAddress(ip, plen)
	if err != nil {
		return err
	}

	err = ap.SetVlan(vid)
	if err != nil {
		return err
	}

	if vid == "0" {
		err = ap.DisableVlan()
		if err != nil {
			return nil
		}
	}

	return ap.SetTesterMacAddress(mac)
}

func (ap *AddressPool) SetTesterAddress6(vid, ip, plen, mac string) error {
	err := ap.SetTesterIP6Address(ip, plen)
	if err != nil {
		return err
	}

	err = ap.SetVlan6(vid)
	if err != nil {
		return err
	}

	if vid == "0" {
		err = ap.DisableVlan6()
		if err != nil {
			return nil
		}
	}

	return ap.SetTesterMacAddress6(mac)
}

func (ap *AddressPool) SetTesterIPAddress(ip, plen string) error {
	cmd := fmt.Sprintf("AgtEthernetAddressPool SetTesterIpAddresses %s %s %s %s %s", ap.Handler, ip, plen, "1", "1")
	_, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set tester ip address for pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	return ap.GetTesterIPAddresses()
}

func (ap *AddressPool) SetTesterIP6Address(ip, plen string) error {
	cmd := fmt.Sprintf("AgtEthernetIpv6AddressPool SetTesterAddresses %s %s %s %s %s 0", ap.Handler, ip, plen, "1", "::1")
	_, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set tester ip6 address for pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	return ap.GetTesterIP6Addresses()
}

func (ap *AddressPool) SetTesterMacAddress(mac string) error {
	cmd := fmt.Sprintf("AgtEthernetAddressPool SetTesterMacAddresses %s %s 1", ap.Handler, strings.TrimSpace(mac))
	_, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set tester mac address for pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	ap.Mac = strings.TrimSpace(mac)

	return ap.GetTesterMacAddress()
}

//AgtEthernetIpv6AddressPool SetTesterMacAddresses
func (ap *AddressPool) SetTesterMacAddress6(mac string) error {
	cmd := fmt.Sprintf("AgtEthernetIpv6AddressPool SetTesterMacAddresses %s %s 1", ap.Handler, strings.TrimSpace(mac))
	_, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set tester mac address for ipv6 pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	ap.Mac = strings.TrimSpace(mac)

	return ap.GetTesterMacAddress6()
}

func (ap *AddressPool) GetTesterMacAddress() error {
	cmd := fmt.Sprintf("AgtEthernetAddressPool GetTesterMacAddresses %s", ap.Handler)
	res, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set tester mac address for pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)

	fields := strings.Split(res, " ")
	if len(fields) != 2 {
		return fmt.Errorf("Get test mac address with: %s", res)
	}

	ap.Mac = strings.TrimSpace(fields[0])
	return nil
}

//AgtEthernetIpv6AddressPool GetTesterMacAddresses
func (ap *AddressPool) GetTesterMacAddress6() error {
	cmd := fmt.Sprintf("AgtEthernetIpv6AddressPool GetTesterMacAddresses %s", ap.Handler)
	res, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set tester mac address for ipv6 pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)

	fields := strings.Split(res, " ")
	if len(fields) != 2 {
		return fmt.Errorf("Get test mac address with: %s", res)
	}

	ap.Mac = strings.TrimSpace(fields[0])
	return nil
}

//GetTesterIpAddresses

func (ap *AddressPool) GetTesterIPAddresses() error {
	cmd := fmt.Sprintf("AgtEthernetAddressPool GetTesterIpAddresses %s", ap.Handler)
	res, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot get address pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	res = strings.Replace(res, "\"", "", -1)
	fields := strings.Split(res, " ")
	ap.First = strings.TrimSpace(fields[0])
	ap.Plen = strings.TrimSpace(fields[1])
	ap.Count = strings.TrimSpace(fields[2])
	ap.Step = strings.TrimSpace(fields[3])

	return nil
}

//AgtEthernetIpv6AddressPool GetTesterAddresses
func (ap *AddressPool) GetTesterIP6Addresses() error {
	cmd := fmt.Sprintf("AgtEthernetIpv6AddressPool GetTesterAddresses %s", ap.Handler)
	res, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot get ipv6 address pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	res = strings.Replace(res, "\"", "", -1)
	fields := strings.Split(res, " ")
	ap.First = strings.TrimSpace(fields[0])
	ap.Plen = strings.TrimSpace(fields[1])
	ap.Count = strings.TrimSpace(fields[2])
	ap.Step = strings.TrimSpace(fields[3])

	return nil
}

func (ap *AddressPool) SetVlan(vid string) error {
	cmd := fmt.Sprintf("AgtEthernetAddressPool EnableVlan %s", ap.Handler)
	_, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot enable vlan  for pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	cmd = fmt.Sprintf("AgtEthernetAddressPool SetVlanId %s %s", ap.Handler, strings.TrimSpace(vid))
	_, err = ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set vid  for pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}
	ap.Vlan = strings.TrimSpace(vid)

	return nil
}

func (ap *AddressPool) SetVlan6(vid string) error {
	//AgtEthernetIpv6AddressPool SetTesterVlanId
	cmd := fmt.Sprintf("AgtEthernetIpv6AddressPool SetTesterVlanId %s %s", ap.Handler, strings.TrimSpace(vid))
	_, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set vid  for ipv6 pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}
	ap.Vlan = strings.TrimSpace(vid)

	return nil
}

func (ap *AddressPool) DisableVlan() error {
	cmd := fmt.Sprintf("AgtEthernetAddressPool DisableVlan %s", ap.Handler)
	_, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot disable vlan  for pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	ap.Vlan = ""

	return nil
}

func (ap *AddressPool) DisableVlan6() error {
	cmd := fmt.Sprintf("AgtEthernetIpv6AddressPool SetTesterVlanId %s 0", ap.Handler)
	_, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot disable vlan  for pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	ap.Vlan = ""

	return nil
}

func (ap *AddressPool) EnableTrafficDestination() error {
	cmd := fmt.Sprintf("AgtEthernetAddressPool EnableTrafficDestinations %s", ap.Handler)
	_, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot enable traffic destination  for pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	return nil
}

func (ap *AddressPool) DisableTrafficDestination() error {
	cmd := fmt.Sprintf("AgtEthernetAddressPool DisableTrafficDestinations %s", ap.Handler)
	_, err := ap.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot disable traffic destination  for pool %s port: %s : %s", ap.Handler, ap.Port.Name, err.Error())
	}

	return nil
}
