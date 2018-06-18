b := a[:0]
for _, x := range a {
	if f(x) {
		b = append(b, x)
	}
}
