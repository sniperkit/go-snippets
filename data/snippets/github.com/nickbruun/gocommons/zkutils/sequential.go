package zkutils

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

// Sequence node.
type SequenceNode struct {
	// Name.
	Name string

	// Sequence number.
	SequenceNumber int32
}

func (n SequenceNode) Equals(b SequenceNode) bool {
	return n.SequenceNumber == b.SequenceNumber && n.Name == b.Name
}

// Ascendingly sorted sequence nodes.
type sequenceNodesAscendingly []SequenceNode

func (l sequenceNodesAscendingly) Len() int {
	return len(l)
}

func (l sequenceNodesAscendingly) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l sequenceNodesAscendingly) Less(i, j int) bool {
	return l[i].SequenceNumber < l[j].SequenceNumber
}

// Ascendingly sorted sequence nodes with negative numbers last.
type sequenceNodesAscendinglyNegativeLast []SequenceNode

func (l sequenceNodesAscendinglyNegativeLast) Len() int {
	return len(l)
}

func (l sequenceNodesAscendinglyNegativeLast) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l sequenceNodesAscendinglyNegativeLast) Less(i, j int) bool {
	isn := int64(l[i].SequenceNumber)
	jsn := int64(l[j].SequenceNumber)

	if isn < 0 {
		isn = int64(2147483647+2147483649) + isn
	}

	if jsn < 0 {
		jsn = int64(2147483647+2147483649) + jsn
	}

	return isn < jsn
}

// Get expression for a sequence node.
func sequenceNodeExpr(prefix string) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(`.*?%s(-?\d+)$`, regexp.QuoteMeta(prefix)))
}

// Parse a sequence node.
//
// Returns ErrNodeNotMatch if the node does not match the provided expression.
func parseSequenceNode(name string, expr *regexp.Regexp) (SequenceNode, error) {
	groups := expr.FindStringSubmatch(name)
	if groups == nil || len(groups) == 0 {
		return SequenceNode{}, ErrNodeNotMatch
	}

	idx, err := strconv.ParseInt(groups[1], 10, 32)
	if err != nil {
		return SequenceNode{}, err
	}

	return SequenceNode{
		Name:           name,
		SequenceNumber: int32(idx),
	}, nil
}

// Parse a sequence node.
//
// Returns ErrNodeNotMatch if the node does not match the provided prefix or is
// not a sequence node.
func ParseSequenceNode(name, prefix string) (SequenceNode, error) {
	return parseSequenceNode(name, sequenceNodeExpr(prefix))
}

// Parse a list of sequence nodes.
//
// Ignores any node that is not a sequentially numbered node. If a prefix is
// provided, any node where the sequence number is not immediately preceeded by
// the prefix is also ignored.
func ParseSequenceNodes(names []string, prefix string) (nodes []SequenceNode) {
	expr := sequenceNodeExpr(prefix)
	nodes = make([]SequenceNode, 0, len(names))

	for _, n := range names {
		if sn, err := parseSequenceNode(n, expr); err == nil {
			nodes = append(nodes, sn)
		}
	}

	return
}

// Sort a list of sequence nodes.
//
// Sorts the sequence nodes in a semi-overflow safe manner. Sequence numbers
// will follow the following order:
//
//     0 .. 2147483647
//     -2147483648 .. -1
//     0 .. 2147483647
//     ..
//
// Sorting makes the assumption, that sequence numbers will never be too far
// apart in the natural, overflowing sequence. Thus, the sorting is performed
// with the following ordering:
//
// - Ordered ascendingly if the set of sequence numbers are in the range
//   [0 ; 2147483647]
// - Ordered ascendingly with negative sequence numbers ordered after the
//   positive sequence numbers if the set of sequence numbers are in the range
//   [0 ; 2147483647] and [-2147483648 ; -1073741824].
// - Ordered ascendingly if the set of sequence numbers are in the range
//   [-2147483648 ; 1073741824].
func SortSequenceNodes(nodes []SequenceNode) {
	if len(nodes) < 2 {
		return
	}

	// Determine the extent of the nodes.
	snMin := nodes[0].SequenceNumber
	snMax := nodes[0].SequenceNumber

	for i, n := range nodes {
		if i > 0 {
			if n.SequenceNumber < snMin {
				snMin = n.SequenceNumber
			} else if n.SequenceNumber > snMax {
				snMax = n.SequenceNumber
			}
		}
	}

	// Order ascendingly if this is a simple case.
	if snMin >= int32(-1073741824) || snMax < 0 {
		sort.Sort(sequenceNodesAscendingly(nodes))
	} else {
		// Order ascendingly with negative numbers ordered last.
		sort.Sort(sequenceNodesAscendinglyNegativeLast(nodes))
	}
}
