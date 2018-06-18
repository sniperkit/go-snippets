people := []struct {
	Name string
	Age  int
}{
	{"Gopher", 7},
	{"Alice", 55},
	{"Vera", 24},
	{"Bob", 75},
}

// sort people by age
sort.Slice(people, func(i, j int) bool {
	return people[i].Age < people[j].Age
})
