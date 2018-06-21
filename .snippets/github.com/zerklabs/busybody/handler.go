package busybody

import "github.com/zerklabs/busybody/protocol"

type Handler interface {
	HandleMessage(msg *protocol.Message) error
}

type HandlerFunc func(msg *protocol.Message) error

func (h HandlerFunc) HandleMessage(msg *protocol.Message) error {
	return h(msg)
}
