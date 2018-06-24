package metadata

import (
	"net/http"
	"strings"

	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/api/types"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/labstack/echo"
)

func AllHandle(c echo.Context) error {
	result := []string{
		"id",
		"hostname",
		"fqdn",
		"user-data",
		"vendor-data",
		"public-keys",
		"region",
		"interfaces/",
		"dns/",
		"tags/",
	}

	if c.Request().Method == "HEAD" {
		return c.NoContent(http.StatusOK)
	}

	return c.String(http.StatusOK, strings.Join(result, "\n"))
}

func AllJsonHandle(c echo.Context) error {
	serverid := c.Get("SERVERID").(int)

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		return err
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		return err
	}
	defer s.End()

	id := s.GetID(serverid)
	hostname := s.GetHostname(serverid)
	fqdn := s.GetFQDN(serverid)
	userData := s.GetUserData(serverid)
	vendorData := s.GetVendorData(serverid)
	publicKeys := s.GetPublicKeys(serverid)
	region := s.GetRegion(serverid)

	interfaces := make(map[string][]map[string]interface{})

	nicTypes := s.GetInterfaces(serverid)
	for _, itype := range nicTypes {
		index := s.GetInterfacesType(serverid, itype)
		for _, i := range index {
			mac := s.GetInterfaceMACAddress(serverid, itype, i)
			nicType := s.GetInterfaceType(serverid, itype, i)

			mMap := map[string]interface{}{
				"mac":  mac,
				"type": nicType,
			}
			r := s.GetEnumeratedInterface(serverid, itype, i)
			for _, ipver := range r {
				if ipver == "ipv4" {
					ip4 := map[string]string{
						"ip_address": s.GetInterfaceIPv4Address(serverid, itype, i),
						"netmask":    s.GetInterfaceIPv4Netmask(serverid, itype, i),
						"gateway":    s.GetInterfaceIPv4Gateway(serverid, itype, i),
					}

					mMap["ipv4"] = ip4
				}
			}

			if s.FloatingIPExists(serverid, itype, i) {
				floatingip := map[string]string{
					"ip_address": s.GetInterfaceFloatingIPAddress(serverid, itype, i),
					"netmask":    s.GetInterfaceFloatingIPNetmask(serverid, itype, i),
					"gateway":    s.GetInterfaceFloatingIPGateway(serverid, itype, i),
				}

				mMap["floating_ip"] = floatingip
			}

			interfaces[itype] = append(interfaces[itype], mMap)
		}
	}

	d := s.GetDNSIndex(serverid)
	dns := make(map[string][]string)

	for _, elem := range d {
		if elem == "nameservers" {
			dns["nameservers"] = s.GetDNSNameservers(serverid)
		}
		if elem == "searches" {
			dns["searches"] = s.GetDNSSearchDomains(serverid)
		}
		if elem == "options" {
			dns["options"] = s.GetDNSOptions(serverid)
		}
	}

	tags := s.GetTags(serverid)

	cloudinit := types.CloudInit{
		ID:       id,
		Hostname: hostname,
		FQDN:     fqdn,
	}

	if len(userData) > 0 {
		cloudinit.UserData = userData
	}
	if len(vendorData) > 0 {
		cloudinit.VendorData = vendorData
	}
	if len(publicKeys) > 0 {
		cloudinit.PublicKeys = publicKeys
	}
	if len(region) > 0 {
		cloudinit.Region = region
	}
	if len(nicTypes) > 0 {
		cloudinit.Interfaces = interfaces
	}
	if len(dns) > 0 {
		cloudinit.DNS = dns
	}
	if len(tags) > 0 {
		cloudinit.Tags = tags
	}

	if c.Request().Method == "HEAD" {
		return c.NoContent(http.StatusOK)
	}

	return c.JSON(http.StatusOK, &cloudinit)
}
