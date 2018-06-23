package list

type packed []int

func (p packed) Value() int {
	return p[0]
}

func (p packed) Next() List {
	if len(p) == 0 {
		return nil
	}
	return p[1:]
}

func Pack(l List) List {
	if l == nil {
		return nil
	}
	var p packed
	for ; l != nil; l = l.Next() {
		p = append(p, l.Value())
	}
	return p
}
