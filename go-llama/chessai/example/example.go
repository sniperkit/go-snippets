package main

import (
	"fmt"
	"github.com/gophergala/go-llama/chessai"
	"github.com/gophergala/go-llama/chessverifier"
	"math/rand"
)

func main() {
	fmt.Println("Test AI running")

	PropUsername := "testAI"
	PropPassword := "testAI"
	VersesAi := true
	FirstUse := false
	chessai.Make(PropUsername, PropPassword, VersesAi, FirstUse, RandSolver, IncomingChat)

	addr := "ws://chess.maycontainawesome.com:80/ws"
	host := "http://localhost"

	chessai.Run(addr, host) //this is a blocking call
}

func RandSolver(game chessverifier.GameState) []byte {
	//best solver ever
	availableMoves := chessverifier.GetAllValidMoves(&game, chessai.IsWhite())
	return availableMoves[rand.Intn(len(availableMoves)-1)]
}

func IncomingChat(messageId int) {
	//echo the same message back
	chessai.SendChat(messageId)
}
