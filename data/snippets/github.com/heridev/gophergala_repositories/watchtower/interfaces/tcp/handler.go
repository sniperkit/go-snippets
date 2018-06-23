package tcp

import (
	"encoding/binary"
	"net"
	"sync"
	"time"

	"github.com/gophergala/watchtower/channels"
	"github.com/gophergala/watchtower/messages"
	"github.com/gophergala/watchtower/users"
)

type state int

const (
	stateNotConnected  state = 0
	stateRegistered    state = 1
	stateJoinedChannel state = 2
	stateDisconnected  state = 3

	errorMessageType       = byte('E')
	ackMessageType         = byte('A')
	registerMessageType    = byte('R')
	joinChannelMessageType = byte('J')
	sendMessageType        = byte('M')
)

var (
	bufferPool = sync.Pool{
		New: func() interface{} { return make([]byte, 1024) },
	}
)

func Handle(c net.Conn) error {
	var s state
	var u *user

	for {
		buffer := bufferPool.Get().([]byte) // Borrow a buffer from the pool
		// Read the incoming connection into the buffer.
		_, err := c.Read(buffer)
		if err != nil {
			return err
		}

		switch s {
		case stateNotConnected:
			// Waiting for the sender to register
			if buffer[0] == registerMessageType {
				// Register the user
				u = newUser()
				userID, err := users.Register(u)
				if err != nil {
					// Reply with an error message
					c.Write([]byte{errorMessageType})
				}

				// Reply with an ack + the user ID
				resp := make([]byte, 5, 5)
				resp[0] = ackMessageType
				binary.LittleEndian.PutUint32(resp[1:], userID)

				c.Write(resp)

				// Join channel
				s = stateJoinedChannel
			}

		case stateRegistered:
			// Waiting for the user to join a channel
			if buffer[0] == joinChannelMessageType {
				// Join the channel
				userID := binary.LittleEndian.Uint32(buffer[1:5])
				channelID := binary.LittleEndian.Uint32(buffer[5:9])
				channels.Join(userID, channelID)

				// Reply with an ack
				resp := make([]byte, 5, 5)
				resp[0] = ackMessageType
				binary.LittleEndian.PutUint32(resp[1:], channelID)
				c.Write(resp)

				s = stateJoinedChannel
			}

			time.Sleep(time.Second)

		case stateJoinedChannel:
			// Waiting for sender to send or receive messages (or join more channels)
			if buffer[0] == joinChannelMessageType {
				// Join the channel
				userID := binary.LittleEndian.Uint32(buffer[1:5])
				channelID := binary.LittleEndian.Uint32(buffer[5:9])
				channels.Join(userID, channelID)

				// Reply with an ack
				resp := make([]byte, 5, 5)
				resp[0] = ackMessageType
				binary.LittleEndian.PutUint32(resp[1:], channelID)
				c.Write(resp)
			}

			if buffer[0] == sendMessageType {
				// Find the channel
				userID := binary.LittleEndian.Uint32(buffer[1:5])
				channelID := binary.LittleEndian.Uint32(buffer[5:9])
				messageLength := binary.LittleEndian.Uint32(buffer[9:13])
				message := buffer[13:(13 + int(messageLength))]

				// Send message to channel
				m := messages.NewBroadcastMessage(userID, channelID, string(message))
				channels.Send(m, channelID)
			}

			select {
			case msg := <-u.messageQueue:
				var resp []byte
				resp = append(resp, sendMessageType)
				resp = append(resp, msg.Bytes()...)
				c.Write(resp)
			default:
			}

		case stateDisconnected:
			// Clean up connection and exit
			return c.Close()

		default:
			panic("invalid state of tcp connection, bailing out")
		}

		bufferPool.Put(buffer)            // Put the buffer back in the pool
		time.Sleep(time.Millisecond * 10) // no need to check in constantly
	}
}
