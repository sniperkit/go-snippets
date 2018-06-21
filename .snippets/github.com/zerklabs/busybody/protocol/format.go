package protocol

import "fmt"

const (
	HelloMessage     int = 0
	PingMessage      int = 1
	PingReqMessage   int = 2
	PingReplyMessage int = 3
	PingRelayMessage int = 4
	StandardMessage  int = 5
)

const (
	NoCompression      int = 0
	SnappyCompression  int = 1
	DeflateCompression int = 2
	ZlibCompression    int = 3
)

//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// /                                                               /
// \                           Header                              \
// /                                                               /
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// /                                                               /
// \                           Content                             \
// /                                                               /
// |                                                               |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//

func NewMessage(msgtype int, comptype int, id string) *Message {
	return &Message{
		Header: buildMessageHeader(msgtype, comptype, id),
		buf:    make([]byte, 0),
		off:    0,
	}
}

func Decode(msg []byte) (*Message, error) {
	if len(msg) == 0 {
		return nil, fmt.Errorf("empty message")
	}

	n, header, err := rebuildHeader(msg)
	if err != nil {
		return nil, err
	}

	protocol := &Message{
		Header: header,
		buf:    make([]byte, 0),
		off:    0,
	}

	// n also accounts for the NULSEP byte sequence
	if len(msg) > n {
		_, err := protocol.Write(msg[n:])
		if err != nil {
			return nil, err
		}
	}

	// return protocol, nil
	return protocol, nil
}
