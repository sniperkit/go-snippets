package list

type fib [2]int

func (l fib) Value() int {
	return l[0]
}

func (l fib) Next() List {
	return fib{l[1], l[0] + l[1]}
}
