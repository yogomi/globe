package point

import (
	"github.com/go-gl/mathgl/mgl32"
)

func LongestDistance(vectors []mgl32.Vec3) float32 {
	var max_distance float32 = 0.0
	for i, v1 := range vectors {
		for _, v2 := range vectors[i + 1:] {
			distance := v1.Sub(v2).Len()
			if max_distance < distance {
				max_distance = distance
			}
		}
	}
	return max_distance
}
