package product

import (
	"github.com/lucasvmiguel/medium_deps/product"
	"github.com/lucasvmiguel/medium_deps/search_a"
)

type SearchAProd struct{}

func (s SearchAProd) Query() []product.Product {
	sa := search_a.SearchA{}

	productsSearchA := sa.Query()

	return transformProducts(productsSearchA)
}

func transformProducts(products []search_a.Product) []product.Product {

	//transform search_a pkg products to product pkg products
	return []product.Product{}
}
