package intchess

import (
	"errors"
	"fmt"
	"github.com/gophergala/go-llama/chessverifier"
	"time"
)

type ChessGame struct {
	Id            int64                    `json:"game_id"`
	CreatedAt     NullTime                 `json:"commenced_at"`
	FinishedAt    NullTime                 `json:"finished_at"`
	Status        string                   `json:"game_status"`
	GameMoves     []GameMove               `json:"game_moves"`
	WhitePlayer   *User                    `sql:"-" json:"white"` //these should be used just for JSONing and not SQLing
	WhitePlayerId int64                    `json:"-"`             //and vice versa
	BlackPlayer   *User                    `sql:"-" json:"black"` //
	BlackPlayerId int64                    `json:"-"`             //
	Winner        *UserRankChange          `sql:"-" json:"winner",omitempty`
	Loser         *UserRankChange          `sql:"-" json:"loser","omitempty`
	BoardStatus   chessverifier.BoardState `sql:"-" json:"board_status"`
	WhiteAccept   bool                     `json:"-" sql:"-"`
	BlackAccept   bool                     `json:"-" sql:"-"`
	WhiteConn     *Connection              `json:"-" sql:"-"`
	BlackConn     *Connection              `json:"-" sql:"-"`
}

type UserRankChange struct {
	Id          int64    `json:"-"`
	UserId      int64    `json:"user_id"`
	ChessGameId int64    `json:"-"`
	RankChange  int      `json:"rank_change"`
	CreatedAt   NullTime `json:"-"`
}

type GameMove struct {
	Id          int64    `json:"id"`
	ChessGameId int64    `json:"-"`
	Move        string   `json:"move"`
	CreatedAt   NullTime `json:"created_at"`
}

//Function to return a new game object
func NewGame(white *Connection, black *Connection) ChessGame {
	g := ChessGame{
		Status:        "new",
		BoardStatus:   chessverifier.StartBoardState,
		WhitePlayer:   white.User,
		WhiteConn:     white,
		WhitePlayerId: white.User.Id,
		BlackPlayer:   black.User,
		BlackConn:     black,
		BlackPlayerId: black.User.Id,
		CreatedAt:     NullTime{Time: time.Now(), Valid: true},
	}
	return g
}

//Function to mark a player has accepted the game
//if both players have now accepted, this returns true and creates the game formally
//by sending it to the database
func (c *ChessGame) PlayerAccepts(player *User) (GameStarts bool) {
	GameStarts = false
	if player == c.WhitePlayer {
		c.WhiteAccept = true
	} else if player == c.BlackPlayer {
		c.BlackAccept = true
	}
	if c.WhiteAccept && c.BlackAccept {
		GameStarts = true
		c.BeginGame()
	}
	return
}

func (c *ChessGame) BeginGame() {
	fmt.Println("A game has been commenced.")
	c.Status = "in_progress"
	c.Create()
	c.SendMoveUpdate()
}

//Function to see if a game has expired (ie has not begun 15 seconds after prompting users if they are ready to play)
//If true, it notifies the players that the game expired
func (c *ChessGame) Expired() bool {
	if c.Status == "new" && c.CreatedAt.Time.Before(time.Now().Add(time.Second*-15)) {
		if c.WhiteConn != nil && GS.connections[c.WhiteConn] {
			c.WhiteConn.SendGameExpired()
		}
		if c.BlackConn != nil && GS.connections[c.BlackConn] {
			c.BlackConn.SendGameExpired()
		}
		return true
	}
	return false
}

//Function to create a chess game into the database
func (c *ChessGame) Create() {
	//this weird adding and deleting variables is needed for GORM. Yay.
	dbGorm.Create(c)
}

func (c *ChessGame) LoadAllFromId(id int64) {

}

func (c *ChessGame) SendMoveUpdate() {
	c.WhiteConn.SendGameUpdate(c, "game_move_update")
	c.BlackConn.SendGameUpdate(c, "game_move_update")
}

//This checks to see if the game is over, and if it is over, ends the game
func (c *ChessGame) End() bool {
	//check disconnections
	if c.Status == "in_progress" {
		if !c.ClientsStillConnected() {
			c.Status = "disconnected"
			c.EndGame()
			return true
		}
	} else if c.Status == "victory" || c.Status == "stalemate" {
		return true
	}
	//check checkmate, etc
	return false
}

//This forcibly ends a game (without setting the status)
func (c *ChessGame) EndGame() {
	fmt.Println("Game between " + c.WhitePlayer.Username + " and " + c.BlackPlayer.Username + " ends with status " + c.Status)
	if c.WhiteConn != nil && GS.connections[c.WhiteConn] {
		c.WhiteConn.SendGameUpdate(c, "game_over")
	}
	if c.BlackConn != nil && GS.connections[c.BlackConn] {
		c.BlackConn.SendGameUpdate(c, "game_over")
	}
	c.FinishedAt = NullTime{Time: time.Now(), Valid: true}
	if err := dbGorm.Save(c).Error; err != nil {
		fmt.Println("DBError at chessgame:141: " + err.Error())
	}

}

func (c *ChessGame) ClientsStillConnected() bool {
	ok := true
	if c.WhiteConn == nil {
		return false
	}
	if c.BlackConn == nil {
		return false
	}
	wVal, wOk := GS.connections[c.WhiteConn]
	bVal, bOk := GS.connections[c.BlackConn]
	if !(wVal && wOk && bVal && bOk) {
		ok = false
	}
	return ok
}

func (c *ChessGame) AttemptMove(moveStr string, moverConn *Connection) error {
	//first, make sure that it's the right player moving
	if len(c.GameMoves)%2 == 0 && moverConn != c.WhiteConn || len(c.GameMoves)%2 == 1 && moverConn != c.BlackConn {
		//wrong player attempted to move!
		c.SendMoveUpdate() //allows for redrawing
		return errors.New("not your turn")
	}
	//now verify that the move is allowable
	//first, we need to convert all our moves from []GameMove to [][]Byte
	bMoves := make([][]byte, len(c.GameMoves))
	for index, move := range c.GameMoves {
		bMoves[index] = []byte(move.Move)
	}

	//get the current board state
	curGameState := chessverifier.GetBoardState(&bMoves)

	curMove := []byte(moveStr)

	//check if this move is valid'
	valid := chessverifier.IsMoveValid(&curGameState, &curMove)
	if !valid {
		c.SendMoveUpdate() //allows for redrawing
		return errors.New("Invalid move attempted")
	}

	chessverifier.MakeMove(&curGameState, &curMove)
	c.GameMoves = append(c.GameMoves, GameMove{Move: moveStr})
	c.BoardStatus = curGameState.Board
	GameOver := false

	if len(c.GameMoves)%2 == 0 {
		//that was the end of BLACK's turn
		mate, check := chessverifier.IsMate(&curGameState, true)
		if mate == true && check == true {
			//black just won
			c.GameFinished(true, false)
			GameOver = true
		} else if mate == true && check == false {
			//stalemate
			c.GameFinished(false, false)
			GameOver = true
		}
	} else {
		//that was the end of WHITE's turn
		mate, check := chessverifier.IsMate(&curGameState, false)
		if mate == true && check == true {
			//white just won
			c.GameFinished(true, true)
			GameOver = true
		} else if mate == true && check == false {
			//stalemate
			c.GameFinished(false, false)
			GameOver = true
		}

	}
	if GameOver == false {
		c.SendMoveUpdate()
	}
	return nil
}

func (c *ChessGame) GameFinished(notStalemate bool, whiteWon bool) {
	if notStalemate == false {
		c.Status = "stalemate"
	}
	if whiteWon == true {
		c.Status = "victory"
		c.Winner, c.Loser = c.WhiteConn.User.WonGame(c, c.BlackConn.User)
	} else {
		c.Status = "victory"
		c.Winner, c.Loser = c.BlackConn.User.WonGame(c, c.WhiteConn.User)
	}
	c.WhitePlayer = c.WhiteConn.User //so JSON version can see new rank
	c.BlackPlayer = c.BlackConn.User

	c.EndGame()
}

func (c *ChessGame) Chat(messageId int, sender *Connection) {
	if c.WhiteConn != nil && c.BlackConn != nil && sender != nil && sender.User != nil {
		c.WhiteConn.SendChat(messageId, sender.User)
		c.BlackConn.SendChat(messageId, sender.User)
	}
}

func (c *ChessGame) GetValidMoves(loc string) [][]byte {
	//now verify that the move is allowable
	//first, we need to convert all our moves from []GameMove to [][]Byte
	bMoves := make([][]byte, len(c.GameMoves))
	for index, move := range c.GameMoves {
		bMoves[index] = []byte(move.Move)
	}

	//get the current board state
	curGameState := chessverifier.GetBoardState(&bMoves)
	x, y := chessverifier.GetSquareIndices([]byte(loc))
	if x == -1 {
		return nil
	}
	return chessverifier.GetValidMoves(&curGameState, x, y)
}
