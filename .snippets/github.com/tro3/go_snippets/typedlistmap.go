package snippets

func (x TypeAList) MapToTypeB(fn func(TYPEA) TYPEB) []TYPEB {
	y := make([]TYPEB, 0)
	for _, val := range x {
		y = append(y, fn(val))
	}
	return y
}
