package quorem

func quoRem1(x, y int) (quo, rem int) {
	quo = x / y
	rem = x % y
	return
}

func quoRem2(x, y int) (quo, rem int) {
	quo = x / y
	rem = x - quo*y
	return
}

var QuoRem = quoRem2 // faster function
