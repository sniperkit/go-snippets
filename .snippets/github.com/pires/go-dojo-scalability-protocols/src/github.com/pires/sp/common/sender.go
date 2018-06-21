package common

//
type Sender interface {
	Send(message *Message, address string) error
}
