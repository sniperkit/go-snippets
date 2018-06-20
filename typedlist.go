package snippets

import "sort"

type TypeAList []TYPEA

func (x *TypeAList) Push(val TYPEA) {
	*x = append(*x, val)
}

func (x *TypeAList) Pop() TYPEA {
	y := (*x)[len(*x)-1]
	*x = (*x)[0 : len(*x)-1]
	return y
}

func (x *TypeAList) Insert(ind int, val TYPEA) {
	tmp := (*x)[ind:]
	*x = (*x)[0:ind]
	*x = append(*x, val)
	*x = append(*x, tmp...)
}

func (x *TypeAList) Remove(ind int) TYPEA {
	y := (*x)[ind]
	tmp := (*x)[ind+1:]
	*x = (*x)[0:ind]
	*x = append(*x, tmp...)
	return y
}

func (x *TypeAList) Index(val TYPEA) int {
	for ind, cand := range *x {
		if val == cand {
			return ind
		}
	}
	return -1
}

func (x TypeAList) Filter(fn func(TYPEA) bool) TypeAList {
	y := TypeAList{}
	for _, val := range x {
		if fn(val) {
			y = append(y, val)
		}
	}
	return y
}

// Sorting support, given a Less function

type TYPEASort struct {
	list TypeAList
	less func(a, b TYPEA) bool
}

func (x TYPEASort) Len() int {
	return len(x.list)
}

func (x TYPEASort) Swap(i, j int) {
	x.list[i], x.list[j] = x.list[j], x.list[i]
}

func (x TYPEASort) Less(i, j int) bool {
	return x.less(x.list[i], x.list[j])
}

func (x *TypeAList) Sort(less func(a, b TYPEA) bool) {
	ts := TYPEASort{*x, less}
	sort.Sort(ts)
}
