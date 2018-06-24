package client

import (
	"fmt"

	"github.com/juliengk/stack/client"
)

type Config struct {
	URL        *client.URL
	APIVersion string
}

func New(url, apiversion string) (*Config, error) {
	u, err := client.ParseUrl(url)
	if err != nil {
		return nil, err
	}

	return &Config{
		URL:        u,
		APIVersion: fmt.Sprintf("/metadata/%s", apiversion),
	}, nil
}

func (c *Config) buildPath(extra string) string {
	return fmt.Sprintf("%s%s", c.APIVersion, extra)
}

func (c *Config) get(path string) (string, error) {
	cc := &client.Config{
		Scheme: c.URL.Scheme,
		Host:   c.URL.Host,
		Port:   c.URL.Port,
		Path:   c.buildPath(path),
	}

	req, _ := client.New(cc)

	result := req.Get()
	if result.Error != nil {
		return "", result.Error
	}

	if result.Response.StatusCode != 200 {
		return "", fmt.Errorf("Problem fetching information")
	}

	return string(result.Body), nil
}

func (c *Config) GetAll() (string, error) {
	response, err := c.get("/")

	return response, err
}

func (c *Config) GetAllJSON() (string, error) {
	response, err := c.get(".json")

	return response, err
}

func (c *Config) GetID() (string, error) {
	response, err := c.get("/id")

	return response, err
}

func (c *Config) GetHostname() (string, error) {
	response, err := c.get("/hostname")

	return response, err
}

func (c *Config) GetFQDN() (string, error) {
	response, err := c.get("/fqdn")

	return response, err
}

func (c *Config) GetUserData() (string, error) {
	response, err := c.get("/user-data")

	return response, err
}

func (c *Config) GetVendorData() (string, error) {
	response, err := c.get("/vendor-data")

	return response, err
}

func (c *Config) GetPublicKeys() (string, error) {
	response, err := c.get("/public-keys")

	return response, err
}

func (c *Config) GetRegion() (string, error) {
	response, err := c.get("/region")

	return response, err
}

func (c *Config) GetNetworkInterfacesIndex() (string, error) {
	response, err := c.get("/interfaces/")

	return response, err
}

func (c *Config) GetNetworkInterfaceTypeIndex(itype string) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/", itype)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetNetworkEnumeratedInterfaceIndex(itype string, num int) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/%d/", itype, num)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetNetworkInterfaceMACAddress(itype string, num int) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/%d/mac", itype, num)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetNetworkInterfaceType(itype string, num int) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/%d/type", itype, num)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetNetworkInterfaceIPv4Index(itype string, num int) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/%d/ipv4/", itype, num)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetNetworkInterfaceIPv4Address(itype string, num int) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/%d/ipv4/address", itype, num)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetNetworkInterfaceIPv4Netmask(itype string, num int) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/%d/ipv4/netmask", itype, num)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetNetworkInterfaceIPv4Gateway(itype string, num int) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/%d/ipv4/gateway", itype, num)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetNetworkInterfaceFloatingIPIndex(itype string, num int) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/%d/floating_ip/", itype, num)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetNetworkInterfaceFloatingIPAddress(itype string, num int) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/%d/floating_ip/address", itype, num)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetNetworkInterfaceFloatingIPNetmask(itype string, num int) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/%d/floating_ip/netmask", itype, num)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetNetworkInterfaceFloatingIPGateway(itype string, num int) (string, error) {
	path := fmt.Sprintf("/interfaces/%s/%d/floating_ip/gateway", itype, num)
	response, err := c.get(path)

	return response, err
}

func (c *Config) GetDnsIndex() (string, error) {
	response, err := c.get("/dns/")

	return response, err
}

func (c *Config) GetDnsNameservers() (string, error) {
	response, err := c.get("/dns/nameservers")

	return response, err
}

func (c *Config) GetDnsSearchDomains() (string, error) {
	response, err := c.get("/dns/searchdomains")

	return response, err
}

func (c *Config) GetDnsOptions() (string, error) {
	response, err := c.get("/dns/options")

	return response, err
}

func (c *Config) GetTags() (string, error) {
	response, err := c.get("/tags")

	return response, err
}

func (c *Config) GetKeysIndex() (string, error) {
	response, err := c.get("keys/")

	return response, err
}

func (c *Config) GetKeyName(name string) (string, error) {
	path := fmt.Sprintf("/keys/%s", name)
	response, err := c.get(path)

	return response, err
}
