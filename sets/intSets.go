package sets

type IntSet struct {
	set map[int]bool
}

func NewIntSet() *IntSet {
	return &IntSet{make(map[int]bool)}
}

func (set *IntSet) Add(s int) bool {
	_, found := set.set[s]
	set.set[s] = true
	return !found
}

func (set *IntSet) Remove(s int) {
	delete(set.set, s)
}

func (set *IntSet) Contains(s int) bool {
	_, found := set.set[s]
	return found
}

func (set *IntSet) Size() int {
	return len(set.set)
}