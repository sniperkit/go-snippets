package fixtures

const (
	EnumA1 = iota
	EnumA2
)

const (
	EnumB1 = Foo(iota + 1)
	EnumB2
	EnumB2
)

const (
	EnumC1 = Foo(1)
	EnumC2 = Foo(2)
)
