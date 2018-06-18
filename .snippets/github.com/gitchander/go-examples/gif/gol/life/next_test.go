package life

import (
	"math/rand"
	"testing"
	"time"
)

func TestNext(t *testing.T) {

	ps := randPoints()

	t.Logf("count: %d", len(ps))

	var (
		ps1 = ps.Clone()
		ps2 = ps.Clone()
		ps3 = ps.Clone()
		ps4 = ps.Clone()
	)

	n1 := Nexter1{}
	n2 := Nexter2{}
	n3 := Nexter3{}
	n4 := Nexter4{}

	for i := 0; i < 100; i++ {

		ps1 = n1.Next(ps1)
		ps2 = n2.Next(ps2)
		ps3 = n3.Next(ps3)
		ps4 = n4.Next(ps4)

		if !ps1.Equal(ps2) {
			t.Fatal("(1 & 2) not equal!")
		}
		if !ps2.Equal(ps3) {
			t.Fatal("(2 & 3) not equal!")
		}
		if !ps3.Equal(ps4) {
			t.Fatal("(3 & 4) not equal!")
		}

		//t.Logf("count: %d", len(ps1))
	}
}

func randPoints() Points {
	var (
		nX = 50
		nY = 50
	)
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	var ps []Point
	for y := 0; y < nY; y++ {
		for x := 0; x < nX; x++ {
			if (r.Int() & 1) == 1 {
				ps = append(ps, Point{x, y})
			}
		}
	}
	return ps
}

var benchPoints = randPoints()

func BenchmarkNext1(b *testing.B) {
	ps := benchPoints.Clone()
	n := Nexter1{}
	for i := 0; i < b.N; i++ {
		ps = n.Next(ps)
	}
}

func BenchmarkNext2(b *testing.B) {
	ps := benchPoints.Clone()
	n := Nexter2{}
	for i := 0; i < b.N; i++ {
		ps = n.Next(ps)
	}
}

func BenchmarkNext3(b *testing.B) {
	ps := benchPoints.Clone()
	n := Nexter3{}
	for i := 0; i < b.N; i++ {
		ps = n.Next(ps)
	}
}

func BenchmarkNext4(b *testing.B) {
	ps := benchPoints.Clone()
	n := Nexter4{}
	for i := 0; i < b.N; i++ {
		ps = n.Next(ps)
	}
}
