package idea

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"../point"
)

type Idea interface {
	Initialize(program uint32)
	BindVertexes(vertexes []mgl32.Vec3)
	Draw() error
}

func NewIdea() Idea {
	object := &idea{}
	return object
}

type idea struct {
	vao uint32
	vbo uint32
	vertexes []mgl32.Vec3
	program uint32
}

func (idea_itself *idea) BindVertexes(vertexes []mgl32.Vec3) {
	fmt.Println(vertexes)
	gl.BindVertexArray(idea_itself.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, idea_itself.vbo)

	float_array := point.VectorToArrayOfFloat32(vertexes)
	gl.BufferData(gl.ARRAY_BUFFER,
						len(float_array) * 4,
						gl.Ptr(float_array),
						gl.DYNAMIC_DRAW)

	idea_itself.vertexes = vertexes
}

func (idea_itself *idea) Initialize(program uint32) {
	fmt.Println("initialize idea")
	idea_itself.program = program
	fmt.Println(idea_itself.program)
	gl.GenVertexArrays(1, &(idea_itself.vao))
	gl.BindVertexArray(idea_itself.vao)
	gl.GenBuffers(1, &(idea_itself.vbo))
	gl.BindBuffer(gl.ARRAY_BUFFER, idea_itself.vbo)

	vert_attrib := uint32(gl.GetAttribLocation(idea_itself.program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vert_attrib)
	gl.VertexAttribPointer(vert_attrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
}

func (idea_itself *idea) Draw() error {
	gl.BindVertexArray(idea_itself.vao)

	gl.LineWidth(1.5);
	gl.DrawArrays(gl.LINE_STRIP, 0, int32(len(idea_itself.vertexes)))
	return nil
}
