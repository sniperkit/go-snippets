// util.go, jpad 2016

package vptree

import (
//	"fmt"
//	"math"
//	"math/rand"
//	"sort"
)

// Info contains structural information about a vp-tree.
type Info struct {
	NNodes   int // number of nodes
	NLeaves  int // number of leaves
	MaxDepth int
}

//// distancePoints ////////////////////////////////////////////////////////////

type distancePoints struct {
	vp     Point   // vantage point
	points []Point // sample points to sort, relative to reference
}

func (p distancePoints) Len() int {
	return len(p.points)
}

func (p distancePoints) Less(i, j int) bool {
	dI := p.vp.Distance(p.points[i])
	dJ := p.vp.Distance(p.points[j])
	return dI < dJ
}

func (p distancePoints) Swap(i, j int) {
	p.points[i], p.points[j] = p.points[j], p.points[i]
}

//// approxInfo ////////////////////////////////////////////////////////////////

type approxInfo struct {
	k     int // number of nearest neighbors to find
	nn    int // maximal number of nodes to visit
	count int // nodes visited
}

//// utilities /////////////////////////////////////////////////////////////////

// selectKth reorders the elements in points, in such a way that the element
// at the k-th position is the element that would be in that position in a
// sorted sequence. This element is returned.
// The other elements are left without any specific order, except that none of
// the elements preceding k-th are greater than it, and none of the elements
// following it are less.
//
// This is Hoare's QuickSelect algorithm, without recursion.
func selectKth(p distancePoints, k int) Point {
	if len(p.points) < k {
		return nil
	}

	from := 0
	to := len(p.points) - 1
	for from < to {
		r := from // read pos
		w := to   // write pos
		mid := p.points[(r+w)/2]
		for r < w {
			if p.points[r].Distance(p.vp) >= mid.Distance(p.vp) {
				// swap large values to the end
				p.points[w], p.points[r] = p.points[r], p.points[w]
				w--
			} else {
				// small values stay where they are
				r++
			}
		}

		// if we stepped up (r++) we must step down
		if p.points[r].Distance(p.vp) > mid.Distance(p.vp) {
			r--
		}

		if k <= r {
			to = r
		} else {
			from = r + 1
		}
	}
	return p.points[k]
}

func indent(width uint) string {
	bytes := make([]byte, width)
	for i := range bytes {
		bytes[i] = ' '
	}
	return string(bytes)
}
