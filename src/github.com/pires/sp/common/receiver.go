package common

//
type Receiver interface {
	Receive(address string) error
}
