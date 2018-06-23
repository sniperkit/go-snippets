package main

import (
	"command"
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"os"
	"rut"
)

var (
	handle *pcap.Handle
	err    error
)

var CTX = context.Background()

func main() {

	/*
		dev, err := rut.New(&rut.RUT{
			Name:     "SWITCH",
			Device:   "V8",
			IP:       "10.71.20.191",
			Port:     "23",
			Username: "admin",
			Protocol: "telnet",
			Hostname: "V8500_2",
			Password: "",
			SFU:      "A",
		})
	*/
	dev, err := rut.New(&rut.RUT{
		Name:     "SWITCH",
		Device:   "V5",
		IP:       "10.71.20.55",
		Port:     "23",
		Username: "admin",
		Protocol: "telnet",
		Hostname: "SWITCH",
		Password: "",
		SFU:      "A",
	})

	if err != nil {
		panic(err)
	}

	err = dev.Init()
	if err != nil {
		panic(err)
	}

	_, err = dev.RunCommand(CTX, &command.Command{
		Mode: "normal",
		CMD:  " config terminal",
	})
	if err != nil {
		fmt.Println(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dev.TCPDUMP("eth05", "", "2.pcap", "")
	fmt.Println("Dump finished")

	if err = dev.FTPPut("/etc/2.pcap", "10.71.1.161", "pi", "raspberry", pwd); err != nil {
		panic(err)
	}

	handle, err = pcap.OpenOffline("2.pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Loop through packets in file
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println(packet)
	}
}
