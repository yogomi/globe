package point

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func HexagonVertex() []float32 {
	var hexagon_vertex = []float32 {
		0.0, 1.0, 0.0,
	}

	h := mgl32.QuatRotate(math.Pi / 3, mgl32.Vec3{0.0, 0.0, 1.0})
	v := h.Rotate(mgl32.Vec3{0.0, 1.0, 0.0})

	var fv = make([]float32, 3)
	fv[0], fv[1], fv[2] = v.Elem()
	hexagon_vertex = append(hexagon_vertex, fv...)

	var f = func() {
		v = h.Rotate(v)
		fv = make([]float32, 3)
		fv[0], fv[1], fv[2] = v.Elem()
		hexagon_vertex = append(hexagon_vertex, fv...)
	}

	f()
	f()
	f()
	f()
	f()

	return hexagon_vertex
}
