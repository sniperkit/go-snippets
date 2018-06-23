package binaryTree

type Execute func(interface{}) interface{}

//PreOrder traverses binary tree in preorder
func (b *BinaryTree) PreOrder(execute Execute) interface{} {
	value := execute(b.Value())
	if b.Left() != nil {
		b.Left().PreOrder(execute)
	}
	if b.Right() != nil {
		b.Right().PreOrder(execute)
	}
	return value
}

//InOrder traverses binary tree in inorder
func (b *BinaryTree) InOrder(execute Execute) interface{} {
	if b.Left() != nil {
		b.Left().InOrder(execute)
	}
	value := execute(b.Value())
	if b.Right() != nil {
		b.Right().InOrder(execute)
	}
	return value
}

//PosOrder traverses binary tree in postorder
func (b *BinaryTree) PostOrder(execute Execute) interface{} {
	if b.Left() != nil {
		b.Left().PostOrder(execute)
	}
	if b.Right() != nil {
		b.Right().PostOrder(execute)
	}
	return execute(b.Value())
}
