package hexagonshell

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func HexagonVertex(size float32) []mgl32.Vec3 {
	var hexagon_vertex = make([]mgl32.Vec3, 0)
	v := mgl32.Vec3{0.0, size, 0.0}
	hexagon_vertex = append(hexagon_vertex, v)

	h := mgl32.QuatRotate(math.Pi / 3, mgl32.Vec3{0.0, 0.0, 1.0})

	var f = func() {
		v = h.Rotate(v)
		hexagon_vertex = append(hexagon_vertex, v)
	}

	f()
	f()
	f()
	f()
	f()

	return hexagon_vertex
}
