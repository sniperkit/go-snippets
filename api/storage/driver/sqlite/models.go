package sqlite

import (
	"time"
)

type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"created_at"`
}

type ServerConfig struct {
	Model
	Key   string
	Value string
}

type Data struct {
	Model
	Name        string `gorm:"unique;"`
	Type        string
	Value       string
	Description string
}

type Profile struct {
	Model
	Name  string `gorm:"unique;"`
	Datas []Data `gorm:"many2many:profile_datas;"`
}

type IP struct {
	Model
	IPAddress string `gorm:"unique;"`
	Netmask   string
	Gateway   string
	Version   int
	Type      string
}

type Interface struct {
	Model
	Index        int
	MACAddress   string `gorm:"unique;"`
	IP           IP
	IPID         uint
	FloatingIP   IP
	FloatingIPID uint
}

type Host struct {
	Model
	IsEnabled  bool
	Name       string      `gorm:"unique;"`
	FQDN       string      `gorm:"unique;"`
	Interfaces []Interface `gorm:"many2many:host_interfaces;"`
	Profile    Profile
	ProfileID  uint
}
