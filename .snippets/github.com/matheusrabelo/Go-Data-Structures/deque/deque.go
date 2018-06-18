package deque

type Node struct {
	value    interface{}
	previous *Node
	next     *Node
}

type Deck struct {
	first *Node
	last  *Node
	size  int
}

//PushFirst pushes given element to the front of the deck and returns the created Node
func (d *Deck) PushFirst(value interface{}) *Node {
	n := new(Node)
	n.value = value
	n.previous = nil
	n.next = d.first
	d.first = n
	d.size = d.size + 1
	return n
}

//PushLast pushes given element to the end of the deck and returns the created Node
func (d *Deck) PushLast(value interface{}) *Node {
	n := new(Node)
	n.value = value
	n.previous = d.last
	n.next = nil
	d.last = n
	d.size = d.size + 1
	return n
}

//PopFirst removes the front element of the deck and return his value
func (d *Deck) PopFirst() interface{} {
	if d.size == 0 {
		return nil
	}
	value := d.first.value
	d.first = d.first.next
	d.size = d.size - 1
	return value
}

//PopLast removes the last element of the deck and return his value
func (d *Deck) PopLast() interface{} {
	if d.size == 0 {
		return nil
	}
	value := d.last.value
	d.last = d.last.previous
	d.size = d.size - 1
	return value
}

//First returns the value of the first element of the deck
func (d *Deck) First() interface{} {
	if d.size == 0 {
		return nil
	}
	return d.first.value
}

//Last returns the value of the last element of the deck
func (d *Deck) Last() interface{} {
	if d.size == 0 {
		return nil
	}
	return d.last.value
}

//Size returns the size of the Deck
func (d *Deck) Size() int {
	return d.size
}

//Empty returns true if the deck is empty
func (d *Deck) Empty() bool {
	return d.size == 0
}

//New creates and returns a new deck instance
func New() *Deck {
	q := new(Deck)
	q.first = nil
	q.last = nil
	q.size = 0
	return q
}
