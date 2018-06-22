package fixtures

func foo() {
	return
}

func bar(a int, b string) (string, int, error) {
	return b, a, nil
}

func baz(a, b int, _ int, args ...string) {
	return
}

func (A) foo() {
	return
}

func (a A) bar() (i1, i2 int, s string) {
	return 1, 2, "x"
}
