package types

type ServerConfig struct {
	Key   string
	Value string
}

type Data struct {
	ID          uint
	Name        string
	Type        string
	Value       string
	Description string
}

type IP struct {
	ID        uint
	IPAddress string
	Netmask   string
	Gateway   string
}

type Interface struct {
	Index      int
	MACAddress string
	IP         IP
	FloatingIP IP
}

type Host struct {
	ID         uint
	Enabled    bool
	Name       string
	FQDN       string
	UUID       string
	Interfaces []string
	Profile    string
}

type CloudInit struct {
	ID         int                 `json:"id"`
	Hostname   string              `json:"hostname"`
	FQDN       string              `json:"fqdn"`
	UserData   string              `json:"user_data"`
	VendorData string              `json:"vendor_data"`
	PublicKeys []string            `json:"public_keys"`
	Region     string              `json:"region"`
	Interfaces interface{}         `json:"interfaces"`
	DNS        map[string][]string `json:"dns"`
	Tags       []string            `json:"tags"`
}
