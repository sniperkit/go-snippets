package magicsquare

func Permute(elems ...int) [][]int {
	if len(elems) == 0 {
		return [][]int{}

	}
	ret := make([][]int, 0)
	for _, elem := range elems {
		p := []int{elem}
		perms := Permute(rest(elems, elem)...)
		if len(perms) == 0 {
			ret = append(ret, p)
		} else {
			for _, perm := range perms {
				ret = append(ret, append(p, perm...))
			}
		}
	}
	return ret
}

func rest(nums []int, n int) []int {
	var elems []int
	for _, num := range nums {
		if num != n {
			elems = append(elems, num)
		}
	}
	return elems
}
