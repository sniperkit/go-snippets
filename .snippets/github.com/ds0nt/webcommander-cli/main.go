package main

import (
	"os"
	"log"
	"github.com/codegangsta/cli"
	"github.com/gorilla/websocket"
	"encoding/json"
)

const (
  socketBufferSize  = 1024
  messageBufferSize = 256
)

type socketClient struct {
  socket *websocket.Conn
  send   chan message
}

type message struct {
  Type string `json:"type"`
  Payload interface{} `json:"payload"`
}


type processor struct {
	In  chan *message
	Out chan *message
}

func newMessage(t string, payload interface{}) *message {
	return &message{t, payload}
}

func newProcessor(client *socketClient) *processor {
	in := make(chan *message, messageBufferSize)
	out := make(chan *message, messageBufferSize)
	return &processor{in, out}
}

func (p *processor) resultsTo(client *socketClient) {
	for msg := range p.Out {
    client.send <- *msg
	}
}

func newClient(socket *websocket.Conn) *socketClient {
  return &socketClient{socket, make(chan message, messageBufferSize)}
}

func (c *socketClient) read(processor *processor) {
  for {
    _, msg, err := c.socket.ReadMessage()
    if err == nil {
      var f message
      err := json.Unmarshal(msg, &f)
      if err != nil {
        log.Printf("Evil JSON Detected: %v, %v", err, string(msg))
        continue
      }
      log.Printf("Message Received: %v", f)
			processor.In <-&f
    } else {
      break
    }
  }
  c.socket.Close()
}

func (c *socketClient) write() {
  for msg := range c.send {
    bytes, err := json.Marshal(&msg)
    if err != nil {
      log.Printf("Client Write Json Marshal Error: %v, %v", err, msg)
    }
    if err := c.socket.WriteMessage(websocket.TextMessage, bytes); err != nil {
      log.Printf("Client Write Error: %v", err)
      break
    } else {
      log.Printf("Client Send: %v", msg)
    }
  }
}

func dial(url string) (*socketClient, *processor) {
	log.Printf("Dialing: %s", url)
	d := websocket.Dialer{}
	ws, _, err := d.Dial(url, nil)
	if err != nil {
		log.Panicf("An error occured while dialing: %s", err)
	}

  client := newClient(ws)
	processor := newProcessor(client)

	go processor.resultsTo(client)
  go client.read(processor)
  go client.write()
	return client, processor
}

func main() {

	log.Println("Web Commander CLI Client")
	// https://github.com/codegangsta/cli
	app := cli.NewApp()
	app.Name = "webcommander-cli"
	app.Usage = "Web Commander CLI Client"
	app.Commands = []cli.Command{
		{
			Name:    "respond",
			Action: func(c *cli.Context) {
				client, processor := dial(c.Args().First())
				bot := newRespondBot(processor, "ds0nt-bot", "Du hast mich?")
				go bot.run()
				client.read(processor)
			},
		},
		{
			Name:    "shorten",
			Action: func(c *cli.Context) {
				client, processor := dial(c.Args().First())
				bot, err := newShortenBot(processor, "ds0nt-bot")
				if err != nil {
					log.Fatalf("There was an error starting the shorten bot: %v", err)
				}
				go bot.run()
				client.read(processor)
			},
		},
		{
			Name:    "insight",
			Action: func(c *cli.Context) {
				client, processor := dial(c.Args().First())
				bot := newInsightBot(processor, "ds0nt-bot")
				go bot.run()
				client.read(processor)
			},
		},
	}
	app.Run(os.Args)
}
