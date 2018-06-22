package fixtures

func slices() {
	var a = [...]int{1, 2}
	b := []int{3, 4}
	d := [][]int{
		{6},
	}

	b = append(b, 5)
	b = append(b, a[:]...)
	b = append(b, d[0]...)

	b = b[1:]
	b = b[1:2:2]
	b[0] += len(b) + cap(b)

	c := b
	c = nil
	c = make([]int, 3)
	c = make([]int, 0, 3)

	for i := range b {
		_ = b[i]
	}
	for i, v := range b {
		_, _ = i, v
	}
}
