package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type User struct {
	viewX, viewY float32
	health       float32
	spawned      bool

	inputMutex  sync.Mutex
	keys        map[string]bool
	chatMessage string
	Username    string
	pong        bool

	connState sync.RWMutex
	connected bool
	messages  chan *UserMessage
}

func wsHandler(s *websocket.Conn) {
	log.Println("User connected from", s.RemoteAddr())
	defer log.Println("User disconnected from", s.RemoteAddr())
	var u User
	u.keys = make(map[string]bool)
	u.chatMessage = ""
	u.Username = ""
	u.pong = true

	u.connected = true
	u.messages = make(chan *UserMessage)

	go u.handleInput(s)
	go u.handleOutpit(s)
	defer u.disconnect()

	NewUser <- &u

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for range ticker.C {
		u.inputMutex.Lock()
		pong := u.pong
		u.pong = false
		u.inputMutex.Unlock()

		if !pong || !u.Connected() {
			s.Close()
			return
		}
		m := UserMessage{"ping", nil}
		u.Send(&m)
	}
}

func (u *User) handleOutpit(w io.Writer) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

	var err error
	for m := range u.messages {
		err = encoder.Encode(m)
		if err != nil {
			break
		}

		if err != nil {
			break
		}
		_, err = w.Write(buf.Bytes())
		if err != nil {
			break
		}
		buf.Reset()
	}
	if err != nil {
		log.Println("User error sending message,", err)
	}
	u.disconnect()
	for range u.messages {
	}
}

type UserMessage struct {
	Event string
	Data  interface{}
}

func (u *User) Send(m *UserMessage) {
	u.connState.RLock()
	if u.connected {
		u.messages <- m
	}
	u.connState.RUnlock()
}

func (u *User) disconnect() {
	u.connState.Lock()
	if u.connected {
		u.connected = false
		close(u.messages)
	}
	u.connState.Unlock()
}

func (u *User) Connected() bool {
	u.connState.RLock()
	defer u.connState.RUnlock()
	return u.connected
}

func (u *User) handleInput(r io.Reader) {
	defer u.disconnect()
	decoder := json.NewDecoder(r)

	for {
		m := make(map[string]interface{})
		err := decoder.Decode(&m)
		if err != nil {
			log.Println("Error reading message from client,", err)
			return
		}

		event, ok := m["Event"].(string)
		if !ok {
			log.Println("Unable to cast event from user to string")
			return
		}

		u.inputMutex.Lock()
		switch event {
		case "chatMessage":
			if msg, ok := m["Message"].(string); ok {
				u.chatMessage = msg
			}
		case "username":
			if name, ok := m["User"].(string); ok {
				u.Username = name
			}
		case "pong":
			u.pong = true
		case "team":
			if tstr, ok := m["Team"].(string); ok {
				var t team
				if tstr == "gophers" {
					t = TeamGophers
				} else if tstr == "pythons" {
					t = TeamPythons
				}
				if !u.spawned && t != TeamNone && u.Username != "" {
					go NewPlayerShip(u, t)
					u.spawned = true
				}
			}
		default:
			if event[1:] == " down" {
				u.keys[event[:1]] = true
			} else if event[1:] == " up" {
				u.keys[event[:1]] = false
			} else {
				log.Println("Unkown user message:", m)
			}
		}
		u.inputMutex.Unlock()
	}
}

func (u *User) Key(k string) bool {
	u.inputMutex.Lock()
	defer u.inputMutex.Unlock()
	return u.keys[k]
}

func (u *User) GetChatMessage() *UserMessage {
	u.inputMutex.Lock()
	defer u.inputMutex.Unlock()
	if u.chatMessage != "" {
		log.Println(u.Username, ":", u.chatMessage)
		type Chatmsg struct {
			User    string
			Message string
		}
		msg := UserMessage{
			"chatMessage",
			Chatmsg{
				u.Username,
				u.chatMessage,
			},
		}
		u.chatMessage = ""
		return &msg
	}
	return nil
}

func (u *User) View(x, y float32) {
	u.viewX = x
	u.viewY = y
}

func (u *User) render(overworld *Overworld, planetInfos []PlanetInfo,
	shipInfos []shipInfo, wait chan *User) {
	type ScreenUpdate struct {
		ViewX            float32
		ViewY            float32
		Objs             []RenderInfo
		Planets          []PlanetInfo
		Ships            []shipInfo
		PlanetAllegance  string
		AllegancePercent float32
		Health           float32
		EngineSound bool
	}

	var s ScreenUpdate
	s.ViewX = u.viewX
	s.ViewY = u.viewY
	s.Health = u.health
	entities := overworld.query(nil, u.viewX, u.viewY, 1000)
	s.Planets = planetInfos
	s.Objs = make([]RenderInfo, len(entities))
	s.Ships = shipInfos
	s.EngineSound = false

	for i, entity := range entities {
		s.Objs[i] = entity.RenderInfo()
		if planet, ok := entity.(*Planet); ok {
			s.AllegancePercent, s.PlanetAllegance = planet.Allegance()
		}
		if sound, ok := entity.(*Sound); ok {
			m := UserMessage{"sound", sound.name}
			u.Send(&m)
		}
		if player, ok := entity.(*PlayerShip); ok {
			if player.vx != 0 || player.vy != 0{
				s.EngineSound = true
			}
		}
		if pirate, ok := entity.(*PirateShip); ok {
			if pirate.vx != 0 || pirate.vy != 0 {
				s.EngineSound = true
			}
		}
	}

	if u.Connected() {
		wait <- u
	} else {
		wait <- nil
	}

	m := UserMessage{
		"screenUpdate", s,
	}
	u.Send(&m)
}
