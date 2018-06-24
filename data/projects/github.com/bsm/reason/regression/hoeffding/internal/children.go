package internal

import (
	"github.com/bsm/reason/internal/sparsedense"
)

// ForEach iterates over a node-set
func (m *SplitNode_Children) ForEach(iter func(int, int64) bool) {
	if m.Dense != nil {
		for i, nodeRef := range m.Dense {
			if nodeRef > 0 {
				if !iter(i, nodeRef) {
					break
				}
			}
		}
	} else if m.Sparse != nil {
		for i, nodeRef := range m.Sparse {
			if !iter(int(i), nodeRef) {
				break
			}
		}
	}
}

// Len returns the size
func (m *SplitNode_Children) Len() int {
	if m.Dense != nil {
		n := 0
		m.ForEach(func(_ int, _ int64) bool { n++; return true })
		return n
	}
	return len(m.Sparse)
}

// GetRef returns a single nodeRef at index
func (m *SplitNode_Children) GetRef(index int) int64 {
	if index < 0 {
		return 0
	}

	if m.Dense != nil && index < len(m.Dense) {
		return m.Dense[index]
	} else if m.Sparse != nil {
		return m.Sparse[int64(index)]
	}
	return 0
}

// SetRef stores a nodeRef at an index
func (m *SplitNode_Children) SetRef(index int, nodeRef int64) {
	if index < 0 {
		return
	}

	if m.Dense != nil {
		m.setDense(index, nodeRef)
		return
	}

	if n := int64(index + 1); n > m.SparseCap {
		m.SparseCap = n
	}
	if m.Sparse == nil {
		m.Sparse = make(map[int64]int64, 1)
	}
	m.Sparse[int64(index)] = nodeRef
	if sparsedense.BetterOffDense(len(m.Sparse), int(m.SparseCap)) {
		m.convertToDense()
	}
}

func (m *SplitNode_Children) setDense(index int, nodeRef int64) {
	if n := index + 1; n > cap(m.Dense) {
		dense := make([]int64, n, 2*n)
		copy(dense, m.Dense)
		m.Dense = dense
	} else if n > len(m.Dense) {
		m.Dense = m.Dense[:n]
	}
	m.Dense[index] = nodeRef
}

func (m *SplitNode_Children) convertToDense() {
	dense := make([]int64, int(m.SparseCap))
	m.ForEach(func(i int, nodeRef int64) bool {
		dense[i] = nodeRef
		return true
	})
	m.Dense = dense
	m.Sparse = nil
	m.SparseCap = 0
}
