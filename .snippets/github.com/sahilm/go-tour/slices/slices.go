package slices

// Pic creates a slice of slices by iterating dy times, iterating dx times in each iteration.
// The value of elem[y][x] = x*y.
func Pic(dx, dy int) [][]uint8 {
	var pic [][]uint8

	for y := 0; y < dy; y++ {
		xs := []uint8{}
		for x := 0; x < dx; x++ {
			xs = append(xs, uint8(x*y))
		}
		pic = append(pic, xs)
	}
	return pic
}
