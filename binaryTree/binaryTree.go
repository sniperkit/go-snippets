package binaryTree

type BinaryTree struct {
	left  *BinaryTree
	right *BinaryTree
	value interface{}
}

//Value returns the value of BinaryTree
func (b *BinaryTree) Value() interface{} {
	return b.value
}

//Left returns the left element of BinaryTree
func (b *BinaryTree) Left() *BinaryTree {
	return b.left
}

//Right returns the right element of BinaryTree
func (b *BinaryTree) Right() *BinaryTree {
	return b.right
}

//PutValue changes the value of BinaryTree
func (b *BinaryTree) PutValue(value interface{}) {
	b.value = value
}

//PutLeft creates a new BinaryTree in the left side of BinaryTree
func (b *BinaryTree) PutLeft(value interface{}) *BinaryTree {
	n := new(BinaryTree)
	n.value = value
	n.left = nil
	n.right = nil
	b.left = n
	return n
}

//PutRight creates a new BinaryTree in the right side of BinaryTree
func (b *BinaryTree) PutRight(value interface{}) *BinaryTree {
	n := new(BinaryTree)
	n.value = value
	n.left = nil
	n.right = nil
	b.right = n
	return n
}

//New creates and returns a new BinaryTree instance
func New(value interface{}) *BinaryTree {
	b := new(BinaryTree)
	b.left = nil
	b.right = nil
	b.value = value
	return b
}
