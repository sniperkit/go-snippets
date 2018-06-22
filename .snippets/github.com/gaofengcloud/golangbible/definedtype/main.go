package main

import (
	"fmt"
	"strings"
)

type Server struct {
	Name string
}

type Servers []Server

func ListServers() Servers {
	return []Server{
		{Name: "Server 1"},
		{Name: "Server 2"},
		{Name: "Host 1"},
		{Name: "Host 2"},
	}
}

func (s Servers) Filter(name string) Servers {
	filterd := make(Servers, 0)

	for _, server := range s {
		if strings.Contains(server.Name, name) {
			filterd = append(filterd, server)
		}
	}

	return filterd
}

func main() {
	servers := ListServers()
	filterd := servers.Filter("Host")
	fmt.Printf("Filtered servers: %v", filterd)
}
