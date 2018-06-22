package goxa

import (
	"log"
	"net"
)

func Listen(proto string, host string) (XA, error) {

	var xa XA
	xa.Queue = make(map[string][]byte)

	lstn, err := net.Listen(proto, host)
	log.Printf("Listened to %s://%s.\n", proto, host)

	if err == nil {

		xa.ListenHandle = lstn

		for {

			conn, err := lstn.Accept()

			if err == nil {
				log.Printf("Accepted connection.\n")
				xa.ConnHandle = conn

				return xa, err
			}
		}
	}

	return xa, err
}

func Close(xa XA) error {

	err := xa.ListenHandle.Close()
	log.Printf("Listener closed.\n")

	return err
}
