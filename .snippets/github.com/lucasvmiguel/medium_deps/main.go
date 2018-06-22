package main

import (
	"fmt"

	"github.com/lucasvmiguel/medium_deps/product"
	search_a_prod "github.com/lucasvmiguel/medium_deps/search_a/product"
	search_b_prod "github.com/lucasvmiguel/medium_deps/search_b/product"
)

func main() {
	searchA := search_a_prod.SearchAProd{}

	productsFromSearchA := product.Query(searchA)

	searchB := search_b_prod.SearchBProd{}

	productsFromSearchB := product.Query(searchB)

	fmt.Println(productsFromSearchA)
	fmt.Println(productsFromSearchB)
}
