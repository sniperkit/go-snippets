package game

import (
	"fmt"
	"strings"

	"golang.org/x/net/websocket"
)

type Client struct {
	Username string
	*websocket.Conn
}

func NewClient(username string, conn *websocket.Conn) *Client {
	return &Client{username, conn}
}

type Clients []*Client

func (clients Clients) Contains(c *Client) bool {
	for _, b := range clients {
		if b == c {
			return true
		}
	}
	return false
}
func (clients Clients) Serialize() string {
	clientReps := []string{}
	for index, client := range clients {
		clientReps = append(clientReps, fmt.Sprintf("%d#%s", index, client.Username))
	}
	return strings.Join(clientReps, "&")
}
