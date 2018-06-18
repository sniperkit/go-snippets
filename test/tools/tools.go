package tools

type cnt struct {
	N int
}

func New() *cnt {
	return new(cnt)
}