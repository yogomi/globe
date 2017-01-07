package triangleporyhedron

import (
	"fmt"
	"math"
	"time"
	"errors"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"../idea"
)

type TrianglePoryhedron interface {
	Transport(v mgl32.Vec3)
	Rotate(start, end mgl32.Vec3, radian float32)
	Copy() TrianglePoryhedron
	Draw() error
	Vertexes() []mgl32.Vec3

	Segmentalize() error
}

type triangle_poryhedron struct {
	core idea.Idea
	radius float32
	center mgl32.Vec3
}

func NewTrianglePoryhedron(radius float64, program uint32) TrianglePoryhedron {
	poryhedron := &triangle_poryhedron{idea.NewIdea(),
					float32(radius),
					mgl32.Vec3{0.0, 0.0, 0.0}}

	icosahedron_vertexes := icosahedronVertex(radius)

	poryhedron.core.Initialize(icosahedron_vertexes[0: 3], program)
	poryhedron.core.SetPrimitiveType(gl.LINE_LOOP)

	for i := 1; i < 20; i++ {
		poryhedron.core.AddVertexes(icosahedron_vertexes[i * 3: (i + 1) * 3])
	}
	return poryhedron
}

func (poryhedron *triangle_poryhedron) Transport(v mgl32.Vec3) {
	poryhedron.core.Transport(v)
	poryhedron.center = poryhedron.center.Add(v)
}

func (poryhedron *triangle_poryhedron) Rotate(start, end mgl32.Vec3, radian float32) {
	poryhedron.core.Rotate(start, end, radian)

	axis := end.Sub(start).Normalize()
	roter := mgl32.QuatRotate(radian, axis)
	poryhedron.center = roter.Rotate(poryhedron.center.Sub(start)).Add(start)
}

func (poryhedron *triangle_poryhedron) Copy() TrianglePoryhedron {
	p := &triangle_poryhedron{poryhedron.core.Copy(),
					poryhedron.radius,
					poryhedron.center}
	return p
}

func (poryhedron *triangle_poryhedron) Draw() error {
	return poryhedron.core.Draw()
}

func (poryhedron *triangle_poryhedron) Vertexes() []mgl32.Vec3 {
	return poryhedron.core.Vertexes()
}

func (poryhedron *triangle_poryhedron) Segmentalize() error {
	vertexes := poryhedron.Vertexes()

	split_triangle := func(a []mgl32.Vec3) ([]mgl32.Vec3, error) {
		if len(a) != 3 {
			return nil, errors.New("vectors length it not 3")
		}

		b := make([]mgl32.Vec3, 3)
		b[0] = a[0].Add(a[1]).Sub(poryhedron.center).Normalize().Mul(poryhedron.radius).Add(poryhedron.center)
		b[1] = a[1].Add(a[2]).Sub(poryhedron.center).Normalize().Mul(poryhedron.radius).Add(poryhedron.center)
		b[2] = a[2].Add(a[0]).Sub(poryhedron.center).Normalize().Mul(poryhedron.radius).Add(poryhedron.center)

		result_vectors := make([]mgl32.Vec3, 3 * 4)

		f := func(i int, v1, v2, v3 mgl32.Vec3) {
			result_vectors[i * 3] = v1
			result_vectors[i * 3 + 1] = v2
			result_vectors[i * 3 + 2] = v3
		}
		f(0, b[0], b[1], b[2])
		f(1, a[0], b[0], b[2])
		f(2, a[1], b[1], b[0])
		f(3, a[2], b[2], b[1])
		return result_vectors, nil
	}

	now := time.Now()

	initialized := false
	new_core := idea.NewIdea()

	for i := 0; i < (len(vertexes) / 3); i++ {
		new_vertexes, err := split_triangle(vertexes[3 * i: 3 * i + 3])
		if err != nil {
			return err
		}
		for j := 0; j < len(new_vertexes) / 3; j++ {
			if initialized {
				new_core.AddVertexesWithoutReplace(new_vertexes[j * 3: (j + 1) * 3])
			} else {
				fmt.Println("side length =", new_vertexes[j * 3].Sub(new_vertexes[j * 3 + 1]).Len())
				new_core.Initialize(new_vertexes[j * 3: (j + 1) * 3],
								poryhedron.core.Program())
				new_core.SetPrimitiveType(gl.LINE_LOOP)
				initialized = true
			}
		}
	}
	new_core.ReplaceBuffer()

	fmt.Println("end  ", time.Now().Sub(now))
	poryhedron.core = new_core
	return nil
}

func icosahedronVertex(radius float64) []mgl32.Vec3 {
	var a, g float64
	a = radius * math.Sqrt(2.0) / math.Sqrt(5.0 + math.Sqrt(5.0))
	g = ((1.0 + math.Sqrt(5.0)) / 2.0) * a

	alpha := []mgl32.Vec3{{float32(a), float32(g), 0},
					{float32(-a), float32(g), 0},
					{float32(-a), float32(-g), 0},
					{float32(a), float32(-g), 0}}
	beta := []mgl32.Vec3{{0, float32(a), float32(g)},
					{0, float32(-a), float32(g)},
					{0, float32(-a), float32(-g)},
					{0, float32(a), float32(-g)}}
	ganma := []mgl32.Vec3{{float32(g), 0, float32(a)},
					{float32(g), 0, float32(-a)},
					{float32(-g), 0, float32(-a)},
					{float32(-g), 0, float32(a)}}

	icosahedron := make([]mgl32.Vec3, 3 * 20)

	f := func(i int, v1, v2, v3 mgl32.Vec3) {
		icosahedron[i * 3] = v1
		icosahedron[i * 3 + 1] = v2
		icosahedron[i * 3 + 2] = v3
	}

	f(0, alpha[0], alpha[1], beta[0])
	f(1, beta[0], beta[1], ganma[0])
	f(2, ganma[0], ganma[1], alpha[0])
	f(3, alpha[1], alpha[0], beta[3])
	f(4, beta[1], beta[0], ganma[3])
	f(5, ganma[1], ganma[0], alpha[3])
	f(6, alpha[2], alpha[3], beta[1])
	f(7, beta[2], beta[3], ganma[1])
	f(8, ganma[2], ganma[3], alpha[1])
	f(9, alpha[3], alpha[2], beta[2])
	f(10, beta[3], beta[2], ganma[2])
	f(11, ganma[3], ganma[2], alpha[2])
	f(12, alpha[0], beta[0], ganma[0])
	f(13, alpha[0], ganma[1], beta[3])
	f(14, alpha[1], beta[3], ganma[2])
	f(15, alpha[1], ganma[3], beta[0])
	f(16, alpha[2], beta[1], ganma[3])
	f(17, alpha[2], ganma[2], beta[2])
	f(18, alpha[3], ganma[0], beta[1])
	f(19, alpha[3], ganma[1], beta[2])

	return icosahedron
}
