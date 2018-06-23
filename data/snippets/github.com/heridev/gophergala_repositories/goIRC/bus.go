package main

import (
	"fmt"
	"net"
	"sync"
)

type EventType int
type Mode int

const (
	UserJoin EventType = iota
	UserPart
	PrivMsg
	Topic
)

const (
	Voice Mode = iota
	Moderator
	None
)

type Subscriber interface {
	OnEvent(*Event)
	GetInfo() string
}

type EventBus struct {
	subscribers map[EventType][]Subscriber
	channel     *Channel
	sync.Mutex
}

type Event struct {
	event_type EventType
	event_data string
}

type Channel struct {
	name  string
	topic string
	mode  map[string]Mode
}

func (u *User) GetInfo() string {
	return u.Nick
}

func (c *Channel) GetInfo() string {
	return c.name
}

func (u *User) Write(line string) {
	u.Conn.Write([]byte(line + "\r\n"))
}

func (u *User) WriteLines(lines []string) {
	for _, v := range lines {
		u.Write(v)
	}
}

func (u *User) OnEvent(event *Event) {
	switch event.event_type {
	case UserJoin:
		//fmt.Printf("%q(%d)> %q\n", s.Nick, event.event_type, event.event_data)
		_, err := u.Conn.Write([]byte(event.event_data))
		if err != nil {
			fmt.Println("Not looking too good")
		}
	case PrivMsg:
		_, err := u.Conn.Write([]byte(event.event_data))
		if err != nil {
			fmt.Println("Not looking too good")
		}
	case Topic:
		_, err := u.Conn.Write([]byte(event.event_data))
		if err != nil {
			fmt.Println("Not looking too good")
		}
	default:
		u.Conn.Write([]byte(event.event_data))
	}
}
func (bus *EventBus) Publish(event *Event) {
	fmt.Printf("\npublishing -%d- data: %q\n", event.event_type, event.event_data)
	for _, subscriber := range bus.subscribers[event.event_type] {
		go subscriber.OnEvent(event) //currently slower than without the goroutine
	}
	fmt.Println("done publishing")
}

func (bus *EventBus) Subscribe(event_type EventType, subscriber Subscriber) {
	bus.subscribers[event_type] = append(bus.subscribers[event_type], subscriber)
}

func (bus *EventBus) Unsubscribe(event_type EventType, subscriber Subscriber) {
	//find the index
	i := -1

	for index, val := range bus.subscribers[event_type] {
		if val.GetInfo() == subscriber.GetInfo() {
			i = index
			break
		}
	}

	if i > -1 { //we found someone
		cur := bus.subscribers[event_type]
		endIndex := i + 1 //will break if index is last element!
		cur = append(cur[0:i], cur[endIndex:]...)

		bus.Lock() //lock the eventbus while we remove the subscriber from the array
		bus.subscribers[event_type] = cur
		bus.Unlock()
	}

}

var buses map[string]*EventBus

func init() {
	loadMessages()
}
func main() {
	// init event bus map
	buses := make(map[string]*EventBus)

	ln, err := net.Listen("tcp", ":3030")

	if err != nil {
		panic("Listen not WORKING")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic("nope not Accepting")
		}
		go handleConnection(conn, buses)
	}

}
