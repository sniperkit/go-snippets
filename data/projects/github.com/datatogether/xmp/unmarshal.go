package xmp

import (
	"encoding/xml"
)

// Unmarshal parses a stream of bytes into an XMPPacket
func Unmarshal(data []byte) (*XMPPacket, error) {
	packet := &XMPPacket{}
	err := xml.Unmarshal(data, &packet)
	if err != nil {
		return nil, err
	}

	return packet, err
}
