package channels

// A Channel has an ID and one or more subscribers
type Channel struct {
	id          uint32
	subscribers map[uint32]struct{}
}
