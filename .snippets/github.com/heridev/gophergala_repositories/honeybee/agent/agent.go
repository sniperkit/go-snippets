package agent

import (
	"encoding/json"
	"fmt"
	"github.com/alicebob/procspy"
	"github.com/golang/protobuf/proto"
	"github.com/gophergala/honeybee/protobee"
	"net"
	"time"
)

type ConAux struct {
	Transport     string
	LocalAddress  string
	LocalPort     uint32
	RemoteAddress string
	RemotePort    uint32
	Pid           uint32
	Name          string
}

func sendDataToDest(data []byte, dst *string) {
	conn, err := net.Dial("tcp", *dst)
	if err != nil {
		fmt.Println("error", err)
	}
	n, err := conn.Write(data)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Printf("Sent %d bytes\n", n)
}

func listener(c <-chan []byte, dst *string) {
	for {
		msg := <-c
		var conAuxSlice []ConAux
		json.Unmarshal(msg, &conAuxSlice)

		fmt.Println("unmarshalled", conAuxSlice)

		connections := new(protobee.Connections)
		connections.Connection = []*protobee.Connection{}

		for _, value := range conAuxSlice {
			con := new(protobee.Connection)
			con.Transport = proto.String(value.Transport)
			con.LocalAddress = proto.String(value.LocalAddress)
			con.LocalPort = proto.Uint32(value.LocalPort)
			con.RemoteAddress = proto.String(value.RemoteAddress)
			con.RemotePort = proto.Uint32(value.RemotePort)
			con.Pid = proto.Uint32(value.Pid)
			con.Name = proto.String(value.Name)
			connections.Connection = append(connections.Connection, con)
		}
		//connections
		pb, err := proto.Marshal(connections)
		if err != nil {
			fmt.Println("error", err)
		}
		sendDataToDest(pb, dst)
		//time.Sleep(time.Second * 2)
	}
}

func startMonitor(channel chan<- []byte, scanningSeconds int64) {
	for {
		cs, err := procspy.Connections(true)
		if err != nil {
			panic(err)
		}
		js, err := json.Marshal(cs)
		if err != nil {
			fmt.Println(err)
		}
		channel <- js
		time.Sleep(time.Second * time.Duration(scanningSeconds))
	}

}

func Run() {
	fmt.Println("Start agent")

	var c chan []byte = make(chan []byte)

	dst := "127.0.0.1:2110"
	go listener(c, &dst)

	startMonitor(c, 1)
}
