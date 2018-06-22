//Package chessai is an application framework for creating AIs to connect to the go-llama intchess internet chess server.
//It allows for users to define their AI in terms of processing moves.
//Here is a trivial (and mostly non-functional) example:
// package main

// import (
// 	"fmt"
// 	"github.com/gophergala/go-llama/chessai"
// 	"github.com/gophergala/go-llama/chessverifier"
// )

// func main() {
// 	fmt.Println("Test AI running")

// 	PropUsername := "testAI"
// 	PropPassword := "testAI"
// 	VersesAi := true
// 	FirstUse := false
// 	chessai.Make(PropUsername, PropPassword, VersesAi, FirstUse, MySolver, IncomingChat)

// 	addr := "ws://192.168.1.25:8080/ws"
// 	host := "http://localhost"

// 	chessai.Run(addr, host) //this is a blocking call
// }

// func MySolver(game chessverifier.GameState) []byte {
// 	//best solver ever
// 	return []byte("a2-a3")
// 	//valid moves - who cares about those PFFT
// }

// func IncomingChat(messageId int) {
// 	//echo the same message back
// 	chessai.SendChat(messageId)
// }

//As can be seen, developers need to pass in both a Solver and a IncomingChat function to the chessai.Make function.
//Both the chessai and the chessverifier packages will be needed.

package chessai

import (
	"encoding/json"
	"fmt"
	"github.com/gophergala/go-llama/chessverifier"
	"github.com/gophergala/go-llama/intchess"
	"golang.org/x/net/websocket"
	"log"
)

type AI struct {
	PropUsername     string
	PropPassword     string
	VersesAi         bool
	FirstUse         bool
	ws               *websocket.Conn
	User             *intchess.User
	sendMessages     chan string
	receivedMessages chan string
	Solve            func(chessverifier.GameState) []byte
	Chat             func(int)
	CurrentlyWhite   bool
}

var ai AI

//Prepares the AI engine for operation.
//Must pass in a proposed username and a proposed password
//VersesAi is only checked if FirstUse is true, which marks if the AI has been used before (ie whether or not it should register on the server)
//Solve is the developer's AI function which determines what move to make given a board state
//Chat is the developer's AI function which listens to incoming chat requests.
func Make(PropUsername string, PropPassword string, VersesAi bool, FirstUse bool, Solve func(chessverifier.GameState) []byte, Chat func(int)) {
	ai.PropPassword = PropPassword
	ai.PropUsername = PropUsername
	ai.FirstUse = FirstUse
	ai.Solve = Solve
	ai.sendMessages = make(chan string, 256)
	ai.receivedMessages = make(chan string, 256)
	ai.Chat = Chat
}

//Runs the AI engine.
//addr is the remote address of the server (should be something like ws://example.com/)
//host is the local address (usually http://localhost/)
func Run(addr string, host string) {

	ws, err := websocket.Dial(addr, "", host)

	if err != nil {
		log.Fatal(err)
	}

	ai.ws = ws

	defer close(ai.sendMessages)
	defer close(ai.receivedMessages)

	go ai.Reader()
	go ai.Writer()
	ai.Runner()
}

//SendChat is a thread-safe function which can be called at any time to send a chat message
//by the developer.
func SendChat(messageId int) {
	request := intchess.APIGameChatRequest{
		Type:      "game_chat_request",
		MessageId: messageId,
	}
	msg, _ := json.Marshal(request)
	ai.sendMessages <- string(msg)
}

//The AI framework websocket reader
func (a *AI) Reader() {
	for {
		var reply string
		err := websocket.Message.Receive(a.ws, &reply)
		if err != nil {
			break
		}

		a.receivedMessages <- reply
	}
	a.ws.Close()
}

//The AI framework websocket writer
func (a *AI) Writer() {
Loop:
	for {
		for msg := range a.sendMessages {
			err := websocket.Message.Send(a.ws, msg)
			if err != nil {
				break Loop
			}
		}
		//how to detect if broken if no messages to send?
		if !a.ws.IsClientConn() {
			break Loop
		}
	}
	a.ws.Close()
}

//The main operation loop of the AI framework
func (a *AI) Runner() {
	if a.FirstUse {
		a.attemptCreateAndAuthenticate(a.PropUsername, a.PropPassword, a.VersesAi)
	} else {
		a.attemptAuthentication(a.PropUsername, a.PropPassword)
	}
Loop:
	for {
		for msg := range a.receivedMessages {
			a.DecodeMessage([]byte(msg))
		}
		if !a.ws.IsClientConn() {
			break Loop
		}
	}

}

//Send an authentication request to the server
func (a *AI) attemptAuthentication(username string, proposed_password string) {
	request := intchess.APIAuthenticationRequest{
		Type:      "authentication_request",
		Username:  username,
		UserToken: proposed_password,
	}
	msg, _ := json.Marshal(request)
	ai.sendMessages <- string(msg)
	return
}

//Send a create and authenticate to the server
func (a *AI) attemptCreateAndAuthenticate(username string, proposed_password string, verses_ai bool) {
	r := intchess.APISignupRequest{
		Type:      "signup_request",
		Username:  username,
		UserToken: proposed_password,
		IsAi:      true,
		VersesAi:  verses_ai,
	}
	msg, _ := json.Marshal(r)
	ai.sendMessages <- string(msg)
	return
}

//Decode a message from the server and call appropriate sub function (if neccessary)
func (a *AI) DecodeMessage(message []byte) {
	var t intchess.APITypeOnly
	if err := json.Unmarshal(message, &t); err != nil {
		fmt.Println("Just receieved a message I couldn't decode:")
		fmt.Println(string(message))
		fmt.Println("Exact error: " + err.Error())
		fmt.Println()
		return
	}

	switch t.Type {
	case "authentication_response":
		var r intchess.APIAuthenticationResponse
		if err := json.Unmarshal(message, &r); err != nil {
			return
		}
		if r.Response == "ok" && r.User != nil {
			a.User = r.User
		} else {
			log.Fatalf("Could not sign in.")
		}
	case "game_request":

		//accept all game requests
		a.SendGameAccept()
	case "game_response_rejection":
		//bleh, ignore it
		break
	case "game_response_accepted":
		//bleh, ignore it
		break
	case "game_move_update":
		var r intchess.APIGameOutput
		if err := json.Unmarshal(message, &r); err != nil {
			return
		}
		a.processMove(r.Game)
	case "game_over":
		//do we need to do anything?
		break
	case "game_chat":
		var r intchess.APIGameChat
		if err := json.Unmarshal(message, &r); err != nil {
			return
		}
		if r.From.Username != a.User.Username {
			a.Chat(r.MessageId)
		}
	default:
		log.Printf("I'm not familiar with type %s\n", t.Type)
	}
}

//Send a game accept message to the user
func (a *AI) SendGameAccept() {
	r := intchess.APIGameResponse{
		Type:     "game_response",
		Response: "ok",
	}
	msg, _ := json.Marshal(r)
	ai.sendMessages <- string(msg)
}

func IsWhite() bool {
	return ai.CurrentlyWhite
}

//Process move function, redirect request to developer's solver function and form request from their move
func (a *AI) processMove(g *intchess.ChessGame) {
	if g.WhitePlayer.Username == a.User.Username {
		a.CurrentlyWhite = true
	} else {
		a.CurrentlyWhite = false
	}
	if len(g.GameMoves)%2 == 0 && g.WhitePlayer.Username == a.User.Username || len(g.GameMoves)%2 == 1 && g.BlackPlayer.Username == a.User.Username {
		//is yer turn

		//now verify that the move is allowable
		//first, we need to convert all our moves from []GameMove to [][]Byte
		bMoves := make([][]byte, len(g.GameMoves))
		for index, move := range g.GameMoves {
			bMoves[index] = []byte(move.Move)
		}

		//get the current board state
		curGameState := chessverifier.GetBoardState(&bMoves)

		//get what to do
		r := intchess.APIGameMoveRequest{
			Type: "game_move_request",
			Move: string(a.Solve(curGameState)),
		}

		msg, _ := json.Marshal(r)
		ai.sendMessages <- string(msg)
	}
}
