package chessverifier

import (
	"testing"
)

var game GameState

func init() {

}

func Test_moveEqual(t *testing.T) {
	a := []byte{'a', '1'}
	b := []byte{'a', '1'}
	if moveEqual(&a, &b) != true {
		t.Errorf("A did not equal B when it should have")
	}
	a = []byte{'a', '1'}
	b = []byte{'a', '2'}
	if moveEqual(&a, &b) == true {
		t.Errorf("A equalled B when it should not have")
	}
	a = []byte{'a', '1'}
	b = []byte{'a'}
	if moveEqual(&a, &b) == true {
		t.Errorf("A equalled B when it should not have")
	}
	a = []byte{'a', '1'}
	b = []byte{'b', '1'}
	if moveEqual(&a, &b) == true {
		t.Errorf("A equalled B when it should not have")
	}
}

func Test_getSquareIndices(t *testing.T) {
	a := []byte{'a', '1'}
	x, y := getSquareIndices(a)
	if x != 0 && y != 0 {
		t.Errorf("a1 did not map correctly to 0,0")
	}
	a = []byte{'c', '3'}
	x, y = getSquareIndices(a)
	if x != 2 && y != 2 {
		t.Errorf("a1 did not map correctly to 2,2")
	}
	a = []byte{'h', '8'}
	x, y = getSquareIndices(a)
	if x != 7 && y != 7 {
		t.Errorf("h8 did not map correctly to 7,7")
	}
}

func Test_NewGame(t *testing.T) {
	game = NewGame()
	if game.board[0][0][0] != 'W' {
		t.Errorf("0,0,0 in boardState not W")
	}
}

func Test_MakeMove(t *testing.T) {
	move := []byte{'a', '2', '-', 'a', '4'} //moves the left-most white pawn forward two spaces
	MakeMove(&game, &move)
	oldX, oldY := getSquareIndices(move[0:2])
	newX, newY := getSquareIndices(move[3:5])
	if len(game.board[oldX][oldY]) != 0 {
		t.Errorf("Old location not cleared")
	}
	if len(game.board[newX][newY]) != 3 {
		t.Errorf("New location not filled")
	}
	if game.board[newX][newY][0] != 'W' && game.board[newX][newY][1] != 'P' && game.board[newX][newY][0] != '1' {
		t.Errorf("Doesn't appear to have put the right piece to a4 (%v, %v)", newX, newY)
	}
}

func Test_occupied(t *testing.T) {
	game = NewGame()

	if occupied(&game, [2]byte{4, 4}) {
		t.Errorf("Square was occupied when it should not have been")
	}
	if !occupied(&game, [2]byte{1, 1}) {
		t.Errorf("Square was not occupied when it should have been")
	}
}

func Test_GetValidMoves(t *testing.T) {

}

func Test_GetAllValidMoves(t *testing.T) {

}

func Test_IsMoveValid(t *testing.T) {

}

func Test_GetBoardState(t *testing.T) {

}
