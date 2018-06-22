package fixtures

var a int

var (
	a, b int = 1, 2
	c        = 3
	_        = "x"
)

func foo() {
	a := 5
	_ = a
	a, b, _ = 1, 2, "x"
}
