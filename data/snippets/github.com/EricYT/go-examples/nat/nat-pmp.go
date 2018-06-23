package main

import (
	"fmt"
	"log"
	"net"

	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

const internalPort int = 3233
const externalPort int = 3233
const lifeTime int = 10
const netType string = "tcp"

func listenInternalPort(port int) error {
	addr, err := net.ResolveTCPAddr(netType, fmt.Sprintf(":%d", internalPort))
	if err != nil {
		return err
	}

	log.Printf("[listenInternalPort] listen local address: %s", addr)
	listener, err := net.ListenTCP(netType, addr)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			go handleConnection(conn)
		}
	}()

	return nil
}

func handleConnection(conn net.Conn) {
	var buf []byte = make([]byte, 20)
	n, _ := conn.Read(buf)
	log.Printf("[handleConnection] read (%d) buffer: %s from remote: %s", n, string(buf), conn.RemoteAddr())
}

func sendPingMsg(targetIP net.IP, port int) error {
	target := fmt.Sprintf("%s:%d", targetIP, port)
	log.Printf("[sendPingMsg] dial target: %s", target)

	conn, err := net.Dial(netType, target)
	if err != nil {
		return err
	}
	log.Printf("[sendPingMsg] ready to send ping message to target: %s", target)
	n, err := conn.Write([]byte("hello,world"))
	if err != nil {
		log.Printf("[sendPingMsg] send ping message to target: %s error: %s", target, err)
		return err
	}
	defer conn.Close()
	log.Printf("[sendPingMsg] send ping message to target: %s success %d", target, n)
	return nil
}

func main() {
	// listen local internal port
	err := listenInternalPort(internalPort)
	if err != nil {
		panic(err)
	}

	//select {}

	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		panic(err)
	}

	client := natpmp.NewClient(gatewayIP)
	response, err := client.GetExternalAddress()
	if err != nil {
		panic(err)
	}
	log.Println("[main] External IP address:", response.ExternalIPAddress)
	extIPAddr := response.ExternalIPAddress

	externalIPAdd := net.IPv4(extIPAddr[0], extIPAddr[1], extIPAddr[2], extIPAddr[3])

	// add port mapping
	result, err := client.AddPortMapping(netType, internalPort, externalPort, lifeTime)
	if err != nil {
		panic(err)
	}
	log.Printf("[main] Add Port Mapping externalIPAddress: %s internalPort: %d externalPort: %d lifeTime: %d \n", externalIPAdd, internalPort, externalPort, lifeTime)

	mapExternalPort := result.MappedExternalPort
	log.Printf("[main] NAT mapping external port: %d", mapExternalPort)

	// send ping message
	//err = sendPingMsg(externalIPAdd, int(mapExternalPort))
	//if err != nil {
	//	panic(err)
	//}

	select {}
}
