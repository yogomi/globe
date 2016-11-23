package point

import (
	"github.com/go-gl/mathgl/mgl32"
)

func RotateAroundVector(v_star mgl32.Vec3, end mgl32.Vec3, vectors []mgl32.Vec3) []mgl32.Vec3 {
	return vectors
}

func VectorToArrayOfFloat32(array_of_vector []mgl32.Vec3) []float32 {
	var array_of_float32 = make([]float32, 0)
	for _, v := range(array_of_vector) {
		var fv = make([]float32, 3)
		fv[0], fv[1], fv[2] = v.Elem()
		array_of_float32 = append(array_of_float32, fv...)
	}
	return array_of_float32
}
