package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gophergala/watchtower/interfaces/tcp"
)

func handleUDPPacket(packet []byte) error {
	return nil
}

func main() {
	// Flags
	var cores int
	var tcpPort int
	var udpPort int
	var secretKey string

	// Parse flags
	flag.StringVar(&secretKey, "secret", "supersecret", "The secret key clients must use to register")
	flag.IntVar(&cores, "cores", 1, "Number of cores (GOMAXPROCS) Watchtower can utilize")
	flag.IntVar(&tcpPort, "tcp", -1, "Port to host the TCP server on (defaults to OFF)")
	flag.IntVar(&udpPort, "udp", -1, "Port to host the UDP server on (defaults to OFF)")
	flag.Parse()

	// config.SetSecret(secretKey)

	if tcpPort != -1 {
		laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", tcpPort))
		if nil != err {
			panic(err)
		}

		listener, err := net.ListenTCP("tcp", laddr)
		if nil != err {
			panic(err)
		}

		go func() {
			// TODO: Clean TCP server shutdown
			for {
				conn, err := listener.Accept()
				if err != nil {
					// handle error
				}
				go tcp.Handle(conn)
			}
		}()
	}

	if udpPort != -1 {
		addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", udpPort))
		sock, _ := net.ListenUDP("udp", addr)

		// TODO: Clean UDP server shutdown
		go func() {
			for {
				buffer := make([]byte, 1024)
				plen, _, err := sock.ReadFromUDP(buffer)
				if err != nil {
					// Error reading what the client send - TODO: Handle this better
					panic(err)
				}

				go handleUDPPacket(buffer[0:plen])
			}
		}()
	}

	if tcpPort == -1 && udpPort == -1 {
		fmt.Println("you must active at least one endpoint (--http [PORT], --tcp [PORT] or --udp [PORT])")
		os.Exit(-1)
	}

	// Sleep while waiting for shutdown signals
	// TODO: Handle shutdown signals cleanly...
	for {
		time.Sleep(time.Second)
	}
}
