package fixtures

func maps() {
	var m map[string]int
	m = make(map[string]int)
	m = make(map[string]int, 3)
	m = map[string]int{
		"a": 1, x: 2,
	}

	v := m["a"]
	m["foo"] = v + len(m)

	v, ok := m[x]
	_, _ = v, ok

	for k := range m {
		v := m[k]
		_, _ = k, v
	}

	for k, v := range m {
		_, _ = k, v
	}
}
