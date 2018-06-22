package similarity

import (
	"fmt"
	"math"
	"strings"
)

var similarity float64
var datasource *PrefixMap

// Match ...
type Match struct {
	Value      string
	Similarity float64
}

// Print ...
func (match *Match) Print() {
	fmt.Printf("match: \t%s\tsimilarity: %.2f\t\n", match.Value, match.Similarity)
}

// Exec execute the search for similarity
func Exec(dictionary []string, input string, float float64) []Match {
	datasource = New()
	similarity = float
	// here we populate the datasource
	for _, words := range dictionary {
		parts := strings.Split(strings.ToLower(words), " ")
		for _, part := range parts {
			datasource.Insert(part, words)
		}
	}

	values := datasource.GetByPrefix(strings.ToLower(input))

	results := make([]Match, 0, len(values))
	for _, v := range values {
		value := v.(string)
		s := ComputeSimilarity(len(value), len(input), LevenshteinDistance(strings.ToLower(value), strings.ToLower(input)))
		if s >= similarity {
			m := Match{value, s}
			fmt.Printf("match: \t%s\tsimilarity: %.2f\t\n", m.Value, m.Similarity)
			results = append(results, m)
		}
	}

	fmt.Printf("Result for target similarity: %.2f\n", similarity)
	return results
}

// ComputeSimilarity ...
func ComputeSimilarity(w1Len, w2Len, ld int) float64 {
	maxLen := math.Max(float64(w1Len), float64(w2Len))

	return 1.0 - float64(ld)/float64(maxLen)
}

//LevenshteinDistance ...
func LevenshteinDistance(source, destination string) int {
	vec1 := make([]int, len(destination)+1)
	vec2 := make([]int, len(destination)+1)

	w1 := []rune(source)
	w2 := []rune(destination)

	// initializing vec1
	for i := 0; i < len(vec1); i++ {
		vec1[i] = i
	}

	// initializing the matrix
	for i := 0; i < len(w1); i++ {
		vec2[0] = i + 1

		for j := 0; j < len(w2); j++ {
			cost := 1
			if w1[i] == w2[j] {
				cost = 0
			}
			min := minimum(vec2[j]+1,
				vec1[j+1]+1,
				vec1[j]+cost)
			vec2[j+1] = min
		}

		for j := 0; j < len(vec1); j++ {
			vec1[j] = vec2[j]
		}
	}

	return vec2[len(w2)]
}

func minimum(value0 int, values ...int) int {
	min := value0
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

// Node is a single node within
// the map
type Node struct {
	// true if this node is a leaf node
	IsLeaf bool

	// the reference to the parent node
	Parent *Node

	// the children nodes
	Children []*Node

	// private
	key    string
	isRoot bool
	data   []interface{}
}

// PrefixMap type
type PrefixMap Node

func newNode() (m *Node) {
	m = new(Node)

	m.IsLeaf = false
	m.Parent = nil

	return
}

// New returns a new empty map
func New() *PrefixMap {
	m := newNode()
	m.isRoot = true

	return (*PrefixMap)(m)
}

// Depth returns the depth of the
// current node within the map
func (m *Node) Depth() int {
	depth := 0
	parent := m.Parent
	for parent != nil {
		depth++
		parent = parent.Parent
	}

	return depth
}

// This method traverses the map to find an appropriate node
// for the given key. Optionally, if no node is found, one is created.
//
// Returns an additional bool indicating if the node key is an exact match
// for the given key parameter or false if it is the closest match found
// Algorithm: BFS
func (m *Node) nodeForKey(key string, createIfMissing bool) (*Node, bool) {
	var lastNode *Node = nil
	var currentNode = m

	// holds the next children to explore
	var children []*Node

	var lcpI int // last lcp index

	for currentNode != nil && len(key) > 0 {
		// root is special case for us
		// since it doesn't hold any information
		if currentNode.isRoot {
			if len(currentNode.Children) == 0 {
				break
			}
			children = currentNode.Children
			if len(children) > 0 {
				currentNode, children = children[0], children[1:]
				continue
			}
			break
		}

		lcpI = lcpIndex(key, currentNode.key)

		// current node is not the one
		if lcpI == -1 {
			if len(children) > 0 {
				currentNode, children = children[0], children[1:]
				continue
			}
			break
		}

		// key matches current node: returning it
		if lcpI == len(key)-1 && lcpI == len(currentNode.key)-1 {
			return currentNode, true
		}
		key = key[lcpI+1:]

		// in this case the key we are looking for is a substring
		// of the current node key so we need to split the node
		if len(key) == 0 {
			if createIfMissing == true {
				currentNode.split(lcpI + 1)
				return currentNode, true
			}
			return currentNode, false
		}

		// current node key is a substring of the requested
		// key so we go deep in the tree from here
		if lcpI == len(currentNode.key)-1 {
			lastNode = currentNode
			children = currentNode.Children
			if len(children) == 0 {
				break
			}
			currentNode, children = children[0], children[1:]
			continue
		}

		if createIfMissing == true {
			currentNode.split(lcpI + 1)
			lastNode = currentNode
			break
		}

		// Important Case: given key partially matches with
		// current node key.
		//
		// This means we have to split the existing node
		// into two nodes and append the new content accordingly
		//
		// e.g.
		// Key to be inserted: 'string'
		// Node found: 'stringmap'
		// => we need to split 'stringmap' into 'string' and 'map'
		//    in order to be able to set a value for the key 'string'
		//    and still maintain the value(s) associated with 'stringmap'
		//    in the new 'map' node
		//
		// States can be represented as following:
		//
		// State 1:
		//
		//         o (root)
		//         |
		//         o (stringmap) = (some values)
		//
		// State 2 after inserting key 'string' into the map:
		//
		//        o (root)
		//        |
		//        o (string) = (some values associated with 'string' key)
		//        |
		//        o (map)    = (some values associated with 'stringmap' key)
	}

	if createIfMissing == true {
		newNode := newNodeWithKey(key)
		if lastNode == nil {
			return m.appendNode(newNode), true
		}
		return lastNode.appendNode(newNode), true
	}

	return lastNode, false
}

func (m *Node) split(index int) {
	rightKey := m.key[index:]
	leftKey := m.key[:index]
	subNode := m.copyNode()
	subNode.key = rightKey
	subNode.Parent = m
	subNode.IsLeaf = true

	// adjusting children parent
	for _, child := range subNode.Children {
		child.Parent = subNode
	}

	m.key = leftKey
	m.Children = []*Node{subNode}
	m.data = []interface{}{}
	m.IsLeaf = false
}

func (m *Node) copyNode() *Node {
	n := &Node{}
	*n = *m
	return n
}

func newNodeWithKey(key string) *Node {
	n := newNode()
	n.key = key
	return n
}

func (m *Node) appendNode(n *Node) *Node {
	m.Children = append(m.Children, n)
	n.IsLeaf = true
	n.Parent = m
	return n
}

// Insert inserts a new value in the map for the specified key
// If the key is already present in the map, the value is appended
// to the values list associated with the given key
func (m *PrefixMap) Insert(key string, values ...interface{}) {
	mNode := (*Node)(m)
	n, _ := mNode.nodeForKey(key, true)
	n.data = append(n.data, values...)
}

// Replace replaces the value(s) for the given key in the map
// with the give ones. If no such key is present, this method
// behaves the same as Insert
func (m *PrefixMap) Replace(key string, values ...interface{}) {
	mNode := (*Node)(m)
	n, _ := mNode.nodeForKey(key, true)
	n.data = values
}

// Contains checks if the given key is present in the map
// In this case, an exact match case is considered
// If you're interested in prefix-based check: ContainsPrefix
func (m *PrefixMap) Contains(key string) bool {
	mNode := (*Node)(m)
	retrievedNode, exactMatch := mNode.nodeForKey(key, false)
	return retrievedNode != nil && exactMatch
}

// Get returns the data associated with the given key in the map
// or nil if no such key is present in the map
func (m *PrefixMap) Get(key string) []interface{} {
	mNode := (*Node)(m)
	retrievedNode, exactMatch := mNode.nodeForKey(key, false)
	if !exactMatch {
		return nil
	}

	return retrievedNode.data
}

// GetByPrefix returns a flattened collection of values
// associated with the given prefix key
func (m *PrefixMap) GetByPrefix(key string) []interface{} {
	mNode := (*Node)(m)
	retrievedNode, _ := mNode.nodeForKey(key, false)
	if retrievedNode == nil {
		return []interface{}{}
	}

	// now, fetching all the values (DFS)
	stack := NewStack()
	values := []interface{}{}
	stack.Push(retrievedNode)
	for stack.Size() > 0 {
		node := stack.Pop().(*Node)
		values = append(values, node.data...)
		for _, c := range node.Children {
			stack.Push(c)
		}
	}

	return values
}

// ContainsPrefix checks if the given prefix is present as key in the map
func (m *PrefixMap) ContainsPrefix(key string) bool {
	mNode := (*Node)(m)
	retrievedNode, _ := mNode.nodeForKey(key, false)
	return retrievedNode != nil
}

// Key Retrieves current node key
// complexity: MAX|O(log(N))| where N
// is the number of nodes in the map.
// Number of nodes in the map cannot exceed
// number of keys + 1.
func (m *Node) Key() string {
	node := m
	k := make([]byte, 0, len(m.key))
	for node != nil && node.isRoot != true {
		key := string(node.key) // triggering a copy here
		k = append([]byte(key), k...)
		node = node.Parent
	}
	return string(k)
}

// PrefixCallback is invoked by EachPrefix for each prefix reached
// by the traversal. The callback has the ability to affect the traversal.
// Returning skipBranch = true will make the traversal skip the current branch
// and jump to the sibling node in the map. Returning halt = true, instead,
// will halt the traversal altogether.
type PrefixCallback func(prefix Prefix) (skipBranch bool, halt bool)

// Prefix holds prefix information
// passed to the PrefixCallback instance by
// the EachPrefifx method.
type Prefix struct {
	node *Node

	// The current prefix string
	Key string

	// The values associated to the current prefix
	Values []interface{}
}

// Depth returns the depth of the corresponding
// node for this prefix in the map.
func (p *Prefix) Depth() int {
	return p.node.Depth()
}

// EachPrefix iterates over the prefixes contained in the
// map using a DFS algorithm. The callback can be used to skip
// a prefix branch altogether or halt the iteration.
func (m *PrefixMap) EachPrefix(callback PrefixCallback) {
	mNode := (*Node)(m)
	stack := NewStack()
	prefix := []byte{}

	skipsubtree := false
	halt := false
	addedLengths := NewStack()
	lastDepth := mNode.Depth()

	stack.Push(mNode)
	for stack.Size() != 0 {
		node := stack.Pop().(*Node)
		if !node.isRoot {
			// if we are now going up
			// in the radix (e.g. we have
			// finished with the current branch)
			// then we adjust the current prefix
			currentDepth := node.Depth()
			if lastDepth >= node.Depth() {
				var length = 0
				for i := 0; i < (lastDepth-currentDepth)+1; i++ {
					length += addedLengths.Pop().(int)
				}
				prefix = prefix[:len(prefix)-length]
			}
			lastDepth = currentDepth
			prefix = append(prefix, node.key...)
			addedLengths.Push(len(node.key))

			// building the info
			// data to pass to the callback
			info := Prefix{
				node:   node,
				Key:    string(prefix),
				Values: node.data,
			}

			skipsubtree, halt = callback(info)
			if halt {
				return
			}
			if skipsubtree {
				continue
			}
		}
		for i := 0; i < len(node.Children); i++ {
			stack.Push(node.Children[i])
		}
	}
}

// -------- auxiliary functions -------- //

// LCP: Longest Common Prefix
// Implementation freely inspired from:
// https://rosettacode.org/wiki/Longest_common_prefix#Go
//
// returns the lcp and the index of the last
// character matching
//
func lcpIndex(strs ...string) int {
	if len(strs) < 2 {
		return -1
	}
	// Special cases first
	switch len(strs) {
	case 0:
		return -1
	case 1:
		return 0
	}
	// LCP of min and max (lexigraphically)
	// is the LCP of the whole set.
	min, max := strs[0], strs[0]
	part := strs[1:]
	for i := 0; i < len(part); i++ {
		s := part[i]
		switch {
		case len(s) < len(min):
			min = s
		case len(s) >= len(max):
			max = s
		}
	}
	for i := 0; i < len(min) && i < len(max); i++ {
		if min[i] != max[i] {
			return i - 1
		}
	}
	// In the case where lengths are not equal but all bytes
	// are equal, min is the answer ("foo" < "foobar").
	return len(min) - 1
}

// Stack ...
type Stack struct {
	size             int
	currentPage      []interface{}
	pages            [][]interface{}
	offset           int
	capacity         int
	pageSize         int
	currentPageIndex int
}

const sDefaultAllocPageSize = 4096

// NewStack Creates a new Stack object with
// an underlying default block allocation size.
// The current default allocation size is one page.
// If you want to use a different block size use
//  NewStackWithCapacity()
func NewStack() *Stack {
	stack := new(Stack)
	stack.currentPage = make([]interface{}, sDefaultAllocPageSize)
	stack.pages = [][]interface{}{stack.currentPage}
	stack.offset = 0
	stack.capacity = sDefaultAllocPageSize
	stack.pageSize = sDefaultAllocPageSize
	stack.size = 0
	stack.currentPageIndex = 0

	return stack
}

// NewStackWithCapacity makes it easy to specify
// a custom block size for inner slice backing the
// stack
func NewStackWithCapacity(cap int) *Stack {
	stack := new(Stack)
	stack.currentPage = make([]interface{}, cap)
	stack.pages = [][]interface{}{stack.currentPage}
	stack.offset = 0
	stack.capacity = cap
	stack.pageSize = cap
	stack.size = 0
	stack.currentPageIndex = 0

	return stack
}

// Push pushes a new element to the stack
func (s *Stack) Push(elem ...interface{}) {
	if elem == nil || len(elem) == 0 {
		return
	}

	// ensures enough pages are ready to be
	// filled in and then fills them from
	// the provided elem array
	if s.size+len(elem) > s.capacity {
		newPages := len(elem) / s.pageSize
		if len(elem)%s.pageSize != 0 {
			newPages++
		}

		// appending new empty pages
		for newPages > 0 {
			page := make([]interface{}, s.pageSize)
			s.pages = append(s.pages, page)
			s.capacity += len(page)
			newPages--
		}
	}

	// now that we have enough pages
	// we can start copying the elements
	// into the stack
	s.size += len(elem)
	for len(elem) > 0 {
		available := len(s.currentPage) - s.offset
		min := len(elem)
		if available < min {
			min = available
		}
		copy(s.currentPage[s.offset:], elem[:min])
		elem = elem[min:]
		s.offset += min
		if len(elem) > 0 {
			// page fully filled; move to the next one
			s.currentPage = s.pages[s.currentPageIndex+1]
			s.currentPageIndex++
			s.offset = 0
		}
	}
}

// Pop pops the top element from the stack
// If the stack is empty it returns nil
func (s *Stack) Pop() (elem interface{}) {
	if s.size == 0 {
		return nil
	}

	s.offset--
	s.size--

	if s.offset < 0 {
		s.offset = s.pageSize - 1

		s.currentPage, s.pages = s.pages[s.currentPageIndex-1], s.pages[:s.currentPageIndex+1]
		s.capacity -= s.pageSize
		s.currentPageIndex--
	}

	elem = s.currentPage[s.offset]

	return
}

// Top ...
func (s *Stack) Top() (elem interface{}) {
	if s.size == 0 {
		return nil
	}

	off := s.offset - 1
	if off < 0 {
		page := s.pages[len(s.pages)-1]
		elem = page[len(page)-1]
		return
	}
	elem = s.currentPage[off]
	return
}

// Size The current size of the stack
func (s *Stack) Size() int {
	return s.size
}
