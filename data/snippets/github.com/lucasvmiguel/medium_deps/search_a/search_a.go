package search_a

type SearchA struct{}
type Product struct {
	ID   int
	name string
}

func (s *SearchA) Query() []Product {
	// code…
	return []Product{}
}
