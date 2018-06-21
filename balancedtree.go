/*
<!--
Copyright (c) 2016 Christoph Berger. Some rights reserved.
Use of this text is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "Balancing a binary search tree"
description = "This article describes a basic tree balancing technique, coded in Go, and applied to the binary search tree from last week's article."
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-08-13"
publishdate = "2016-08-13"
domains = ["Algorithms And Data Structures"]
tags = ["Tree", "Balanced Tree", "Binary Tree", "Search Tree"]
categories = ["Tutorial"]
+++

Only a well-balanced search tree can provide optimal search performance. This article adds automatic balancing to the binary search tree from the previous article.

<!--more-->

> Get the Balance Right!
>
> ~ Depeche Mode


## How a tree can get out of balance

As we have seen in last week's article, search performance is best if the tree's height is small. Unfortunately, without any further measure, our simple binary search tree can quickly get out of shape - or never reach a good shape in the first place.

The picture below shows a balanced tree on the left and an extreme case of an unbalanced tree at the right. In the balanced tree, element #6 can be reached in three steps, whereas in the extremely unbalanced case, it takes six steps to find element #6.

![Tree Shapes](BinTreeShapes.png)

Unfortunately, the extreme case can occur quite easily: Just create the tree from a sorted list.

```go
tree.Insert(1)
tree.Insert(2)
tree.Insert(3)
tree.Insert(4)
tree.Insert(5)
tree.Insert(6)
```

According to `Insert`'s logic, each new element is added as the right child of the rightmost node, because it is larger than any of the elements that were already inserted.

We need a way to avoid this.


## A Definition Of "Balanced"

For our purposes, a good working definition of "balanced" is:

> The heights of the two child subtrees of any node differ by at most one.
>
> (Wikipedia: [AVL-Tree](https://en.wikipedia.org/wiki/AVL_tree))

Why "at most one"? Shouldn't we demand *zero* difference for perfect balance? Actually, no, as we can see on this very simple two-node tree:

![Two-node tree](TwoNodeTree.png)

The left subtree is a single node, hence the height is 1, and the right "subtree" is empty, hence the height is zero. There is no way to make both subtrees exactly the same height, except perhaps by adding a third "fake" node that has no other purpose of providing perfect balance. But we would gain nothing from this, so a height difference of 1 is perfectly acceptable.

Note that our definition of *balanced* does not include the *size* of the left and right subtrees of a node. That is, the following tree is completely fine:

![No Weight Balance](BinTreeNoWeightBalance.png)

The left subtree is considerably larger than the right one; yet for either of the two subtrees, any node can be reached with at most four search steps. And the heights of both subtrees differs only by one.


## How to keep a tree in balance

Now that we know what balance means, we need to take care of always keeping the tree in balance. This task consists of two parts: First, we need to be able to detect when a (sub-)tree goes out of balance. And second, we need a way to rearrange the nodes so that the tree is in balance again.


### Step 1. Detecting an imbalance

Balance is related to subtree heights, so we might think of writing a "height" method that descends a given subtree to calculate its height. But this can be come quite costly in terms of CPU time, as these calculations would need to be done repeatedly as we try to determine the balance of each subtee and each subtree's subree, and so on.

Instead, we store a "balance factor" in each node. This factor is an integer that tells the height difference between the node's right and left subtrees, or more formally (this is just maths, no Go code):

    balance_factor := height(right_subtree) - height(left_subtree)

Based on our definition of "balanced", the balance factor of a balanced tree can be -1, 0, or +1. If the balance factor is outside that range (that is, either smaller than -1 or larger than +1), the tree is out of balance and needs to be rebalanced.

After inserting or deleting a node, the balance factors of all affected nodes and parent nodes must be updated.

*For brevity, this article only handles the `Insert` case.*

Here is how `Insert` maintains the balance factors:

1. First, `Insert` descends recursively down the tree until it finds a node `n` to append the new value. `n` is either a leaf (that is, it has no children) or a half-leaf (that is, it has exactly one (direct) child).
2. If `n` is a leaf, adding a new child node increases the height of the subtree `n` by 1. If the child node is added to the left, the balance of `n` changes from 0 to -1. If the child is added to the right, the balance changes from 0 to 1.
2. `Insert` now adds a new child node to node `n`.
3. The height increase is passed back to `n`'s parent node.
4. Depending on whether `n` is the left or the right child, the parent node adjusts its balance accordingly.

**If the balance factor of a node changes to +2 or -2, respectively, we have detected an imbalance.** At this point, the tree needs rebalancing.

HYPE[Balance Factors](BalanceFactors.html)


### Removing the imbalance

Let's assume a node `n` that has one left child and no right child. `n`'s left child has no children; otherwise, the tree at node `n` would already be out of balance. (The following considerations also apply to inserting below the *right* child in a mirror-reversed way, so we can focus on the left-child scenario here.)

Now let's insert a new node below the left child of `n`.

Two scenarios can happen:


#### 1. The new node was inserted as the *left* child of `n`'s left child.

Since `n` has no right children, its balance factor is now -2. (Remember, the balance is defined as "height of right tree minus height of left tree".)
This is an easy case. All we have to do is to "rotate" the tree:

1. Make the left child node the root node.
2. Make the former root node the new root node's right child.

Here is a visualization of these steps (click "Rotate"):

HYPE[Rotation](Rotation.html)

The balance is restored, and the tree's sort order is still intact.

Easy enough, isn't it? Well, only until we look into the other scenario...


#### 2. The new node was inserted as the *right* child of `n`'s left child.

This looks quite similar to the previous case, so let's try the same rotation here. Click "Single Rotation" in the diagram below and see what happens:

HYPE[Double Rotation](DoubleRotation.html)

The tree is again unbalanced; the root node's balance factor changed from -2 to +2. Obviously, a simple rotation as in case 1 does not work here.

Now try the second button, "Double Rotation". Here, the unbalanced node's left subtree is rotated first, and now the situation is similar to case 1. Rotating the tree to the right finally rebalances the tree and retains the sort order.


#### Two more cases and a summary

The two cases above assumed that the unbalanced node's balance factor is -2. If the balance factor is +2, the same cases apply in an analogous way, except that everything is mirror-reversed.


To summarize, here is a scenario where all of the above is included - double rotation as well as reassigning a child node/tree to a rotated node.

HYPE[Re-balance](Rebalance.html)


## The Code

Now, after all this theory, let's see how to add the balancing into the code from the previous article.

First, we set up two helper functions, `min` and `max`, that we will need later.


### Imports, helper functions, and globals
*/

//
package main

import (
	"fmt"
	"strings"
)

// `min` is like math.Min but for int.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// `max` is math.Max for int.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// `Node` gets a new field, `bal`, to store the height difference between the node's subtrees.
type Node struct {
	Value string
	Data  string
	Left  *Node
	Right *Node
	bal   int // height(n.Right) - height(n.Left)
}

/* ### The modified `Insert` function
 */

// `Insert` takes a search value and some data and inserts a new node (unless a node with the given
// search value already exists, in which case `Insert` only replaces the data).
//
// It returns:
//
// * `true` if the height of the tree has increased.
// * `false` otherwise.
func (n *Node) Insert(value, data string) bool {
	// The following actions depend on whether the new search value is equal, less, or greater than
	// the current node's search value.
	switch {
	case value == n.Value:
		n.Data = data
		return false // Node already exists, nothing changes
	case value < n.Value:
		// If there is no left child, create a new one.
		if n.Left == nil {
			// Create a new node.
			n.Left = &Node{Value: value, Data: data}
			// If there is no right child, the new child node has increased the height of this subtree.
			if n.Right == nil {
				// The new left child is the only child.
				n.bal = -1
			} else {
				// There is a left and a right child. The right child cannot have children;
				// otherwise the tree would already have been out of balance at `n`.
				n.bal = 0
			}
		} else {
			// The left child is not nil. Continue in the left subtree.
			if n.Left.Insert(value, data) {
				// If the subtree's balance factor has become either -2 or 2, the subtree must be rebalanced.
				if n.Left.bal < -1 || n.Left.bal > 1 {
					n.rebalance(n.Left)
				} else {
					// If no rebalancing occurred, the left subtree has grown by one: Decrease the balance of the current node by one.
					n.bal--
				}
			}
		}
	// This case is analogous to `value < n.Value`, except that everything is mirrored.
	case value > n.Value:
		if n.Right == nil {
			n.Right = &Node{Value: value, Data: data}
			if n.Left == nil {
				n.bal = 1
			} else {
				n.bal = 0
			}
		} else {
			if n.Right.Insert(value, data) {
				if n.Right.bal < -1 || n.Right.bal > 1 {
					n.rebalance(n.Right)
				} else {
					n.bal++
				}
			}
		}
	}
	if n.bal != 0 {
		return true
	}
	// No more adjustments to the ancestor nodes required.
	return false
}

/* ### The new `rebalance()` method and its helpers `rotateLeft()`, `rotateRight()`, `rotateLeftRight()`, and `rotateRightLeft`.

 **Important note: Many of the assumptions about balances, left and right children, etc, as well as much of the logic usde in the functions below, apply to the `Insert` operation only. For `Delete` operations, different rules and operations apply.** As noted earlier, this article focuses on `Insert` only, to keep the code short and clear.
 */

// `rotateLeft` takes a child node and rotates the child node's subtree to the left.
func (n *Node) rotateLeft(c *Node) {
	fmt.Println("rotateLeft " + c.Value)
	// Save `c`'s right child.
	r := c.Right
	// `r`'s left subtree gets reassigned to `c`.
	c.Right = r.Left
	// `c` becomes the left child of `r`.
	r.Left = c
	// Make the parent node (that is, the current one) point to the new root node.
	if c == n.Left {
		n.Left = r
	} else {
		n.Right = r
	}
	// Finally, adjust the balances. After a single rotation, the subtrees are always of the same height.
	c.bal = 0
	r.bal = 0
}

// `rotateRight` is the mirrored version of `rotateLeft`.
func (n *Node) rotateRight(c *Node) {
	fmt.Println("rotateRight " + c.Value)
	l := c.Left
	c.Left = l.Right
	l.Right = c
	if c == n.Left {
		n.Left = l
	} else {
		n.Right = l
	}
	c.bal = 0
	l.bal = 0
}

// `rotateRightLeft` first rotates the right child of `c` to the right, then `c` to the left.
func (n *Node) rotateRightLeft(c *Node) {
	// `rotateRight` assumes that the left child has a left child, but as part of the rotate-right-left process,
	// the left child of `c.Right` is a leaf. We therefore have to tweak the balance factors before and after
	// calling `rotateRight`.
	// If we did not do that, we would not be able to reuse `rotateRight` and `rotateLeft`.
	c.Right.Left.bal = 1
	c.rotateRight(c.Right)
	c.Right.bal = 1
	n.rotateLeft(c)
}

// `rotateLeftRight` first rotates the left child of `c` to the left, then `c` to the right.
func (n *Node) rotateLeftRight(c *Node) {
	c.Left.Right.bal = -1 // The considerations from rotateRightLeft also apply here.
	c.rotateLeft(c.Left)
	c.Left.bal = -1
	n.rotateRight(c)
}

// `rebalance` brings the (sub-)tree with root node `c` back into a balanced state.
func (n *Node) rebalance(c *Node) {
	fmt.Println("rebalance " + c.Value)
	c.Dump(0, "")
	switch {
	// Left subtree is too high, and left child has a left child.
	case c.bal == -2 && c.Left.bal == -1:
		n.rotateRight(c)
	// Right subtree is too high, and right child has a right child.
	case c.bal == 2 && c.Right.bal == 1:
		n.rotateLeft(c)
	// Left subtree is too high, and left child has a right child.
	case c.bal == -2 && c.Left.bal == 1:
		n.rotateLeftRight(c)
	// Right subtree is too high, and right child has a left child.
	case c.bal == 2 && c.Right.bal == -1:
		n.rotateRightLeft(c)
	}
}

// `Find` stays the same as in the previous article.
func (n *Node) Find(s string) (string, bool) {

	if n == nil {
		return "", false
	}

	switch {
	case s == n.Value:
		return n.Data, true
	case s < n.Value:
		return n.Left.Find(s)
	default:
		return n.Right.Find(s)
	}
}

// `Dump` dumps the structure of the subtree starting at node `n`, including node search values and balance factors.
// Parameter `i` sets the line indent. `lr` is a prefix denoting the left or the right child, respectively.
func (n *Node) Dump(i int, lr string) {
	if n == nil {
		return
	}
	indent := ""
	if i > 0 {
		//indent = strings.Repeat(" ", (i-1)*4) + "+" + strings.Repeat("-", 3)
		indent = strings.Repeat(" ", (i-1)*4) + "+" + lr + "--"
	}
	fmt.Printf("%s%s[%d]\n", indent, n.Value, n.bal)
	n.Left.Dump(i+1, "L")
	n.Right.Dump(i+1, "R")
}

/*
## Tree

Changes to the Tree type:

* `Insert` now takes care of rebalancing the root node if necessary.
* A new method, `Dump`, exist for invoking `Node.Dump`.
* `Delete` is gone.

*/

//
type Tree struct {
	Root *Node
}

func (t *Tree) Insert(value, data string) {
	if t.Root == nil {
		t.Root = &Node{Value: value, Data: data}
		return
	}
	t.Root.Insert(value, data)
	// If the root node gets out of balance,
	if t.Root.bal < -1 || t.Root.bal > 1 {
		t.rebalance()
	}
}

// `Node`'s `rebalance` method is invoked from the parent node of the node that needs rebalancing.
// However, the root node of a tree has no parent node.
// Therefore, `Tree`'s `rebalance` method creates a fake parent node for rebalancing the root node.
func (t *Tree) rebalance() {
	fakeParent := &Node{Left: t.Root, Value: "fakeParent"}
	fakeParent.rebalance(t.Root)
	// Fetch the new root node from the fake parent node
	t.Root = fakeParent.Left
}

func (t *Tree) Find(s string) (string, bool) {
	if t.Root == nil {
		return "", false
	}
	return t.Root.Find(s)
}

func (t *Tree) Traverse(n *Node, f func(*Node)) {
	if n == nil {
		return
	}
	t.Traverse(n.Left, f)
	f(n)
	t.Traverse(n.Right, f)
}

// `Dump` dumps the tree structure.
func (t *Tree) Dump() {
	t.Root.Dump(0, "")
}

/*
### A demo

Using the `Dump` method plus some `fmt.Print...` statements at relevant places, we can watch the code how it inserts new values, rebalancing the subtrees where necessary.

The output of the final `Dump` call should look like this:

```
g[1]
+L--d[0]
    +L--b[0]
        +L--a[0]
        +R--c[0]
    +R--e[1]
        +R--f[0]
+R--i[1]
    +L--h[0]
    +R--k[0]
        +L--j[0]
        +R--l[0]
```

The small letters are the search values. "L" and "R" denote if the child node is a left or a right child. The number in brackets is the balance factor.

If everything works correctly, the `Traverse` method should finally print out the nodes in alphabetical sort order.
*/

//
func main() {
	// The values are sorted in a way that causes two single rotations and a double rotation.
	values := []string{"d", "b", "g", "g", "c", "e", "a", "h", "f", "i", "j", "l", "k"}
	data := []string{"delta", "bravo", "golang", "golf", "charlie", "echo", "alpha", "hotel", "foxtrot", "india", "juliett", "lima", "kilo"}

	tree := &Tree{}
	for i := 0; i < len(values); i++ {
		fmt.Println("Insert " + values[i] + ": " + data[i])
		tree.Insert(values[i], data[i])
		tree.Dump()
		fmt.Println()
	}

	fmt.Print("Sorted values: | ")
	tree.Traverse(tree.Root, func(n *Node) { fmt.Print(n.Value, ": ", n.Data, " | ") })
	fmt.Println()

}

/*
As always, the code is available on GitHub. Using the `-d` flag with `go get` to avoid that the binary gets auto-installed into $GOPATH/bin.

```sh
go get -d github.com/appliedgo/balancedtree
cd $GOPATH/src/github.com/appliedgo/balancedtree
go build
./balancedtree
```

The code is also available on the [Go Playground](https://play.golang.org/p/dd1Z9U90JJ). (Subject to availabilty of the Playground service.)

## Conclusion

Keeping a binary search tree in balance is a bit more involved as it might seem at first. In this article, I have broken down the rebalancing to the bare minimum by removing the `Delete` operation entirely. If you want to dig deeper, here are a couple of useful readings:

[Wikipedia on Tree Rotation](https://en.wikipedia.org/wiki/Tree_rotation): Richly illustrated, concise discussion of the rotation process.

[German Wikipedia on AVL Trees](https://de.wikipedia.org/wiki/AVL-Baum): Sorry, this is German only, but when you scroll down to section 4, "Rebalancierung", there are a couple of detailed diagrams on single and double rotation. Here you can see how the subtree heights change after each rotation.

[GitHub search for Go AVL libs](https://github.com/search?o=desc&q=language%3Ago+avl&s=stars&type=Repositories&utf8=%E2%9C%93): For advanced study :)


That's it. Happy tree planting!

*/
