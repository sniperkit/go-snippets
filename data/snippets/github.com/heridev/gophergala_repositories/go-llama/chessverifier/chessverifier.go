package chessverifier

import (
// "encoding/json"
// "// fmt"
)

type BoardState [8][8][]byte
type GameState struct {
	Board    BoardState
	MoveList [][]byte
}

func NewGame() GameState {
	var game GameState
	game.Board = StartBoardState
	return game
}

// func main() {
// 	// fmt.Printf("Hello world!\n")

// }

func GetValidMoves(game *GameState, x, y int) [][]byte {
	var validMoves [][]byte
	// fmt.Println("GetValidMoves", x, y)
	var piece = game.Board[x][y]
	// fmt.Println(string(piece))
	if len(piece) == 0 {
		return [][]byte{}
	}
	var newMove = []byte{}
	var white = piece[0] == 'W'
	switch piece[1] { //make sure the king is not in check before anything else
	case 'P': //@todo add on passant
		var ymove = 1
		if !white {
			ymove = -1
		}
		var newSquare = [2]int{x, y + ymove}
		var free, taking = canLand(game, newSquare, white) //moving forward
		if free && !taking {
			newMove = getMove([2]int{x, y}, newSquare)
			validMoves = append(validMoves, newMove)
		}
		newSquare[0] = newSquare[0] + 1 //taking to the right
		free, taking = canLand(game, newSquare, white)
		if free && taking {
			newMove = getMove([2]int{x, y}, newSquare)
			validMoves = append(validMoves, newMove)
		}
		newSquare[0] = newSquare[0] - 2 //taking to the left
		free, taking = canLand(game, newSquare, white)
		if free && taking {
			newMove = getMove([2]int{x, y}, newSquare)
			validMoves = append(validMoves, newMove)
		}

		if white && y == 1 {
			var free, taking = canLand(game, [2]int{x, y + 2}, white)
			if free && !taking {
				validMoves = append(validMoves, getMove([2]int{x, y}, [2]int{x, y + 2}))
			}
		} else if !white && y == 6 {
			var free, taking = canLand(game, [2]int{x, y - 2}, white)
			if free && !taking {
				validMoves = append(validMoves, getMove([2]int{x, y}, [2]int{x, y - 2}))
			}
		}

		//checking for on passant, right then left
		var left, right = canOnpassant(game, x, y, white)
		if right {
			validMoves = append(validMoves, getMove([2]int{x, y}, [2]int{x + 1, y + ymove}))
		}
		if left {
			validMoves = append(validMoves, getMove([2]int{x, y}, [2]int{x - 1, y + ymove}))
		}

	case 'R':
		for _, direction := range squareDirections {
			validMoves = append(validMoves, moveDirection(game, x, y, direction, white)...)
		}

	case 'N':
		var moveDiffs = knightMoves
		var free bool
		for i := range moveDiffs {
			free, _ = canLand(game, [2]int{x + moveDiffs[i][0], y + moveDiffs[i][1]}, white)
			if free {
				var newMove = getMove([2]int{x, y}, [2]int{x + moveDiffs[i][0], y + moveDiffs[i][1]})
				// fmt.Println("newMove=", string(newMove), i, x, y, moveDiffs[i])
				var testGame = testMove(*game, &newMove)
				if !IsCheck(&testGame, white) {
					validMoves = append(validMoves, newMove)
				}
			}
		}

	case 'B':
		for _, direction := range diagonalDirections {
			validMoves = append(validMoves, moveDirection(game, x, y, direction, white)...)
		}

	case 'Q':
		for _, direction := range append(squareDirections, diagonalDirections...) {
			validMoves = append(validMoves, moveDirection(game, x, y, direction, white)...)
		}

	case 'K':
		var moveDiffs [8][2]int = [8][2]int{[2]int{1, 0}, [2]int{1, 1}, [2]int{0, 1},
			[2]int{-1, 1}, [2]int{-1, 0}, [2]int{-1, -1}, [2]int{0, -1}, [2]int{1, -1}}

		for i := range moveDiffs {
			// fmt.Println("moveDiff", i, moveDiffs[i])
			var free, _ = canLand(game, [2]int{x + moveDiffs[i][0], y + moveDiffs[i][1]}, white)
			// fmt.Println(free)
			if free {
				var newMove = getMove([2]int{x, y}, [2]int{x + moveDiffs[i][0], y + moveDiffs[i][1]})
				var testGame = testMove(*game, &newMove)
				if !IsCheck(&testGame, white) {
					validMoves = append(validMoves, newMove)
				}
			}
		}

		//check for castling
		if !IsCheck(game, white) {
			// fmt.Println("castling")
			var king, left, right bool = true, true, true //left and right for each of the rooks
			var row byte
			var rownum int
			if white {
				row = '1'
				rownum = 0
			} else {
				row = '8'
				rownum = 7
			}
			for _, move := range game.MoveList {
				// fmt.Println(string(move))
				if move[1] == row {
					if move[0] == 'e' && move[1] == row {
						king = false
						// fmt.Println("king")
						break
					} else if move[0] == 'a' && move[1] == row && left {
						// fmt.Println("left")
						left = false
					} else if move[0] == 'h' && move[1] == row && right {
						// fmt.Println("right")
						right = false
					}
				}
			}
			// fmt.Println("can Castle")
			if king {
				if left && len(game.Board[1][rownum]) == 0 && len(game.Board[2][rownum]) == 0 && len(game.Board[3][rownum]) == 0 {
					// fmt.Println("row left clear")
					var testGame = testMove(*game, &[]byte{'e', row, '-', 'd', row})
					if !IsCheck(&testGame, white) {
						MakeMove(&testGame, &[]byte{'d', row, '-', 'c', row})
						if !IsCheck(&testGame, white) {
							validMoves = append(validMoves, []byte{'e', row, '-', 'c', row})
						}
					}
				}
				if right && len(game.Board[5][rownum]) == 0 && len(game.Board[6][rownum]) == 0 {
					// fmt.Println("row right clear")
					var testGame = testMove(*game, &[]byte{'e', row, '-', 'f', row})
					if !IsCheck(&testGame, white) {
						MakeMove(&testGame, &[]byte{'f', row, '-', 'g', row})
						if !IsCheck(&testGame, white) {
							validMoves = append(validMoves, []byte{'e', row, '-', 'g', row})
						}
					}
				}
			}
		}
	}
	// fmt.Println(validMoves)
	var freeMoves [][]byte
	var testGame GameState
	for _, move := range validMoves {
		testGame = testMove(*game, &move)
		if !IsCheck(&testGame, white) {
			freeMoves = append(freeMoves, move)
		}
	}
	return freeMoves
}

func moveDirection(game *GameState, x, y int, direction [2]int, white bool) (validMoves [][]byte) {
	for i := 1; i < 7; i++ {
		var free, taking = canLand(game, [2]int{x + (direction[0] * i), y + (direction[1] * i)}, white)
		// // fmt.Println(free)
		if free {
			validMoves = append(validMoves, getMove([2]int{x, y}, [2]int{x + (direction[0] * i), y + (direction[1] * i)}))
			if taking {
				break
			}
		} else {
			break
		}
	}
	// fmt.Println(validMoves)
	return validMoves
}

func IsCheck(game *GameState, white bool) bool {
	// fmt.Println("IsCheck")
	var pawnDirection = 1 //this is for finding attacking pawns
	if !white {
		pawnDirection = -1
	}
	var x, y int
	var pieceType byte

	for i, row := range game.Board { //find the king
		for j, piece := range row {
			if len(piece) == 3 && piece[1] == 'K' && ((piece[0] == 'W') == white) {
				x = i
				y = j
			}
		}
	}
	// // fmt.Println("Kinglocation=", x, y)
	// printBoard(game)
	if (onBoard([2]int{x + 1, y + pawnDirection}) && len(game.Board[x+1][y+pawnDirection]) != 0 && game.Board[x+1][y+pawnDirection][1] == 'P' && (game.Board[x+1][y+pawnDirection][0] == 'W') != white) ||
		(onBoard([2]int{x - 1, y + pawnDirection}) && len(game.Board[x-1][y+pawnDirection]) != 0 && game.Board[x-1][y+pawnDirection][1] == 'P' && (game.Board[x-1][y+pawnDirection][0] == 'W') != white) {
		return true
	}
	// fmt.Println("checking knights")
	var destination [2]int
	for _, moveDiff := range knightMoves {
		destination = [2]int{x + moveDiff[0], y + moveDiff[1]}
		// // fmt.Println("king destination=", destination)
		if onBoard(destination) && len(game.Board[destination[0]][destination[1]]) != 0 && game.Board[destination[0]][destination[1]][1] == 'N' && (game.Board[destination[0]][destination[1]][0] == 'W') != white {
			return true
		}
	}

	// fmt.Println("checking square")
	for _, direction := range squareDirections {
		for i := 1; i < 8; i++ {
			var canLand, taking = canLand(game, [2]int{x + (direction[0] * i), y + (direction[1] * i)}, white)
			if taking {
				// validMoves = append(validMoves, getMove([2]int{x, y}, [2]int{x + i, y}))
				pieceType = game.Board[x+(direction[0]*i)][y+(direction[1]*i)][1]
				if pieceType == 'Q' || pieceType == 'R' {
					return true
				}
				break
			} else if !canLand {
				break
			}
		}
	}

	// fmt.Println("Checking diagonal")
	for _, direction := range diagonalDirections {
		for i := 1; i < 8; i++ {
			var canLand, taking = canLand(game, [2]int{x + (direction[0] * i), y + (direction[1] * i)}, white)
			if taking {
				// validMoves = append(validMoves, getMove([2]int{x, y}, [2]int{x + i, y}))
				pieceType = game.Board[x+(direction[0]*i)][y+(direction[1]*i)][1]
				if pieceType == 'Q' || pieceType == 'B' {
					return true
				}
				break
			} else if !canLand {
				break
			}
		}
	}
	return false
}

func IsMate(game *GameState, white bool) (mate, check bool) {
	// fmt.Println("IsMate")
	if len(GetAllValidMoves(game, white)) == 0 {
		if IsCheck(game, white) {
			return true, true
		} else {
			return true, false
		}
	}
	return false, false
}

func canOnpassant(game *GameState, x, y int, white bool) (left, right bool) {
	if len(game.MoveList) <= 1 {
		return false, false
	}
	var ymove int
	if white {
		ymove = 1
		if y != 4 {
			return false, false
		}
	} else {
		ymove = -1
		if y != 3 {
			return false, false
		}
	}
	var free, taking = canLand(game, [2]int{x + 1, y + ymove}, white)
	if free && !taking {
		var onpassant = getMove([2]int{x + 1, y + (2 * ymove)}, [2]int{x + 1, y})
		if moveEqual(&game.MoveList[len(game.MoveList)-1], &onpassant) && game.Board[x+1][y][1] == 'P' {
			right = true
		}
	}
	free, taking = canLand(game, [2]int{x - 1, y + ymove}, white)
	if free && !taking {
		var onpassant = getMove([2]int{x - 1, y + (2 * ymove)}, [2]int{x - 1, y})
		if moveEqual(&game.MoveList[len(game.MoveList)-1], &onpassant) && game.Board[x-1][y][1] == 'P' {
			left = true
		}
	}
	return
}

//used for testing if a move will place the king in check so that functions can be passed only the board state after the move
func testMove(game GameState, move *[]byte) GameState {
	// fmt.Println("testMove")
	MakeMove(&game, move)
	return game
}

//white is whether or not the current player is white for handling taking
//if canLand is false taking sould be ignored (it will always be false)
func canLand(game *GameState, square [2]int, white bool) (canLand, taking bool) {
	// // fmt.Println("canLand", square, white)
	if !onBoard(square) { //is the destination off the board?
		return false, false
	}
	var piece = game.Board[square[0]][square[1]]
	var occupied bool = (len(piece) != 0)

	if !occupied { //on the board and empty
		return true, false
	}
	if (piece[0] == 'W') == white { //occupied with a piece the same colour as the piece being moved
		return false, false
	}
	return true, true //occupied with an oponent's piece
}

func onBoard(square [2]int) bool {
	// // fmt.Println("onboard", square, !(square[0] < 0 || square[0] > 7 || square[1] < 0 || square[1] > 7))
	return !(square[0] < 0 || square[0] > 7 || square[1] < 0 || square[1] > 7)
}

//do the line of formatting to save typing and tidy the code
func getMove(source [2]int, dest [2]int) []byte {
	return []byte{byte(source[0] + 'a'), byte(source[1] + '1'), '-', byte(dest[0] + 'a'), byte(dest[1] + '1')}
}

func GetAllValidMoves(game *GameState, white bool) (validMoves [][]byte) {
	// fmt.Println("GetAllValidMoves")
	var piece []byte
	for x := range game.Board {
		for y := range game.Board[x] {
			piece = game.Board[x][y]
			if len(piece) != 0 && (piece[0] == 'W') == white {
				validMoves = append(validMoves, GetValidMoves(game, x, y)...)
			}
		}
	}
	return
}

func IsMoveValid(game *GameState, move *[]byte) bool {
	// fmt.Println("IsMoveValid")
	var x, y = GetSquareIndices((*move)[0:2])
	var moveList = GetValidMoves(game, x, y)
	for _, testmove := range moveList {
		if moveEqual(move, &testmove) {
			return true
		}
	}
	return false
}

func GetBoardState(moveList *[][]byte) GameState {
	// fmt.Println("GetBoardState")
	var game GameState = NewGame()
	for moveNum := range *moveList {
		MakeMove(&game, &(*moveList)[moveNum])
	}
	return game
}

//This function is used to apply a move (eg "a2-a4") to the board
//note that this does not check if the move is valid because it causes and infinite loop XD
//you therefor need to make sure the move is allowed before trying it
func MakeMove(game *GameState, move *[]byte) { //@todo add ugrading pawns
	// fmt.Println("MakeMove")
	ox, oy := GetSquareIndices((*move)[0:2])
	nx, ny := GetSquareIndices((*move)[3:5])
	// fmt.Println(ox, oy, nx, ny)
	var piece = game.Board[ox][oy]
	// fmt.Println("piece=", string(piece))
	var white = piece[0] == 'W'
	var pieceType = piece[1]
	if pieceType == 'P' {
		var left, right = canOnpassant(game, ox, oy, white)
		if (left || right) && (nx != ox) {
			game.Board[nx][oy] = []byte{}
		} else if white && ny == 7 {
			// fmt.Println("upgrading", nx, ny)
			game.Board[ox][oy] = []byte{}
			var queencount = countQueens(game, white)
			game.Board[nx][ny] = []byte{'W', 'Q', byte(queencount + 1)}
			return
		} else if !white && ny == 0 {
			// fmt.Println("upgrading", nx, ny)
			game.Board[ox][oy] = []byte{}
			var queencount = countQueens(game, white)
			game.Board[nx][ny] = []byte{'B', 'Q', byte(queencount + 1)}
			return
		}
	} else if pieceType == 'K' {
		var dist = nx - ox
		if dist == 2 || dist == -2 {
			if nx == 6 {
				game.Board[5][ny] = game.Board[7][ny]
				game.Board[7][ny] = []byte{}
			} else if nx == 2 {
				game.Board[3][ny] = game.Board[0][ny]
				game.Board[0][ny] = []byte{}
			}
		}
	}
	game.Board[ox][oy] = []byte{}
	game.Board[nx][ny] = piece
	game.MoveList = append(game.MoveList, *move)
	// }
}

func countQueens(game *GameState, white bool) int {
	var count int
	for _, row := range game.Board { //find the king
		for _, piece := range row {
			if len(piece) == 3 && piece[1] == 'Q' && (piece[0] == 'W') == white {
				count++
			}
		}
	}
	// fmt.Println("newQueen =", count+1)
	return count
}

//This function is used to convert a piece location in Algebraic notation
//to a piece location in the internal board 2d slice
func GetSquareIndices(squareID []byte) (x, y int) {
	if len(squareID) != 2 {
		return -1, -1
	}
	// fmt.Println("square", string(squareID))
	y = int(squareID[1] - '1')
	x = int(squareID[0] - 'a')
	return
}

//This function is used to test the equality of two byte slices
func moveEqual(a, b *[]byte) bool {
	if len(*a) != len(*b) {
		return false
	}

	for i := range *a {
		if (*a)[i] != (*b)[i] {
			return false
		}
	}

	return true
}

//Knight is actaully spelt Night did you know?
var StartBoardState BoardState = BoardState{
	[8][]byte{[]byte{'W', 'R', '1'}, []byte{'W', 'P', '1'}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{'B', 'P', '1'}, []byte{'B', 'R', '1'}},
	[8][]byte{[]byte{'W', 'N', '1'}, []byte{'W', 'P', '2'}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{'B', 'P', '2'}, []byte{'B', 'N', '1'}},
	[8][]byte{[]byte{'W', 'B', '1'}, []byte{'W', 'P', '3'}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{'B', 'P', '3'}, []byte{'B', 'B', '1'}},
	[8][]byte{[]byte{'W', 'Q', '1'}, []byte{'W', 'P', '4'}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{'B', 'P', '4'}, []byte{'B', 'Q', '1'}},
	[8][]byte{[]byte{'W', 'K', '1'}, []byte{'W', 'P', '5'}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{'B', 'P', '5'}, []byte{'B', 'K', '1'}},
	[8][]byte{[]byte{'W', 'B', '3'}, []byte{'W', 'P', '6'}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{'B', 'P', '6'}, []byte{'B', 'B', '2'}},
	[8][]byte{[]byte{'W', 'N', '3'}, []byte{'W', 'P', '7'}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{'B', 'P', '7'}, []byte{'B', 'N', '2'}},
	[8][]byte{[]byte{'W', 'R', '2'}, []byte{'W', 'P', '8'}, []byte{}, []byte{}, []byte{}, []byte{}, []byte{'B', 'P', '8'}, []byte{'B', 'R', '2'}},
}

var diagonalDirections = [][2]int{[2]int{1, 1}, [2]int{1, -1}, [2]int{-1, -1}, [2]int{-1, 1}}
var squareDirections = [][2]int{[2]int{0, 1}, [2]int{0, -1}, [2]int{1, 0}, [2]int{-1, 0}}
var knightMoves = [8][2]int{[2]int{2, 1}, [2]int{2, -1}, [2]int{1, 2},
	[2]int{-1, 2}, [2]int{-2, 1}, [2]int{-2, -1}, [2]int{-1, -2}, [2]int{1, -2}}

// func Runtest() {
// 	var game = NewGame()
// 	game.Board[1][0] = []byte{}
// 	game.Board[2][0] = []byte{}
// 	// var board, _ = json.MarshalIndent(game.Board, "", "  ")
// 	// // fmt.Println(string(board))
// 	var moveList = GetValidMoves(&game, 0, 0)
// 	for _, move := range moveList {
// 		// fmt.Println(string(move))
// 	}
// }

// func printBoard(game *GameState) {
// 	var board, _ = json.MarshalIndent(game.Board, "", "  ")
// 	// fmt.Println(string(board))
// 	// // fmt.Printf("%s\n")
// 	// // fmt.Println(game.Board)
// }
