package goxa

import (
	"log"
	"net"
)

func Connect(proto string, host string) (XA, error) {

	var xa XA
	xa.Queue = make(map[string][]byte)

	conn, err := net.Dial(proto, host)
	log.Printf("Connected to %s://%s.\n", proto, host)
	xa.ConnHandle = conn

	return xa, err
}

func Disconnect(xa XA) error {

	err := xa.ConnHandle.Close()
	log.Printf("Disconnected.\n")

	return err
}
