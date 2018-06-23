package twothree

import (
	"strconv"
)

type tree struct {
	root node
}

/// Creates new EmptyTree
func NewTree() *tree {
	return &tree{}
}

func (t *tree) Insert(value int) {
	if t.root == nil {
		t.root = two{k: value}
	} else {
		t.root = t.root.insert(value)
	}
}

/// Checks if the tree contains the given value
func (t *tree) Member(value int) bool {
	return t.root != nil && t.root.member(value)
}

/// Calculates the number of values in the tree
func (t *tree) Length() int {
	if t.root == nil {
		return 0
	} else {
		return t.root.length()
	}
}

/// Function for representing the tree as set e.g "{1,5,8,34,123,56}"
func (t *tree) String() string {
	if t.root != nil {
		return "{" + t.root.String() + "}"
	} else {
		return "{}"
	}
}

type node interface {
	ins(int) node
	insert(int) node
	grow() node
	restore() node
	length() int
	member(int) bool
	String() string
}

type two struct {
	k           int
	left, right node
}

func newTwo(left node, k int, right node) *two {
	return &two{left: left, k: k, right: right}
}

type three struct {
	k1, k2           int
	left, mid, right node
}

type four struct {
	k1, k2, k3              int
	left, mid1, mid2, right node
}

func (t two) member(value int) bool {
	switch {
	case value == t.k:
		return true
	case value < t.k:
		return t.left != nil && t.left.member(value)
	case value > t.k:
		return t.right != nil && t.right.member(value)
	default:
		println("Error! Default case should not be reached")
		return false
	}
}

func (t three) member(value int) bool {
	switch {
	case value == t.k1 || value == t.k2:
		return true
	case value < t.k1:
		return t.left != nil && t.left.member(value)
	case value < t.k2:
		return t.mid != nil && t.mid.member(value)
	case value > t.k2:
		return t.right != nil && t.right.member(value)
	default:
		println("Error! Default case should not be reached")
		return false
	}
}

func (t four) member(value int) bool {
	println("Error! Function four.member must not be called")
	return false
}

func (t two) length() int {
	length := 1
	if t.left != nil {
		length += t.left.length()
	}
	if t.right != nil {
		length += t.right.length()
	}
	return length
}

func (t three) length() int {
	length := 2
	if t.left != nil {
		length += t.left.length()
	}
	if t.mid != nil {
		length += t.mid.length()
	}
	if t.right != nil {
		length += t.right.length()
	}
	return length
}

func (t four) length() int {
	println("Error! Function four.length must not be called on struct four")
	return 0
}

func (t two) insert(value int) node {
	return t.ins(value).restore().grow()
}

func (t three) insert(value int) node {
	return t.ins(value).restore().grow()
}

func (t four) insert(value int) node {
	return t.ins(value).restore().grow()
}

func (t two) ins(value int) node {
	switch {
	case t.k == value:
		return t
	case t.left == nil && t.right == nil:
		if t.k > value {
			return three{k1: value, k2: t.k}
		}
		if t.k < value {
			return three{k1: t.k, k2: value}
		}
	case t.left != nil && t.right != nil:
		if t.k > value {
			t.left = t.left.ins(value).restore()
		}
		if t.k < value {
			t.right = t.right.ins(value).restore()
		}
		return t
	}
	return nil
}

func (t three) ins(value int) node {
	switch {
	case t.k1 == value || t.k2 == value:
		return t
	case t.left == nil && t.mid == nil && t.right == nil:
		if value < t.k1 {
			return four{k1: value, k2: t.k1, k3: t.k2}
		} else if value < t.k2 {
			return four{k1: t.k1, k2: value, k3: t.k2}
		} else if value > t.k2 {
			return four{k1: t.k1, k2: t.k2, k3: value}
		}
		return t
	case t.left != nil && t.mid != nil && t.right != nil:
		if value < t.k1 {
			t.left = t.left.ins(value).restore()
		} else if value < t.k2 {
			t.mid = t.mid.ins(value).restore()
		} else if value > t.k2 {
			t.right = t.right.ins(value).restore()
		}
		return t
	}
	return nil
}

func (t four) ins(value int) node {
	println("Error! Function four.ins must not be called")
	return t
}

func (t two) restore() node {
	if left, ok := t.left.(four); ok {
		return three{left: newTwo(left.left, left.k1, left.mid1), k1: left.k2, mid: newTwo(left.mid2, left.k3, left.right), k2: t.k, right: t.right}

	}
	if right, ok := t.right.(four); ok {
		return three{left: t.left, k1: t.k, mid: newTwo(right.left, right.k1, right.mid1), k2: right.k2, right: newTwo(right.mid2, right.k3, right.right)}
	}
	return t
}

func (t three) restore() node {
	if left, ok := t.left.(four); ok {
		return four{left: newTwo(left.left, left.k1, left.mid1), k1: left.k2, mid1: newTwo(left.mid2, left.k3, left.right), k2: t.k1, mid2: t.mid, k3: t.k2, right: t.right}
	}
	if mid, ok := t.mid.(four); ok {
		return four{left: t.left, k1: t.k1, mid1: newTwo(mid.left, mid.k1, mid.mid1), k2: mid.k2, mid2: newTwo(mid.mid2, mid.k3, mid.right), k3: t.k2, right: t.right}
	}
	if right, ok := t.right.(four); ok {
		return four{left: t.left, k1: t.k1, mid1: t.mid, k2: t.k2, mid2: newTwo(right.left, right.k1, right.mid1), k3: right.k2, right: newTwo(right.mid2, right.k3, right.right)}
	}
	return t
}

func (t four) restore() node {
	return t
}

// Functions for splitting four node in two nodes

func (t four) grow() node {
	return newTwo(newTwo(t.left, t.k1, t.mid1), t.k2, newTwo(t.mid2, t.k3, t.right))
}

func (t three) grow() node {
	return t
}

func (t two) grow() node {
	return t
}

/// Functions for printing the nodes as list joined with ','
/// e.g "1,5,8,34,123,56"

func (t two) String() string {
	result := ""
	if t.left != nil {
		result += t.left.String() + ","
	}
	result += strconv.Itoa(t.k)
	if t.right != nil {
		result += "," + t.right.String()
	}
	return result
}

func (t three) String() string {
	result := ""
	if t.left != nil {
		result += t.left.String() + ","
	}
	result += strconv.Itoa(t.k1) + ","
	if t.mid != nil {
		result += t.mid.String() + ","
	}
	result += strconv.Itoa(t.k2)
	if t.right != nil {
		result += "," + t.right.String()
	}
	return result
}

func (t four) String() string {
	result := ""
	if t.left != nil {
		result += t.left.String() + ","
	}
	result += strconv.Itoa(t.k1) + ","
	if t.mid1 != nil {
		result += t.mid1.String() + ","
	}
	result += strconv.Itoa(t.k2) + ","
	if t.mid2 != nil {
		result += t.mid2.String() + ","
	}
	result += strconv.Itoa(t.k3)
	if t.right != nil {
		result += "," + t.right.String()
	}
	return result
}
