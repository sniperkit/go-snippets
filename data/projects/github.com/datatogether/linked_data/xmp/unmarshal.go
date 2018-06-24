package xmp

import (
	"encoding/xml"
	"time"
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

func unmarshalXmpDate(date string) (time.Time, error) {
	return time.Parse(time.RFC3339, date)
}
