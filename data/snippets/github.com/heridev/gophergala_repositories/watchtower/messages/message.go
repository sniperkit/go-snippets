package messages

const (
	messageSenderSize  = 4
	messageChannelSize = 4
	messageTypeSize    = 2
	messageLengthSize  = 4
	messageHeaderSize  = messageSenderSize + messageChannelSize + messageTypeSize + messageLengthSize
)

// A Message is sent on a channel, either as
// a broadcast or directly to one or more users
// on the channel
type Message interface {
	Sender() uint32
	Channel() uint32
	Content() string
	Receivers() map[uint32]struct{} // nil if broadcast
	JSON() string
	Bytes() []byte
}
