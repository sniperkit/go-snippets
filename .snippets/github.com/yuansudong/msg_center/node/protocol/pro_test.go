package protocol

import (
	"fmt"
	"testing"
)

func TestPacket(t *testing.T) {
	Data := Packet(1104, []byte("hellowrold"))
	cmd, d := UnPacket(Data)
	fmt.Println(cmd, string(d))
}
