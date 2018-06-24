package main

import (
	"testing"
)

func TestAgeDistribution(t *testing.T) {
	for k, v := range map[int]string{
		-2:  "This program is for humans",
		0:   "Still in Mama's arms",
		3:   "Preschool Maniac",
		5:   "Elementary school",
		12:  "Middle school",
		15:  "High school",
		19:  "College",
		23:  "Working for the man",
		66:  "The Golden Years",
		101: "This program is for humans"} {
		if r := ageDistribution(k); r != v {
			t.Errorf("failed: ageDistribution %d is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkAgeDistribution(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ageDistribution(i % 102)
	}
}

var (
	category []string
	age      []int
)

func init() {
	category = []string{"This program is for humans",
		"Still in Mama's arms",
		"Preschool Maniac",
		"Elementary school",
		"Middle school",
		"High school",
		"College",
		"Working for the man",
		"The Golden Years"}
	age = []int{0, 3, 5, 12, 15, 19, 23, 66, 101}
}

func ageDistribution(m int) string {
	var c int
	for c < len(age) && m >= age[c] {
		c++
	}
	return category[c%len(age)]
}
