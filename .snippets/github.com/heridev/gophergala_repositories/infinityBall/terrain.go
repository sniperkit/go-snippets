package main

import (
	"fmt"
	"math"
	"github.com/go-gl/gl"
	"github.com/go-gl/mathgl/mgl64"
)

type Terrain struct {
	Scale mgl64.Vec3
	Pos mgl64.Vec3
	Heights [][]float64
	Verts [][][3]float64
	Norms [][][3]float64
	DrawAsSurface bool
}

func newTerrain(scale mgl64.Vec3, heights [][]float64) *Terrain {
	t := new(Terrain)
	t.Scale = scale
	t.Heights = heights
	t.Verts = calculateVertices(scale, heights)
	t.Norms = calculateNormals(t.Verts)
	t.DrawAsSurface = true
	return t
}

func CreateTerrain(scale mgl64.Vec3, size int) *Terrain {
	return newTerrain(scale, createHM(size))
}

func ReadTerrain(scale mgl64.Vec3) *Terrain {
	return newTerrain(scale, readHeightmap())
}

func readHeightmap() [][]float64 {
	var n int
	fmt.Scan(&n)
	hm := make([][]float64, n)
	for y:=0; y<n; y++ {
		hm[y] = make([]float64, n)
		for x:=0; x<n; x++ {
			fmt.Scan(&hm[y][x])
		}
	}
	return hm
}

func calculateVertices(scale mgl64.Vec3, heights [][]float64) [][][3]float64 {
	h := len(heights)
	w := len(heights[0])
	verts := make([][][3]float64, h)
	for y:=0; y<h; y++ {
		verts[y] = make([][3]float64, w)
		for x:=0; x<w; x++ {
			verts[y][x] = [3]float64{float64(x)*scale[0], heights[y][x]*scale[1], float64(y)*scale[2]}
		}
	}
	return verts
}

func calculateNormals(verts [][][3]float64) [][][3]float64 {
	h := len(verts)
	w := len(verts[0])

	normals := make([][][3]float64, h)
	for y:=0; y<h; y++ {
		normals[y] = make([][3]float64, w)
		for x:=0; x<w; x++ {
			normals[y][x] = [3]float64{0,1,0}
		}
	}
	
	for y:=1; y<h-1; y++ {
		for x:=1; x<w-1; x++ {
			c := verts[y  ][x  ]
			
			l := verts[y  ][x-1]
			d := verts[y-1][x  ]
			r := verts[y  ][x+1]
			u := verts[y+1][x  ]

			n1 := Triangle{c,u,r}.Normal()
			n2 := Triangle{c,l,u}.Normal()
			n3 := Triangle{c,d,l}.Normal()
			n4 := Triangle{c,r,d}.Normal()
		
			normals[y][x] = n1.Add(n2).Add(n3).Add(n4).Mul(.25)
		}
	}
	return normals
}

func (t *Terrain) GetTriangleUnder(pos mgl64.Vec3) Triangle {
	h := len(t.Heights)
	w := len(t.Heights[0])

	v := pos.Sub(t.Pos)
	xf,xr := math.Modf(v.X()/t.Scale.X())
	yf,yr := math.Modf(v.Z()/t.Scale.Z())
	x := int(xf)
	y := int(yf)
	
	if x<0 || y<0 || x>=w || y>=h {
		return Triangle{mgl64.Vec3{math.NaN(),0,0},mgl64.Vec3{0,0,0},mgl64.Vec3{0,0,0}}
	}
	
	if xr < yr {
		return Triangle{t.Verts[y][x], t.Verts[y+1][x], t.Verts[y+1][x+1]}
	} else {
		return Triangle{t.Verts[y+1][x+1], t.Verts[y][x+1], t.Verts[y][x]}
	}
}

func (t *Terrain) Draw() {
	h := len(t.Heights)
	w := len(t.Heights[0])
	
	gl.PushMatrix()

	gl.Translated(t.Pos[0], t.Pos[1], t.Pos[2])
	
	if t.DrawAsSurface {
		gl.Begin(gl.TRIANGLES)
		for y:=0; y<h-1; y++ {
			for x:=0; x<w-1; x++ {
				gl.Vertex3dv(&t.Verts[y  ][x  ])
				gl.Normal3dv(&t.Norms[y  ][x  ])
				gl.Vertex3dv(&t.Verts[y  ][x+1])
				gl.Normal3dv(&t.Norms[y  ][x+1])
				gl.Vertex3dv(&t.Verts[y+1][x+1])
				gl.Normal3dv(&t.Norms[y+1][x+1])
				gl.Vertex3dv(&t.Verts[y+1][x+1])
				gl.Normal3dv(&t.Norms[y+1][x+1])
				gl.Vertex3dv(&t.Verts[y+1][x  ])
				gl.Normal3dv(&t.Norms[y+1][x  ])
				gl.Vertex3dv(&t.Verts[y  ][x  ])
				gl.Normal3dv(&t.Norms[y  ][x  ])
			}
		}
	} else {
		gl.Begin(gl.LINES)
		for y:=0; y<h-1; y++ {
			for x:=0; x<w-1; x++ {
				gl.Vertex3dv(&t.Verts[y  ][x  ])
				gl.Vertex3dv(&t.Verts[y  ][x+1])
				gl.Vertex3dv(&t.Verts[y  ][x  ])
				gl.Vertex3dv(&t.Verts[y+1][x  ])
				gl.Vertex3dv(&t.Verts[y  ][x  ])
				gl.Vertex3dv(&t.Verts[y+1][x+1])
				gl.Vertex3dv(&t.Verts[y  ][x+1])
				gl.Vertex3dv(&t.Verts[y+1][x+1])
				gl.Vertex3dv(&t.Verts[y+1][x  ])
				gl.Vertex3dv(&t.Verts[y+1][x+1])
			}
		}
	}

	gl.End()
	
	gl.PopMatrix()
}

