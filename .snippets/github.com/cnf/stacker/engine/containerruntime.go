package engine

import (
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

// ContainerRuntime is the interface for the container runtine object.
type ContainerRuntime interface {
	Setup
	StartStop
	Commander
}

// Container is the interface for containers.
type Container interface {
	Commander
}

// ContainerData holds information about a container for communicating with modules.
type ContainerData struct {
	ID      string
	Name    string
	Created time.Time
	// Args    []string
	Network *ContainerNetwork
}

// NewContainerData returns an empty ContainerData object.
func NewContainerData() *ContainerData {
	return &ContainerData{
		Network: &ContainerNetwork{
			Published: make(map[Port][]Published),
		},
	}
}

// ContainerNetwork holds network information about a container.
type ContainerNetwork struct {
	IPAddress  string
	Mask       int
	Ports      []Port
	Published  map[Port][]Published
	Hostname   string
	Domainname string
	DNS        []string
}

// Port is a port number.
// the format is fe 80/tcp or 53/udp
type Port string

// Port returns the port number as an integer.
func (p Port) Port() uint16 {
	i, err := strconv.ParseUint(strings.Split(string(p), "/")[0], 10, 16)
	if err != nil {
		return 0
	}
	return uint16(i)
}

// Proto returns the protocol as a string.
// at this time, only tcp and udp are supported.
func (p Port) Proto() string {
	parts := strings.Split(string(p), "/")
	if len(parts) == 1 {
		return "tcp"
	}
	switch strings.ToLower(parts[1]) {
	case "tcp":
		return "tcp"
	case "udp":
		return "udp"
	}
	return "tcp"
}

// Published is a published port, with the ip and port it is published to.
type Published struct {
	IP   net.IP
	Port uint16
}

// NewContainerRuntime is a function definition each Docker must provide
type NewContainerRuntime func() (c ContainerRuntime)

var containerRTList = make(map[string]NewContainerRuntime)

// RegisterContainerRuntime is called to register a ContainerRuntime object
func RegisterContainerRuntime(name string, creator NewContainerRuntime) {
	name = strings.ToLower(name)
	if glog.V(1) {
		glog.Infof("registering Container runtime: %s", name)
	}

	if containerRTList[name] != nil {
		glog.Fatalf("a Container runtime named `%s` already exists", name)
	}
	containerRTList[name] = creator
}
