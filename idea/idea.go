package idea

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"../point"
)

type Idea interface {
	Initialize(vertexes []mgl32.Vec3, program uint32)
	RebindVertexes(vertexes []mgl32.Vec3)
	Vertexes() []mgl32.Vec3
	Transport(v mgl32.Vec3)
	Rotate(start, end mgl32.Vec3, radian float32)
	Copy() Idea
	Draw() error
}

func NewIdea() Idea {
	object := &idea{0, 0, []mgl32.Vec3{}, 0}
	return object
}

type idea struct {
	vao uint32
	vbo uint32
	vertexes []mgl32.Vec3
	program uint32
}

func (idea_itself *idea) RebindVertexes(vertexes []mgl32.Vec3) {
	idea_itself.vertexes = vertexes

	gl.BindVertexArray(idea_itself.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, idea_itself.vbo)

	float_array := point.VectorToArrayOfFloat32(vertexes)
	gl.BufferSubData(gl.ARRAY_BUFFER,
						0,
						len(float_array) * 4,
						gl.Ptr(float_array))
}

func (idea_itself *idea) Vertexes() []mgl32.Vec3 {
	return idea_itself.vertexes
}

func (idea_itself *idea) Initialize(vertexes []mgl32.Vec3, program uint32) {
	fmt.Println("initialize idea")
	idea_itself.program = program
	idea_itself.vertexes = vertexes

	gl.GenVertexArrays(1, &(idea_itself.vao))
	gl.BindVertexArray(idea_itself.vao)
	gl.GenBuffers(1, &(idea_itself.vbo))
	gl.BindBuffer(gl.ARRAY_BUFFER, idea_itself.vbo)

	float_array := point.VectorToArrayOfFloat32(vertexes)
	gl.BufferData(gl.ARRAY_BUFFER,
						len(float_array) * 4,
						gl.Ptr(float_array),
						gl.DYNAMIC_DRAW)

	vert_attrib := uint32(gl.GetAttribLocation(idea_itself.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vert_attrib)
	gl.VertexAttribPointer(vert_attrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
}

func (idea_itself *idea) Transport(vec mgl32.Vec3) {
	for i, v := range idea_itself.vertexes {
		idea_itself.vertexes[i] = v.Add(vec)
	}

	idea_itself.RebindVertexes(idea_itself.vertexes)
}

func (idea_itself *idea) Rotate(start, end mgl32.Vec3, radian float32) {
	axis := end.Sub(start).Normalize()
	roter := mgl32.QuatRotate(radian, axis)

	for i, v := range idea_itself.vertexes {
		idea_itself.vertexes[i] = roter.Rotate(v.Sub(start)).Add(start)
	}

	idea_itself.RebindVertexes(idea_itself.vertexes)
}

func (idea_itself *idea) Copy() Idea {
	object := &idea{0, 0, []mgl32.Vec3{}, 0}
	vertexes := make([]mgl32.Vec3, len(idea_itself.vertexes))
	copy(vertexes, idea_itself.vertexes)
	object.Initialize(vertexes, idea_itself.program)
	return object
}

func (idea_itself *idea) Draw() error {
	gl.BindVertexArray(idea_itself.vao)

	gl.LineWidth(1.5);
	gl.DrawArrays(gl.LINE_STRIP, 0, int32(len(idea_itself.vertexes)))
	return nil
}
