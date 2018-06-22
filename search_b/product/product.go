package product

import (
	"github.com/lucasvmiguel/medium_deps/search_b"

	"github.com/lucasvmiguel/medium_deps/product"
)

type SearchBProd struct{}

func (s SearchBProd) Query() []product.Product {
	sb := search_b.SearchB{}

	productsSearchB := sb.Query()

	return transformProducts(productsSearchB)
}

func transformProducts(products []search_b.Product) []product.Product {

	//transform search_b pkg products to product pkg products
	return []product.Product{}
}
