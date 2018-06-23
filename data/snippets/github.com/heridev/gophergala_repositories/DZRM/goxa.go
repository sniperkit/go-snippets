package goxa

import (
	"net"
)

const UOW string = "UOW "
const CMIT string = "CMIT"
const BACK string = "BACK"
const UUIDLength int = 36

type XA struct {
	ListenHandle net.Listener
	ConnHandle   net.Conn
	Queue        map[string][]byte
	ID           string
}
