package list

import "testing"

type templary []int

var data templary = []int{1, 2, 3, 4, 5}

func (t templary) Value() int {
	return t[0]
}

func (t templary) Next() List {
	return t[1:]
}

func TestLazy1(t *testing.T) {
	var l List
	l = &lazy{
		f: func() (int, List) {
			v := data.Value()
			data, _ = data.Next().(templary)
			return v, data
		},
	}

	for i := 1; i <= 5; i++ {
		v := l.Value()
		l = l.Next()
		if v != i {
			t.Fatalf("lazy %dth not equal %d, got %d", i, i, v)
		}
	}
}
