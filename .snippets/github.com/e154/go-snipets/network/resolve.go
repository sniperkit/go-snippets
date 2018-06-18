package main

import (
	"net"
	"fmt"
)

func main() {

	// var 1
	ip, err := net.ResolveIPAddr("ip4", "www.google.com")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(ip)

	// var 2
	ips, err := net.LookupIP("www.google.com")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(ips)

}