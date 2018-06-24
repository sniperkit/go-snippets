package sqlite

import (
	"fmt"
	"net"

	iplib "github.com/juliengk/go-utils/ip"
	"github.com/kassisol/metadata/api/types"
)

func (c *Config) AddIP(ip, netmask, gateway string) error {
	version := 4
	if net.ParseIP(ip).To4() == nil {
		version = 6
	}

	itype := "public"
	if iplib.IsPrivateSubnet(net.ParseIP(ip)) {
		itype = "private"
	}

	c.DB.Create(&IP{
		IPAddress: ip,
		Netmask:   netmask,
		Gateway:   gateway,
		Version:   version,
		Type:      itype,
	})

	return nil
}

func (c *Config) RemoveIP(ip string) error {
	if c.memberOfInterface(ip) {
		return fmt.Errorf("IP \"%s\" cannot be removed. It is linked to an interface", ip)
	}

	c.DB.Where("ip_address = ?", ip).Delete(IP{})

	return nil
}

func (c *Config) ListIP(filter map[string]string) []types.IP {
	var result []types.IP

	sql := c.DB.Table("ips").Select("id, ip_address, netmask, gateway")

	if v, ok := filter["id"]; ok {
		sql = sql.Where("id = ?", v)
	}

	if v, ok := filter["ip"]; ok {
		sql = sql.Where("ip_address = ?", v)
	}

	if v, ok := filter["netmask"]; ok {
		sql = sql.Where("netmask = ?", v)
	}

	if v, ok := filter["gateway"]; ok {
		sql = sql.Where("gateway = ?", v)
	}

	rows, _ := sql.Rows()
	defer rows.Close()

	for rows.Next() {
		var d_id uint
		var d_ip string
		var d_netmask string
		var d_gateway string

		rows.Scan(&d_id, &d_ip, &d_netmask, &d_gateway)

		ir := types.IP{ID: d_id, IPAddress: d_ip, Netmask: d_netmask, Gateway: d_gateway}

		result = append(result, ir)
	}

	return result
}

func (c *Config) CountIP() int {
	var count int64

	c.DB.Model(&IP{}).Count(&count)

	return int(count)
}

func (c *Config) memberOfInterface(ip string) bool {
	var count int64

	c.DB.Table("profiles").Joins("JOIN ips ON ips.id = interfaces.ip_id").Where("ips.ip_address = ?", ip).Count(&count)

	if count > 0 {
		return true
	}

	return false
}
