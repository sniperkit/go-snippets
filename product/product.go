package product

type Product struct {
	ID   int
	name string
}
type Querier interface {
	Query() []Product
}

func Query(q Querier) []Product {
	return q.Query()
}
