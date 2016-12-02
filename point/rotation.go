package point

import (
	"github.com/go-gl/mathgl/mgl32"
)

func RotateAroundVector(start mgl32.Vec3,
					end mgl32.Vec3,
					radian float32,
					vectors []mgl32.Vec3) []mgl32.Vec3 {
	rotated_vectors := make([]mgl32.Vec3, len(vectors))

	axis := end.Sub(start).Normalize()
	roter := mgl32.QuatRotate(radian, axis)

	for i, v := range vectors {
		rotated_vectors[i] = roter.Rotate(v.Sub(start)).Add(start)
	}
	return rotated_vectors
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
