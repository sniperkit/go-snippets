package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
	"log"
	"time"
)

var (
	device       string = "ens33"
	snapshot_len int32  = 65549
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle
	txhandle       *pcap.Handle
)

func main() {
	recvHdl, err := pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if recvHdl == nil {
		log.Printf("RecvHdl: No device found. %s, %s", device, err.Error())
		return
	}

	var filter string = "proto ospf and not src host 10.71.1.122"
	// Setting Pcap filter for Ospf Pkt
	err = recvHdl.SetBPFFilter(filter)
	if err != nil {
		log.Printf("Unable to set filter on %s", device)
		return
	}

	recv := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := recv.Packets()
	for {   
		select {
		case packet, ok := <-in:
			if ok {
				//server.logger.Info("Got Some Ospf Packet on the Recv Thread")
				log.Printf("Received on packet %#v", packet)
			}
		}
	}
}
