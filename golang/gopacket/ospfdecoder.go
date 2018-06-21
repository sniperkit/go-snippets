package main

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

var (
	DecodeOSPFPktHeader gopacket.Decoder = gopacket.DecodeFunc(decodeOSPFPktHeader)
	DecodeOSPFPktHello  gopacket.Decoder = gopacket.DecodeFunc(decodeOSPFPktHello)
	DecodeOSPFPktDBD    gopacket.Decoder = gopacket.DecodeFunc(decodeOSPFPktDBD)
	DecodeOSPFPktLSR    gopacket.Decoder = gopacket.DecodeFunc(decodeOSPFPktLSR)
	DecodeOSPFPktLSU    gopacket.Decoder = gopacket.DecodeFunc(decodeOSPFPktLSU)
	DecodeOSPFPktLSAck  gopacket.Decoder = gopacket.DecodeFunc(decodeOSPFPktLSAck)
)

var (
	LayerTypeOSPFHeader = gopacket.RegisterLayerType(666, gopacket.LayerTypeMetadata{Name: "OSPF Header", Decoder: DecodeOSPFPktHeader})
	LayerTypeOSPFHello  = gopacket.RegisterLayerType(667, gopacket.LayerTypeMetadata{Name: "OSPF Hello", Decoder: DecodeOSPFPktHello})
	LayerTypeOSPFDBD    = gopacket.RegisterLayerType(668, gopacket.LayerTypeMetadata{Name: "OSPF Database Description", Decoder: DecodeOSPFPktDBD})
	LayerTypeOSPFLSR    = gopacket.RegisterLayerType(669, gopacket.LayerTypeMetadata{Name: "OSPF Link State Request", Decoder: DecodeOSPFPktLSR})
	LayerTypeOSPFLSU    = gopacket.RegisterLayerType(670, gopacket.LayerTypeMetadata{Name: "OSPF Link State Update", Decoder: DecodeOSPFPktLSU})
	LayerTypeOSPFLSAck  = gopacket.RegisterLayerType(671, gopacket.LayerTypeMetadata{Name: "OSPF Link State Acknowledgment", Decoder: DecodeOSPFPktLSAck})
)

func decodeOSPFPktHeader(data []byte, p gopacket.PacketBuilder) error {
	return errors.New("Layer type not currently supported")
}

func decodeOSPFPktHello(data []byte, p gopacket.PacketBuilder) error {
	return errors.New("Layer type not currently supported")
}

func decodeOSPFPktDBD(data []byte, p gopacket.PacketBuilder) error {
	return errors.New("Layer type not currently supported")
}

func decodeOSPFPktLSR(data []byte, p gopacket.PacketBuilder) error {
	return errors.New("Layer type not currently supported")
}

func decodeOSPFPktLSU(data []byte, p gopacket.PacketBuilder) error {
	return errors.New("Layer type not currently supported")
}

func decodeOSPFPktLSAck(data []byte, p gopacket.PacketBuilder) error {
	return errors.New("Layer type not currently supported")
}

type Header struct {
}

func (h *Header) LayerType() gopacket.LayerType { return LayerTypeOSPFHeader }
func (h *Header) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	return nil
}

func (h *Header) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	return nil
}

func (h *Header) CanDecode() gopacket.LayerClass {
	return LayerTypeOSPFHeader
}

type Hello struct {
}

func (h *Hello) LayerType() gopacket.LayerType { return LayerTypeOSPFHello }

func (h *Hello) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	return nil
}

func (h *Hello) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	return nil
}

func (h *Hello) CanDecode() gopacket.LayerClass {
	return LayerTypeOSPFHello
}

type DBD struct {
}

func (d *DBD) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	return nil
}

func (d *DBD) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	return nil
}

func (d *DBD) CanDecode() gopacket.LayerClass {
	return LayerTypeOSPFHeader
}

func (d *DBD) LayerType() gopacket.LayerType { return LayerTypeOSPFDBD }

type LSR struct {
}

func (l *LSR) LayerType() gopacket.LayerType { return LayerTypeOSPFLSR }

func (l *LSR) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	return nil
}

func (l *LSR) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	return nil
}

func (l *LSR) CanDecode() gopacket.LayerClass {
	return LayerTypeOSPFHello
}

type LSU struct {
}

func (l *LSU) LayerType() gopacket.LayerType { return LayerTypeOSPFLSU }

func (l *LSU) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	return nil
}

func (l *LSU) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	return nil
}

func (l *LSU) CanDecode() gopacket.LayerClass {
	return LayerTypeOSPFHello
}

type LSAck struct {
}

func (l *LSAck) LayerType() gopacket.LayerType { return LayerTypeOSPFLSAck }

func (l *LSAck) DecodeFromBytes(data []byte, df gopacket.DecodeFeedback) error {
	return nil
}

func (l *LSAck) SerializeTo(b gopacket.SerializeBuffer, opts gopacket.SerializeOptions) error {
	return nil
}

func (l *LSAck) CanDecode() gopacket.LayerClass {
	return LayerTypeOSPFHello
}

func main() {
}
