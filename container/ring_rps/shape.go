package main

import (
	"container/ring"
	"math/rand"
	"strconv"
	"time"
)

type Shape int

const (
	Rock Shape = iota
	Paper
	Scissors
)

var name_Shape = map[Shape]string{
	Rock:     "Rock",
	Paper:    "Paper",
	Scissors: "Scissors",
}

func (e Shape) String() (s string) {
	if name, ok := name_Shape[e]; ok {
		return name
	}
	return strconv.Itoa(int(e))
}

func (a Shape) Beats(b Shape) bool {
	if (a == Rock) && (b == Scissors) {
		return true
	}
	if (a == Paper) && (b == Rock) {
		return true
	}
	if (a == Scissors) && (b == Paper) {
		return true
	}
	return false
}

var nextShape = func() func() Shape {

	shapes := make(chan Shape)

	go func() {
		vs := []Shape{Rock, Paper, Scissors}
		r := ring.New(len(vs))
		for _, shape := range vs {
			r.Value = shape
			r = r.Next()
		}
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		for {
			n := random.Intn(10 * len(vs))
			for i := 0; i < n; i++ {
				r = r.Next()
			}
			shapes <- r.Value.(Shape)
		}
	}()

	return func() Shape {
		return <-shapes
	}
}()
