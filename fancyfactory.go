package abstractfactory

//
//
//
//	Fancy Factory
//
//
//
type FancyFactory struct{}

func (ff FancyFactory) CreateStore() *Store {
	return NewStore()
}

func (ff FancyFactory) CreateMetalProduct() MetalProduct {
	return NewFancyMetalProduct()
}

func (ff FancyFactory) CreatePaperProduct() PaperProduct {
	return NewFancyPaperProduct()
}

var productId = 0

//
//
//
//	Fancy Metal Product Line
//
//

type FancyMetalProduct struct {
	style     string
	metal     string
	productId int
}

func NewFancyMetalProduct() MetalProduct {
	productId += 1
	return FancyMetalProduct{productId: productId, style: "Fancy", metal: "Iron"}
}

func (this FancyMetalProduct) Style() string {
	return this.style
}

func (this FancyMetalProduct) Metal() string {
	return this.metal
}

func (this FancyMetalProduct) GetProductId() int {
	return this.productId
}

func (this FancyMetalProduct) SetProductId(productId int) {
	this.productId = productId
}

//
//
//
//	Fancy Paper Product Line
//
//

type FancyPaperProduct struct {
	style     string
	tree      string
	productId int
}

func NewFancyPaperProduct() PaperProduct {
	productId += 1
	return FancyPaperProduct{productId: productId, style: "Fancy", tree: "Oak"}
}

func (this FancyPaperProduct) Tree() string {
	return this.tree
}

func (this FancyPaperProduct) Style() string {
	return this.style
}

func (this FancyPaperProduct) GetProductId() int {
	return this.productId
}

func (this FancyPaperProduct) SetProductId(productId int) {
	this.productId = productId
}
