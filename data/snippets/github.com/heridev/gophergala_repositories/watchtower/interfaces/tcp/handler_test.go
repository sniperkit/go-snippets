package tcp

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestHandler(t *testing.T) {
	laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:48879")
	if nil != err {
		t.Error(err)
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if nil != err {
		t.Error(err)
	}

	go func() {
		for {
			// Listen for connections and let Handler handle them
			conn, err := listener.Accept()
			if err != nil {
				t.Error(err)
			}
			go Handle(conn)
		}
	}()

	c, err := net.Dial("tcp", "127.0.0.1:48879")
	if err != nil {
		t.Error(err)
	}

	c.Write([]byte{registerMessageType})
	buffer := bufferPool.Get().([]byte)
	bytes, err := c.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	if buffer[0] != ackMessageType || bytes != 5 {
		t.Error("invalid response")
	}

	// Store the user ID
	userID := buffer[1:5]

	// Let's join channel 1
	req := make([]byte, 9, 9)
	req[0] = joinChannelMessageType            // Request a channel join
	copy(req[1:5], userID)                     // Copy the user ID here
	binary.LittleEndian.PutUint32(req[5:9], 1) // And set the channel ID we want to join
	c.Write(req)

	// Re-use the buffer for the response
	bytes, err = c.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	// Check that we got the correct ACK
	if buffer[0] != ackMessageType || bytes != 5 || binary.LittleEndian.Uint32(buffer[1:5]) != 1 {
		t.Error("invalid response")
	}

	// Then, send a message
	req = make([]byte, 30, 30)
	req[0] = sendMessageType                     // We want to SEND a message now
	copy(req[1:5], userID)                       // Set the user ID again
	binary.LittleEndian.PutUint32(req[5:9], 1)   // And set the channel ID we want to join
	binary.LittleEndian.PutUint32(req[9:13], 10) // And set the channel ID we want to join
	copy(req[13:23], []byte("helloworld"))       // And include a message
	c.Write(req)                                 // And send the request

	// Re-use the buffer for the response
	bytes, err = c.Read(buffer)
	if err != nil {
		t.Error(err)
	}

	// Check that we got the correct ACK
	if buffer[0] != sendMessageType {
		t.Error("invalid response")
	}

	if binary.LittleEndian.Uint32(buffer[1:5]) != binary.LittleEndian.Uint32(userID) {
		t.Error("invalid user id in message")
	}

	if binary.LittleEndian.Uint32(buffer[5:9]) != 1 {
		t.Error("invalid channel for message")
	}

	if binary.LittleEndian.Uint16(buffer[9:11]) != 0 {
		t.Error("invalid use of reserved space")
	}

	if binary.LittleEndian.Uint32(buffer[11:15]) != 10 {
		t.Error("invalid length of message")
	}

	if string(buffer[15:bytes]) != "helloworld" {
		t.Error("invalid message")
	}
}
