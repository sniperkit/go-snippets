package ip

import (
	"net"
)

type IP struct {
        V4      []string
        V6      []string
}

type Interface struct {
	net.Interface
	IP
}

type Interfaces map[string]*Interface

func New() Interfaces {
        return make(Interfaces)
}

func (intfs Interfaces) Get() error {
        ifaces, err := net.Interfaces()
        if err != nil {
		return err
        }

        for _, i := range ifaces {
                addrs, err := i.Addrs()
                if err != nil {
			return err
                }

                for _, addr := range addrs {
                        var ip net.IP

                        switch v := addr.(type) {
                        case *net.IPNet:
                                ip = v.IP
                        case *net.IPAddr:
                                ip = v.IP
                        }

                        intfs.add(i, ip)
                }
        }

	return nil
}

func (intfs Interfaces) GetIntf(intf string) *Interface {
	return intfs[intf]
}

func (intfs Interfaces) IgnoreIntf(ifaces []string) {
	for _, iface := range ifaces {
		if _, ok := intfs[iface]; ok {
			delete(intfs, iface)
		}
	}
}

func (intfs Interfaces) add(iface net.Interface, ip net.IP) {
	ipver := 4
	if !IsIPv4(ip) {
		ipver = 6
	}

        if val, ok := intfs[iface.Name]; ok {
		if ipver == 4 {
			val.V4 = append(val.V4, ip.String())

		} else {
			val.V6 = append(val.V6, ip.String())
		}
        } else {
		if ipver == 4 {
			intfs[iface.Name] = &Interface{iface, IP{V4: []string{ip.String()}}}
		} else {
			intfs[iface.Name] = &Interface{iface, IP{V6: []string{ip.String()}}}
		}
	}
}

func IsIPv4(ip net.IP) bool {
        if ip.To4() == nil {
                return false
        }

        return true
}
