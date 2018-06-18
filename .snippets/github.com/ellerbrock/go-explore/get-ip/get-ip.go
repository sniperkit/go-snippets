package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)

	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			fmt.Println("IP: " + ipv4.String())
		}
	}
}
