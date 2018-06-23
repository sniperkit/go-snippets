package fixtures

func pointers() {
	var a *int
	b := 0
	a = &b
	*a = 1
	a = new(int)
}
