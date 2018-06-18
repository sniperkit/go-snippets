package wmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/yizha/go/w2v"
)

var model = &w2v.Model{
	FeatureSize: 2,
	Word2id: map[string]int{
		"test":  0,
		"word":  1,
		"world": 2,
		"text":  3,
		"a":     4,
		"the":   5,
	},
	Vectors: []w2v.Vector{
		[]float64{0.18, 0.24},
		[]float64{0.09, 0.43},
		[]float64{-0.18, 0.23},
		[]float64{0.04, 0.11},
		[]float64{-0.03, 0.02},
		[]float64{-0.07, 0.05},
	},
}

func TestWMD(t *testing.T) {
	t1 := "A test word"
	t2 := "The text world"
	d1 := strings.Split(t1, " ")
	d2 := strings.Split(t2, " ")
	distance, err := Wmd(d1, d2, model)
	if err != nil {
		t.Error("Wmd() returns error:", err)
		return
	}
	fmt.Printf("wmd between [%v] and [%v] is %v\n", t1, t2, distance)
}
