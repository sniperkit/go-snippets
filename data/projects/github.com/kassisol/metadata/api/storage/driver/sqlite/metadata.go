package sqlite

import (
	"github.com/juliengk/go-utils"
	"github.com/kassisol/metadata/api/storage/driver"
)

func (c *Config) GetIDFromIP(ip string) int {
	var result int

	row := c.DB.Table("hosts").Select("hosts.id").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.ip_id").Where("ips.ip_address = ?", ip).Row()

	row.Scan(&result)

	return result
}

func (c *Config) GetID(srvid int) int {
	h := Host{}

	c.DB.Where("id = ?", srvid).First(&h)

	return int(h.ID)
}

func (c *Config) GetHostname(srvid int) string {
	h := Host{}

	c.DB.Where("id = ?", srvid).First(&h)

	return h.Name
}

func (c *Config) GetFQDN(srvid int) string {
	h := Host{}

	c.DB.Where("id = ?", srvid).First(&h)

	return h.FQDN
}

func (c *Config) GetUserData(srvid int) string {
	return c.getValue(srvid, "user-data")
}

func (c *Config) GetVendorData(srvid int) string {
	return c.getValue(srvid, "vendor-data")
}

func (c *Config) GetPublicKeys(srvid int) []string {
	return c.getValues(srvid, "public-key")
}

func (c *Config) GetRegion(srvid int) string {
	return c.getValue(srvid, "region")
}

func (c *Config) GetInterfaces(srvid int) []string {
	var result []string

	rows, _ := c.DB.Table("hosts").Select("DISTINCT(ips.type)").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.ip_id").Where("hosts.id = ?", srvid).Rows()
	defer rows.Close()

	for rows.Next() {
		var itype string

		rows.Scan(&itype)

		result = append(result, itype)
	}

	return result
}

func (c *Config) GetInterfacesType(srvid int, itype string) []int {
	var result []int

	rows, _ := c.DB.Table("hosts").Select("interfaces.`index`").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.ip_id").Where("hosts.id = ?", srvid).Where("ips.type = ?", itype).Rows()
	defer rows.Close()

	for rows.Next() {
		var index int

		rows.Scan(&index)

		result = append(result, index)
	}

	return result
}

func (c *Config) GetEnumeratedInterface(srvid int, itype string, index int) []string {
	var result []string

	rows, _ := c.DB.Table("hosts").Select("DISTINCT(ips.version)").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.ip_id").Where("hosts.id = ?", srvid).Where("ips.type = ?", itype).Where("interfaces.`index` = ?", index).Rows()
	defer rows.Close()

	for rows.Next() {
		var version int

		rows.Scan(&version)

		if version == 4 {
			result = append(result, "ipv4")
		}
		if version == 6 {
			result = append(result, "ipv6")
		}
	}

	return result
}

func (c *Config) GetInterfaceMACAddress(srvid int, itype string, index int) string {
	var result string

	row := c.DB.Table("hosts").Select("interfaces.mac_address").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.ip_id").Where("hosts.id = ?", srvid).Where("ips.type = ?", itype).Where("interfaces.`index` = ?", index).Row()

	row.Scan(&result)

	return result
}

func (c *Config) GetInterfaceType(srvid int, itype string, index int) string {
	var result string

	row := c.DB.Table("hosts").Select("ips.type").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.ip_id").Where("hosts.id = ?", srvid).Where("ips.type = ?", itype).Where("interfaces.`index` = ?", index).Row()

	row.Scan(&result)

	return result
}

func (c *Config) GetInterfaceIPv4Address(srvid int, itype string, index int) string {
	var result string

	row := c.DB.Table("hosts").Select("ips.ip_address").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.ip_id").Where("hosts.id = ?", srvid).Where("ips.type = ?", itype).Where("interfaces.`index` = ?", index).Row()

	row.Scan(&result)

	return result
}

func (c *Config) GetInterfaceIPv4Netmask(srvid int, itype string, index int) string {
	var result string

	row := c.DB.Table("hosts").Select("ips.netmask").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.ip_id").Where("hosts.id = ?", srvid).Where("ips.type = ?", itype).Where("interfaces.`index` = ?", index).Row()

	row.Scan(&result)

	return result
}

func (c *Config) GetInterfaceIPv4Gateway(srvid int, itype string, index int) string {
	var result string

	row := c.DB.Table("hosts").Select("ips.gateway").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.ip_id").Where("hosts.id = ?", srvid).Where("ips.type = ?", itype).Where("interfaces.`index` = ?", index).Row()

	row.Scan(&result)

	return result
}

func (c *Config) FloatingIPExists(srvid int, itype string, index int) bool {
	var count int64

	c.DB.Table("hosts").Select("ips.ip_address").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.floating_ip_id").Where("hosts.id = ?", srvid).Where("ips.type = ?", itype).Where("interfaces.`index` = ?", index).Count(&count)

	if count > 0 {
		return true
	}

	return false
}

func (c *Config) GetInterfaceFloatingIPAddress(srvid int, itype string, index int) string {
	var result string

	row := c.DB.Table("hosts").Select("ips.ip_address").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.floating_ip_id").Where("hosts.id = ?", srvid).Where("ips.type = ?", itype).Where("interfaces.`index` = ?", index).Row()

	row.Scan(&result)

	return result
}

func (c *Config) GetInterfaceFloatingIPNetmask(srvid int, itype string, index int) string {
	var result string

	row := c.DB.Table("hosts").Select("ips.netmask").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.floating_ip_id").Where("hosts.id = ?", srvid).Where("ips.type = ?", itype).Where("interfaces.`index` = ?", index).Row()

	row.Scan(&result)

	return result
}

func (c *Config) GetInterfaceFloatingIPGateway(srvid int, itype string, index int) string {
	var result string

	row := c.DB.Table("hosts").Select("ips.gateway").Joins("JOIN host_interfaces ON host_interfaces.host_id = hosts.id").Joins("JOIN interfaces ON interfaces.id = host_interfaces.interface_id").Joins("JOIN ips ON ips.id = interfaces.floating_ip_id").Where("hosts.id = ?", srvid).Where("ips.type = ?", itype).Where("interfaces.`index` = ?", index).Row()

	row.Scan(&result)

	return result
}

func (c *Config) GetDNSIndex(srvid int) []string {
	var result []string

	rows, _ := c.DB.Table("hosts").Select("DISTINCT(data.type)").Joins("JOIN profiles ON profiles.id = hosts.profile_id").Joins("JOIN profile_datas ON profile_datas.profile_id = profiles.id").Joins("JOIN data ON data.id = profile_datas.data_id").Where("hosts.id = ?", srvid).Where("data.type LIKE 'dns-%'").Rows()
	defer rows.Close()

	for rows.Next() {
		var dtype string

		rows.Scan(&dtype)

		if dtype == "dns-nameserver" {
			result = append(result, "nameservers")
		}
		if dtype == "dns-searchdomain" {
			result = append(result, "searchdomains")
		}
		if dtype == "dns-option" {
			result = append(result, "options")
		}
	}

	return result
}

func (c *Config) GetDNSNameservers(srvid int) []string {
	return c.getValues(srvid, "dns-nameserver")
}

func (c *Config) GetDNSSearchDomains(srvid int) []string {
	return c.getValues(srvid, "dns-searchdomain")
}

func (c *Config) GetDNSOptions(srvid int) []string {
	return c.getValues(srvid, "dns-option")
}

func (c *Config) GetTags(srvid int) []string {
	return c.getValues(srvid, "tag")
}

func (c *Config) GetKeys(srvid int) []string {
	var result []string

	rows, _ := c.DB.Table("hosts").Select("DISTINCT(data.type)").Joins("JOIN profiles ON profiles.id = hosts.profile_id").Joins("JOIN profile_datas ON profile_datas.profile_id = profiles.id").Joins("JOIN data ON data.id = profile_datas.data_id").Where("hosts.id = ?", srvid).Rows()
	defer rows.Close()

	for rows.Next() {
		var key string

		rows.Scan(&key)

		if !utils.StringInSlice(key, driver.DefaultKeys, false) {
			result = append(result, key)
		}
	}

	return result
}

func (c *Config) GetKey(srvid int, key string) string {
	return c.getValue(srvid, key)
}

func (c *Config) getValue(srvid int, dtype string) string {
	var result string

	row := c.DB.Table("hosts").Select("data.value").Joins("JOIN profiles ON profiles.id = hosts.profile_id").Joins("JOIN profile_datas ON profile_datas.profile_id = profiles.id").Joins("JOIN data ON data.id = profile_datas.data_id").Where("hosts.id = ?", srvid).Where("data.type = ?", dtype).Row()

	row.Scan(&result)

	return result
}

func (c *Config) getValues(srvid int, dtype string) []string {
	var result []string

	rows, _ := c.DB.Table("hosts").Select("data.value").Joins("JOIN profiles ON profiles.id = hosts.profile_id").Joins("JOIN profile_datas ON profile_datas.profile_id = profiles.id").Joins("JOIN data ON data.id = profile_datas.data_id").Where("hosts.id = ?", srvid).Where("data.type = ?", dtype).Rows()
	defer rows.Close()

	for rows.Next() {
		var value string

		rows.Scan(&value)

		result = append(result, value)
	}

	return result
}
