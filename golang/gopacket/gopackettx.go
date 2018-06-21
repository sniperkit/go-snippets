package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

var (
	device       string = "ens33"
	snapshot_len int32  = 65535
	promiscuous  bool   = true
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle
	txhandle     *pcap.Handle
)

func main() {
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	/*
		txhandle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()
	*/

	// Set filter
	//var filter string = "tcp and port 22"
	var filter string = "proto ospf and not src host 10.71.1.122" //Try to capture ospf packet

	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	packetSource := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	for pkt := range packetSource.Packets() {
		// Do something with a packet here.
		fmt.Println(pkt)
		ethLayer := pkt.Layer(layers.LayerTypeEthernet)
		if ethLayer == nil {
			log.Println("Not an Ethernet frame")
			return
		}
		eth := ethLayer.(*layers.Ethernet)

		ipLayer := pkt.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			log.Println("Not an IP packet")
			return
		}

		ospfPkt := ipLayer.LayerPayload()
		ospfData := ospfPkt[24:]
		log.Printf("eth: %#v\n", eth)
		log.Printf("ip:  %#v\n", ipLayer)
		log.Printf("ospf: %#v\n", ospfPkt)
		log.Printf("ospfdata: %#v\n", ospfData)
		ethd := &layers.Ethernet{}
		ethd.DecodeFromBytes(pkt.Data(), nil)
		fmt.Printf(":::::::::::::::%#v\n", ethd)
		/*
			err := txhandle.WritePacketData(pkt.Data())
			if err != nil {
				log.Fatal(err)
			}
		*/
	}
}
