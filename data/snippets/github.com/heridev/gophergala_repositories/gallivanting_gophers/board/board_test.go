package board

// func createFailure(b *Board, x int) {
// 	row := (x / 16) + 1
// 	column := x % 16
// 	b.SetWall(row-1, column-1)
// 	b.SetWall(row-1, column)
// 	b.SetWall(row-1, column+15)
// 	b.SetWall(row, column+15)
// }
//
// func TestBoardValidation(t *testing.T) {
// 	rand.Seed(time.Now().UnixNano())
// 	b := &Board{
// 		Gophers:   make([]int, 0),
// 		Targets:   make([]Target, 0),
// 		locations: make([]int, 0),
// 	}
//
// 	v := rand.Intn(256)
// 	createFailure(b, v)
//
// 	if b.isValid() {
// 		t.Errorf("Board.isValid() did not detect error condition! (block %d)\n", v)
// 	}
// }
//
// func TestBoardMoving(t *testing.T) {
//
// }
