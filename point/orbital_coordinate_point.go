package point

import (
	"math"
)

func GlobeOrbitalCoordinatePoints(step int) []float32 {
	const DIV = 0.0001
	const RADIUS = 1.0
	const TWIST_STRENGTH = 120.0

	max_step := int(RADIUS / DIV * 2 + 1)

	if step > max_step {
		step = max_step
	}

	var globe []float32 = make([]float32, 3 * step)

	for i := 0; i < step; i++ {
		var theta float64 = float64(i) / float64(max_step) * math.Pi
		r := RADIUS * math.Sin(theta)
		globe[3 * i] = float32(r * math.Sin((theta * TWIST_STRENGTH)))
		globe[3 * i + 1] = float32(RADIUS * math.Cos(theta))
		globe[3 * i + 2] = float32(r * math.Cos(theta * TWIST_STRENGTH))
	}

	return globe
}
