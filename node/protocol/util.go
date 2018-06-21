package protocol

import (
	"bytes"
	"encoding/binary"
)

// UnPacket is usedas unpacket the data from remote computer
func UnPacket(RequestData []byte) (cmd int, data []byte) {
	if len(RequestData) < 4 {
		return GUnknownProtoCode, data
	}
	cmd = bytesToInt(RequestData[0:4])
	if cmd < 0 {
		return GUnknownProtoCode, data
	}
	data = RequestData[4:]
	return cmd, data
}

// Packet is used as make cmd and request
func Packet(cmd int, requestContent []byte) []byte {
	tmp := uint32(cmd)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	bytesBuffer.Write(requestContent)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func bytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}
