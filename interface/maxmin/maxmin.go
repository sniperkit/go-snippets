package maxmin

type Interface interface {
	Len() int
	Less(i, j int) bool
}

func IndexOfMax(v Interface) (max int) {
	n := v.Len()
	if n == 0 {
		return -1
	}
	for i := 1; i < n; i++ {
		if v.Less(max, i) {
			max = i
		}
	}
	return max
}

func IndexOfMin(v Interface) (min int) {
	n := v.Len()
	if n == 0 {
		return -1
	}
	for i := 1; i < n; i++ {
		if v.Less(i, min) {
			min = i
		}
	}
	return min
}

func IndexOfMaxMin(v Interface) (max, min int) {
	n := v.Len()
	if n == 0 {
		return -1, -1
	}
	for i := 1; i < n; i++ {
		if v.Less(max, i) {
			max = i
		}
		if v.Less(i, min) {
			min = i
		}
	}
	return max, min
}

type IntSlice []int

func (v IntSlice) Len() int           { return len(v) }
func (v IntSlice) Less(i, j int) bool { return v[i] < v[j] }
