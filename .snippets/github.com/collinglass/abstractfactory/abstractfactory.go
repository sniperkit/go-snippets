package abstractfactory

import (
	"errors"
	"fmt"
)

// Store Factory
//
// implements CreateStore() which
// returns a new store struct
//
// implements CreateMetalProduct() which
// returns a new MetalProduct
//
// implements CreatePaperProduct() which
// returns a new PaperProduct
type StoreFactory interface {
	CreateStore() *Store
	CreateMetalProduct() MetalProduct
	CreatePaperProduct() PaperProduct
}

// Store Product
//
// implements Style() which
// returns a string of the style
//
// implements GetProductId which
// returns the product's Id
//
// implements SetProductId which
// sets a new product Id
type Product interface {
	Style() string
	GetProductId() int
	SetProductId(productId int)
}

// Metal Product
//
// extends Product interface
//
// implements Metal() which
// returns string of the type of metal
type MetalProduct interface {
	Product
	Metal() string
}

// Paper Product
//
// extends Product interface
//
// implements Tree() which
// returns a string of the type of tree
type PaperProduct interface {
	Product
	Tree() string
}

// Store
//
// contains stock array of Products
type Store struct {
	stock *[]Product
}

// NewStore
//
// returns pointer to new Store
func NewStore() *Store {
	return &Store{stock: new([]Product)}
}

// AddStock
//
// adds new Product to store stock
func (store *Store) AddStock(newProduct Product) {
	*store.stock = append(*store.stock, newProduct)
	fmt.Println("Added Stock")
}

// GetStock
//
// returns product with product Id and nil error
//
// or
//
// returns nil and error if product
// does not exist
func (store *Store) GetStock(productId int) (*Product, error) {
	for _, product := range *store.stock {
		if product.GetProductId() == productId {
			return &product, nil
		}
	}
	return nil, errors.New("no product with that id exists")
}

// ErectStore
//
// returns a new Store given a concrete StoreFactory
func ErectStore(factory StoreFactory) *Store {
	store := factory.CreateStore()
	mp := factory.CreateMetalProduct()
	pp := factory.CreatePaperProduct()

	store.AddStock(mp)
	store.AddStock(pp)

	return store
}
