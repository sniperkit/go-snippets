package main

import (
	"unsafe"
	"github.com/go-gl/gl"
	"github.com/go-gl/glu"
	"github.com/go-gl/mathgl/mgl64"
)

type Ball struct {
	Pos mgl64.Vec3
	Velocity mgl64.Vec3
	Radius float64
	Mass float64
	Bounciness float64
	sphere unsafe.Pointer
}

func NewBall(radius, mass, bounciness float64) *Ball {
	ball := new(Ball)
	ball.sphere = glu.NewQuadric()
	ball.Radius = radius
	ball.Mass = mass
	ball.Bounciness = bounciness
	ball.Pos = mgl64.Vec3{0,0,0}
	return ball
}

func (ball *Ball) Draw() {
	gl.PushMatrix()
	gl.Translated(ball.Pos.X(), ball.Pos.Y(), ball.Pos.Z())
	glu.Sphere(ball.sphere, float32(ball.Radius), 10, 10)
	gl.PopMatrix()
}

