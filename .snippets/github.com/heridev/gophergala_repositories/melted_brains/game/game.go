package game

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

const (
	Created Status = iota
	Started
	Ended
)

type MessageChannel chan string
type KillChannel chan bool

var SampleFiles []string

func init() {
	var err error
	SampleFiles, err = filepath.Glob("static/samples/*")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%v\n", SampleFiles)
}

type Game struct {
	Id string
	Status
	Clients
	MessageChannel
	KillChannel
	CodeFile string
}

func (g *Game) Code() string {
	content, _ := ioutil.ReadFile(g.CodeFile)
	return string(content)
}
func (g *Game) KillBroadCast() {
	g.KillChannel <- true
}
func (g *Game) PublishClients() {
	g.Publish(fmt.Sprintf("users:%s", g.Clients.Serialize()))
}

func (g *Game) Add(username string, conn *websocket.Conn) (int, error) {
	if g.Status != Created {
		return 0, errors.New("Already started")
	}
	fmt.Printf("g.Clients : %v %v\n", g.Clients, username)
	g.Clients = append(g.Clients, NewClient(username, conn))
	g.PublishClients()
	id := len(g.Clients) - 1
	websocket.Message.Send(conn, fmt.Sprintf("current_user:%d", id))
	if id == 3 {
		g.Start()

	}
	//TODO: Check number of clients
	//TODO: Start game if full
	return id, nil
}
func (g *Game) RemoveClients(toRemove Clients) {
	newClients := Clients{}
	for _, client := range g.Clients {
		if !toRemove.Contains(client) {
			newClients = append(newClients, client)
		}
	}
	g.Clients = newClients
}

func (g *Game) Publish(event string) {
	g.MessageChannel <- event
}

func (g *Game) Start() {
	g.Status = Started
	g.PublishStart()
}

func (g *Game) PublishStart() {
	g.Publish("start:")
}

func (g *Game) PublishFromUser(userId int, char string) {
	g.Publish(fmt.Sprintf("k:%d#%s", userId, char))
}

func (g *Game) BroadCast() {
	for {
		select {
		case msg := <-g.MessageChannel:
			erroredClients := Clients{}
			for _, client := range g.Clients {
				if err := websocket.Message.Send(client.Conn, msg); err != nil {
					erroredClients = append(erroredClients, client)
				}
			}
			if len(erroredClients) > 0 {
				g.RemoveClients(erroredClients)
			}

		case _ = <-time.After(10 * time.Minute):
			Repository.Delete(g)
			log.Printf("Timing out game %s\n", g.Id)
			return
		}
	}
}

type Status int

func NewGame() *Game {
	game := &Game{Id: NewId(), Clients: Clients{}, MessageChannel: make(MessageChannel, 100), CodeFile: SampleFiles[rand.Intn(len(SampleFiles))]}
	go game.BroadCast()
	return game
}

func NewId() string {
	hash := sha1.Sum([]byte(time.Now().String()))
	return strings.Replace(base64.URLEncoding.EncodeToString(hash[:]), "=", "", -1)
}
