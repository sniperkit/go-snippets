package sorting

func Sort(list []int) {
	n := len(list)
	for i := n / 2; i >= 1; i-- {
		sink(i, n, list)
	}
	for n > 1 {
		swap(1,n,list)
		n--
		sink(1, n, list)
	}
}

func swap(i, n int, list []int) {
	list[i-1], list[n-1] = list[n-1], list[i-1]
}

func sink(p, n int, list []int) {
	for 2*p <= n {
		j := p * 2
		if j < n && list[j-1] < list[j] {
			j++
		}
		if list[p-1] > list[j-1] {
			break
		}
		swap(p,j,list)
		p = j
	}
}
