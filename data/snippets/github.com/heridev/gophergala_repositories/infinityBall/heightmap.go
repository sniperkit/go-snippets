package main

import (
	"fmt"
	"math/rand"
	"math"
	"time"
)

const max float64 = 1.0

const min int = 4
/*
func main() {

	//rand.Seed( time.Now().UTC().UnixNano())

	//size := 9

	s := 17
	m := make([][]float64, s)
	for x:=0; x<s;x++ {
		m[x] = make([]float64, s)
		for y :=0; y<s; y++ {
			m[x][y] = float64(x+y)
		}
	}

	printArray(m)
	tm := translateHM(&m, -1, 1)
	printArray(*tm)
	
	//var A [size][size] float64

	//update(&A,Pair{(size-1)/2,(size-1)/2},(size-1)/2)
	//hm := createHM(size)

	//printArray(*hm)

	
}
*/
type Pair struct {
	row int
	col int
}

func createHM(size int) [][]float64 {
	//const s = size
	rand.Seed( time.Now().UTC().UnixNano())
	//var A [size][size] float64
	A := hmHelper(size)

	update(&A,Pair{(size-1)/2,(size-1)/2},(size-1)/2)
	return A
}

func translateHM(hm *[][]float64, ox , oy int) *[][]float64 {
	size := len(*hm)
	
	
	transHm := hmHelper(size)
	for x:=0+ox; x < size; x++ {
		for y:=0+oy; y<size;y++ {
			if !(x < 0 || y < 0 || (x-ox)>=size || (y-oy)>=size) {
				transHm[x-ox][y-oy] = (*hm)[x][y]
			}
		} 
	}

	printArray(transHm)
	
	update(&transHm,Pair{(size-1)/2,(size-1)/2},(size-1)/2)	
	return &transHm
}

func hmHelper(size int) [][]float64 {

	hm := make([][]float64, size)
	for x:=0; x<size; x++ {
		hm[x] = make([]float64, size)
		for y:=0; y<size; y++ {
			hm[x][y] = math.NaN()
		}
	}
	hm[0][0] = 0
	hm[0][size-1] = 0
	hm[size-1][0] = 0
	hm[size-1][size-1] = 0
	return hm
}

func get4Mid(p1 Pair, p2 Pair, p3 Pair, p4 Pair) Pair{
	return Pair{(p1.row+p2.row+p3.row+p4.row)/4,(p1.col+p2.col+p3.col+p4.col)/4}
}

func getVal(a *[][]float64, p Pair) float64 {
	if math.IsNaN((*a)[p.row][p.col]) {
		return 0.0
	}
	return (*a)[p.row][p.col]
}

func setVal(a *[][]float64, p Pair, val float64) {
	if math.IsNaN((*a)[p.row][p.col]) {
		(*a)[p.row][p.col] = val
	}
}

func isOut(p Pair, size int) bool{
	return p.row < 0 || p.row >= size || p.col < 0 || p.col >= size
}

func diamond(a* [][]float64, p Pair, l int){

		size := len(*a)
		
		up := Pair{p.row-l,p.col}
		down := Pair{p.row+l,p.col}
		left := Pair{p.row,p.col-l}
		right := Pair{p.row,p.col+l}

		diamondAvg := 0.0
		if l > min {
			diamondAvg += (rand.Float64()*max-max/2)*float64(l)/float64(size)
		}


		if isOut(up, size) {
			diamondAvg += (getVal(a, down) + getVal(a, left) + getVal(a, right))/3
		} else if isOut(down, size) {
			diamondAvg += (getVal(a, up) + getVal(a, left) + getVal(a, right))/3
		} else if isOut(left, size) {
			diamondAvg += (getVal(a, up) + getVal(a, down) + getVal(a, right))/3
		} else if isOut(right, size) {
			diamondAvg += (getVal(a, up) + getVal(a, down) + getVal(a, left))/3
		} else {
			diamondAvg += (getVal(a, up) + getVal(a, down) + getVal(a, left) + getVal(a, right))/4
		}

		setVal(a, p, diamondAvg)
}

func square(a* [][]float64, p Pair, l int){
		size := len(*a)

		topLeft			:= Pair{p.row-l,p.col-l}
		topRight		:= Pair{p.row-l,p.col+l}
		bottomLeft	:= Pair{p.row+l,p.col-l}
		bottomRight := Pair{p.row+l,p.col+l}

		squareAvg := (getVal(a, topLeft) + getVal(a, topRight) + getVal(a, bottomLeft) + getVal(a, bottomRight))/4
		//perturburation := rand.Float64()*max
		if l > min {
			squareAvg += (rand.Float64()*max-max/2)*float64(l)/float64(size)
		}
		setVal(a, p, squareAvg)
}

func update(a *[][]float64, p Pair, l int) {
		//fmt.Println(l, p)
		
		if l > 0 {
			square(a,p,l)
			diamond(a,Pair{p.row+l,p.col},l)
			diamond(a,Pair{p.row-l,p.col},l)
			diamond(a,Pair{p.row,p.col-l},l)
			diamond(a,Pair{p.row,p.col+l},l)

			update(a,Pair{p.row-l/2,p.col-l/2},l/2) // TopLeft
			update(a,Pair{p.row-l/2,p.col+l/2},l/2) // TopRight
			update(a,Pair{p.row+l/2,p.col-l/2},l/2) // BottomLeft
			update(a,Pair{p.row+l/2,p.col+l/2},l/2) // BottomRight
		}
}


/*
func isPowerOfTwo (num int) bool {
	for (((num % 2) == 0) && num > 1) 
		num /= 2;
	return (num == 1);
}
*/

func printArray(a[][] float64) {
	size := len(a)

	fmt.Print("[")
	for i := 0; i < size-1; i++ {
		fmt.Print("[")
		for j:= 0; j < size-1; j++ {
			fmt.Print(a[i][j], ",")
		}
		fmt.Println(a[i][size-1], "],")
	}
	fmt.Print("[")
	for j:= 0; j < size-1; j++ {
		fmt.Print(a[size-1][j], ",")
	}
	fmt.Println(a[size-1][size-1], "]]")
}

// 	var a [size+1][size+1]float32
// 	a[0][0] = 1
// 	a[0][size] = 1
// 	a[size][0] = 1
// 	a[size][size] = 1
// 
// 	updateSquare(&a, Square{Pair{0,0},Pair{size,0},Pair{0,size},Pair{size,size}}, size)
// 	printArray(a)
// 
// }
// 
// 
// func get2Mid(p1 Pair, p2 Pair) Pair{
// 	return Pair{(p1.x+p2.x)/2,(p1.y+p2.y)/2}
// }
// 
// 
// func get2Avg(aPtr *[size+1][size+1]float32, p1 Pair, p2 Pair) float32 {
// 	return ((*aPtr)[p1.x][p1.y] + (*aPtr)[p2.x][p2.y])/2
// }
// 
// func get4Avg(aPtr *[size+1][size+1]float32, p1 Pair, p2 Pair, p3 Pair, p4 Pair) float32 {
// 	return ((*aPtr)[p1.x][p1.y] + (*aPtr)[p2.x][p2.y] + (*aPtr)[p3.x][p3.y] + (*aPtr)[p4.x][p4.y])/4
// }
// 
// type Square struct {
// 	nw Pair
// 	ne Pair
// 	sw Pair
// 	se Pair
// }
// 
// func equals(p1 Pair, p2 Pair) bool{
// 	return p1.x == p2.x && p1.y == p2.y
// }
// func updateSquare(aPtr *[size+1][size+1]float32, s Square, w int) {
// 
// 
// 	m := get4Mid(s.nw,s.ne,s.sw,s.se)
// 	u := get2Mid(s.nw, s.ne)
// 	d := get2Mid(s.sw, s.se)
// 	l := get2Mid(s.nw, s.sw)
// 	r := get2Mid(s.ne, s.se)
// 
// 	(*aPtr)[m.x][m.y] = (get4Avg(aPtr, s.nw, s.ne, s.sw, s.se) + rand.Float32()-0.5)/2
// 	(*aPtr)[u.x][u.y] = get2Avg(aPtr, s.nw, s.ne) //+ rand.Float32()*max
// 	(*aPtr)[d.x][d.y] = get2Avg(aPtr, s.sw, s.se) //+ rand.Float32()*max
// 	(*aPtr)[l.x][l.y] = get2Avg(aPtr, s.nw, s.sw) //+ rand.Float32()*max
// 	(*aPtr)[r.x][r.y] = get2Avg(aPtr, s.ne, s.se) //+ rand.Float32()*max
// 
// 	if w > 1 {
// 		updateSquare(aPtr, Square{s.nw, u, m, l}, w/2)
// 		updateSquare(aPtr, Square{u, s.ne, r, m}, w/2)
// 		updateSquare(aPtr, Square{l, m, d, s.sw}, w/2)
// 		updateSquare(aPtr, Square{m, r, s.se, d}, w/2)
// 	}
// 
// }
// 

