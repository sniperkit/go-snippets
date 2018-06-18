package main

import "fmt"

type node struct {
	left, right *node
	val         int
}

func insert(root *node, val int) *node {
	if root == nil {
		root = &node{val: val}
		return root
	} else if val < root.val {
		root.left = insert(root.left, val)
	} else if val > root.val {
		root.right = insert(root.right, val)
	}

	return root
}

func preorder(root *node) {
	if root != nil {
		preorder(root.left)
		fmt.Println(root.val)
		preorder(root.right)
	}
}

func inorder(root *node) {
        if root != nil {
                fmt.Println(root.val)
                inorder(root.left)
                inorder(root.right)
        }
}

func main() {
	var root *node

	vals := []int{5, 10, 1}

	for _, val := range vals {
		root = insert(root, val)
	}

	preorder(root)
}
