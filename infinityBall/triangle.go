package main

import (
	"math"
	"github.com/go-gl/mathgl/mgl64"
)

type Triangle [3]mgl64.Vec3

// make sure that this is right-oriented
func (tri Triangle) Normal() mgl64.Vec3 {
	n := tri[1].Sub(tri[0]).Cross(tri[2].Sub(tri[0]))
	return n.Mul(1/n.Len())
}

func (tri Triangle) Distance(pos mgl64.Vec3) float64 {
	return ScalarProjection(pos, tri.Normal())
}

func ScalarProjection(a, b mgl64.Vec3) float64 {
	return a.Dot(b)/b.Len();
}

func VectorProjection(a, b mgl64.Vec3) mgl64.Vec3 {
	return b.Mul(a.Dot(b)/math.Pow(b.Len(),2))
}

