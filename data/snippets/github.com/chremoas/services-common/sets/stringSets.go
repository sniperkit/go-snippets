package sets

type StringSet struct {
	Set map[string]bool
}

func NewStringSet() *StringSet {
	return &StringSet{make(map[string]bool)}
}

func (set *StringSet) Add(s string) bool {
	if len(s) == 0 {
		return false
	}
	_, found := set.Set[s]
	set.Set[s] = true
	return !found
}

func (set *StringSet) Remove(s string) {
	delete(set.Set, s)
}

func (set *StringSet) Contains(s string) bool {
	_, found := set.Set[s]
	return found
}

func (set *StringSet) Len() int {
	return len(set.Set)
}

func (set *StringSet) FromSlice(slice []string) {
	for s := range slice {
		set.Add(slice[s])
	}
}

func (set *StringSet) ToSlice() (slice []string) {
	for s := range set.Set {
		slice = append(slice, s)
	}

	return slice
}

func (set *StringSet) Intersection(set1 *StringSet) *StringSet {
	var output = NewStringSet()

	for s := range set.Set {
		if set1.Contains(s) {
			output.Add(s)
		}
	}

	return output
}

func (set *StringSet) Difference(set1 *StringSet) *StringSet {
	var output = NewStringSet()

	for s := range set.Set {
		if !set1.Contains(s) {
			output.Add(s)
		}
	}

	return output
}
