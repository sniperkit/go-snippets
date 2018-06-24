package internal

import (
	"bufio"
	"fmt"
	"io"

	common "github.com/bsm/reason/common/hoeffding"
	"github.com/bsm/reason/core"
	internal "github.com/bsm/reason/internal/hoeffding"
	"github.com/bsm/reason/internal/iocount"
	"github.com/bsm/reason/internal/protoio"
	"github.com/bsm/reason/util"
	"github.com/gogo/protobuf/proto"
)

// NewTree inits a brand-new tree
func NewTree(model *core.Model, target string) *Tree {
	t := &Tree{
		Model:  model,
		Target: target,
	}
	t.Root = t.Add(nil) // init root
	return t
}

// Get retrieves a node by its reference
func (t *Tree) Get(nodeRef int64) *Node {
	pos := int(nodeRef - 1)
	if pos > -1 && pos < len(t.Nodes) {
		return t.Nodes[pos]
	}
	return nil
}

// Set sets a node by reference
func (t *Tree) Set(nodeRef int64, n *Node) {
	pos := int(nodeRef - 1)
	if pos > -1 && pos < len(t.Nodes) {
		t.Nodes[pos] = n
	}
}

// Len returns the number of registered nodes
func (t *Tree) Len() int {
	return len(t.Nodes)
}

// Add adds a new leaf node
func (t *Tree) Add(stats *util.Vector) int64 {
	if stats == nil {
		stats = new(util.Vector)
	}

	leaf := &LeafNode{WeightAtLastEval: stats.Weight()}
	kind := &Node_Leaf{Leaf: leaf}
	node := &Node{Kind: kind, Stats: stats}
	t.Nodes = append(t.Nodes, node)
	return int64(len(t.Nodes))
}

// Split splits an existing leaf node
func (t *Tree) Split(leafRef int64, feature string, pre *util.Vector, post *util.VectorDistribution, pivot float64) {
	if orig := t.Get(leafRef); orig == nil || orig.GetLeaf() == nil {
		return
	}

	split := &SplitNode{
		Feature: feature,
		Pivot:   pivot,
	}

	post.ForEach(func(i int, stats *util.Vector) bool {
		split.Children.SetRef(i, t.Add(stats))
		return true
	})

	kind := &Node_Split{Split: split}
	t.Set(leafRef, &Node{Kind: kind, Stats: pre})
}

// Traverse traverses the tree starting at the given node ID
func (t *Tree) Traverse(x core.Example, nodeRef int64, parent *Node, parentIndex int, forEach func(*Node)) (*Node, int64, *Node, int) {
	node := t.Get(nodeRef)
	if node == nil {
		return node, nodeRef, parent, parentIndex
	}
	if forEach != nil {
		forEach(node)
	}

	if split := node.GetSplit(); split != nil {
		feature := t.Model.Feature(split.Feature)

		if nodeIndex := int(split.childCat(feature, x)); nodeIndex > -1 {
			if childRef := split.Children.GetRef(nodeIndex); childRef > 0 {
				return t.Traverse(x, childRef, node, nodeIndex, forEach)
			}
			return nil, nodeRef, node, nodeIndex
		}
	}
	return node, nodeRef, parent, parentIndex
}

// Prune prunes the leaves of a node recursively
func (t *Tree) Prune(nodeRef int64, parent *Node, isObsolete func(*util.Vector, *util.Vector) bool) {
	// Get the node
	node := t.Get(nodeRef)
	if node == nil {
		return
	}

	// Recurse if a split node
	if split := node.GetSplit(); split != nil {
		split.Children.ForEach(func(_ int, childRef int64) bool {
			t.Prune(childRef, node, isObsolete)
			return true
		})
		return
	}

	// Disable if leaf (with a parent) and obsolete
	if parent != nil && isObsolete(node.Stats, parent.Stats) {
		if leaf := node.GetLeaf(); leaf != nil {
			leaf.Disable()
		}
	}
}

// WriteText appends node information to a text document.
func (t *Tree) WriteText(w io.Writer, nodeRef int64, indent, name string) (nw int64, err error) {
	// Get the node
	node := t.Get(nodeRef)
	if node == nil {
		return
	}

	// Print node stats
	var n int
	n, err = fmt.Fprintf(w, indent+name+" [weight:%.0f]\n", node.Weight())
	nw += int64(n)
	if err != nil {
		return
	}

	// Recurse if a split node
	if split := node.GetSplit(); split != nil {
		feat := t.Model.Feature(split.Feature)
		if feat == nil {
			return
		}

		subIndent := indent + "\t"
		split.Children.ForEach(func(i int, childRef int64) bool {
			var nn int64
			nn, err = t.WriteText(w, childRef, subIndent, internal.FormatNodeCondition(feat, i, split.Pivot))
			nw += nn
			return err == nil
		})
		if err != nil {
			return
		}
	}
	return
}

// WriteDOT appends node information to a DOT document.
func (t *Tree) WriteDOT(w io.Writer, nodeRef int64, name, label string) (nw int64, err error) {
	// Get the node
	node := t.Get(nodeRef)
	if node == nil {
		return
	}

	// Write node data
	var n int
	n, err = fmt.Fprintf(w, `  %s [label="%sweight: %.0f"];`+"\n", name, label, node.Weight())
	nw += int64(n)
	if err != nil {
		return
	}

	// Recurse if a split node
	if split := node.GetSplit(); split != nil {
		feat := t.Model.Feature(split.Feature)
		if feat == nil {
			return
		}

		split.Children.ForEach(func(i int, childRef int64) bool {
			subName := fmt.Sprintf("%s_%d", name, i)

			n, err = fmt.Fprintf(w, "  %s -> %s;\n", name, subName)
			nw += int64(n)
			if err != nil {
				return false
			}

			var nn int64
			nn, err = t.WriteDOT(w, childRef, subName, internal.FormatNodeCondition(feat, i, split.Pivot)+`\n`)
			nw += nn
			return err == nil
		})
		if err != nil {
			return
		}
	}
	return
}

// FilterLeaves finds all leaf-nodes and appends them to dst
func (t *Tree) FilterLeaves(dst []*Node) []*Node {
	for _, n := range t.Nodes {
		switch n.GetKind().(type) {
		case *Node_Leaf:
			dst = append(dst, n)
		}
	}
	return dst
}

// Accumulate collects info stats.
func (t *Tree) Accumulate(nodeRef int64, depth int, info *common.TreeInfo) {
	node := t.Get(nodeRef)
	if node == nil {
		return
	}

	info.NumNodes++
	if depth > info.MaxDepth {
		info.MaxDepth = depth
	}

	if split := node.GetSplit(); split != nil {
		split.Children.ForEach(func(_ int, childRef int64) bool {
			t.Accumulate(childRef, depth+1, info)
			return true
		})
	} else if leaf := node.GetLeaf(); leaf != nil {
		if leaf.IsDisabled {
			info.NumDisabled++
		} else {
			info.NumLearning++
		}
	}
}

// WriteTo writes a tree to a Writer.
func (t *Tree) WriteTo(w io.Writer) (int64, error) {
	wc := &iocount.Writer{W: w}
	wp := &protoio.Writer{Writer: bufio.NewWriter(wc)}

	if err := wp.WriteMessageField(1, t.Model); err != nil {
		return wc.N, err
	}
	if err := wp.WriteStringField(2, t.Target); err != nil {
		return wc.N, err
	}
	if err := wp.WriteVarintField(3, uint64(t.Root)); err != nil {
		return wc.N, err
	}
	for _, node := range t.Nodes {
		if err := wp.WriteMessageField(4, node); err != nil {
			return wc.N, err
		}
	}
	return wc.N, wp.Flush()
}

// ReadFrom reads a tree from a Reader.
func (t *Tree) ReadFrom(r io.Reader) (int64, error) {
	rc := &iocount.Reader{R: r}
	rp := &protoio.Reader{Reader: bufio.NewReader(rc)}

	for {
		tag, wire, err := rp.ReadField()
		if err == io.EOF {
			return rc.N, nil
		} else if err != nil {
			return rc.N, err
		}

		switch tag {
		case 1: // model
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			model := new(core.Model)
			if err := rp.ReadMessage(model); err != nil {
				return rc.N, err
			}
			t.Model = model
		case 2: // target
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			str, err := rp.ReadString()
			if err != nil {
				return rc.N, err
			}
			t.Target = str
		case 3: // root
			if wire != proto.WireVarint {
				return rc.N, proto.ErrInternalBadWireType
			}

			u, err := rp.ReadVarint()
			if err != nil {
				return rc.N, err
			}
			t.Root = int64(u)
		case 4: // nodes
			if wire != proto.WireBytes {
				return rc.N, proto.ErrInternalBadWireType
			}

			node := new(Node)
			if err := rp.ReadMessage(node); err != nil {
				return rc.N, err
			}
			t.Nodes = append(t.Nodes, node)
		default:
			return rc.N, fmt.Errorf("hoeffding: unexpected field tag %d", tag)
		}
	}
}
