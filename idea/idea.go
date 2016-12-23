package idea

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"../point"
)

type Idea interface {
	Initialize(vertexes []mgl32.Vec3, program uint32)
	Transport(v mgl32.Vec3)
	Rotate(start, end mgl32.Vec3, radian float32)
	Copy() Idea
	Draw() error

	SetPrimitiveType(primitive_type uint32)
	AddVertexes(vertexes []mgl32.Vec3)
	MergeIdea(o Idea)

	Program() uint32
	Vertexes() []mgl32.Vec3

	// vertexesをOpenGLのBufferへ登録し直す
	rebindBuffer()
	// vertexesの数が変わったらこっちを使わないとダメ。
	replaceBuffer()

	vertexesFirsts() []int32
	vertexesCounts() []int32
}

func NewIdea() Idea {
	object := &idea{0, 0, []mgl32.Vec3{}, gl.LINE_STRIP, []int32{0}, []int32{0}, 0}
	return object
}

type idea struct {
	vao uint32
	vbo uint32
	vertexes []mgl32.Vec3
	primitive_type uint32
	vertexes_firsts []int32
	vertexes_counts []int32
	program uint32
}

func (idea_itself *idea) Initialize(vertexes []mgl32.Vec3, program uint32) {
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

	idea_itself.vertexes_counts = []int32{int32(len(vertexes))}
}

func (idea_itself *idea) Transport(vec mgl32.Vec3) {
	for i, v := range idea_itself.vertexes {
		idea_itself.vertexes[i] = v.Add(vec)
	}

	idea_itself.rebindBuffer()
}

func (idea_itself *idea) Rotate(start, end mgl32.Vec3, radian float32) {
	axis := end.Sub(start).Normalize()
	roter := mgl32.QuatRotate(radian, axis)

	for i, v := range idea_itself.vertexes {
		idea_itself.vertexes[i] = roter.Rotate(v.Sub(start)).Add(start)
	}

	idea_itself.rebindBuffer()
}

func (idea_itself *idea) Copy() Idea {
	object := &idea{0, 0, nil, 0, nil, nil, 0}
	vertexes := make([]mgl32.Vec3, len(idea_itself.vertexes))
	copy(vertexes, idea_itself.vertexes)
	object.Initialize(vertexes, idea_itself.program)

	object.primitive_type = idea_itself.primitive_type
	object.vertexes_firsts = make([]int32, len(idea_itself.vertexes_firsts))
	copy(object.vertexes_firsts, idea_itself.vertexes_firsts)
	object.vertexes_counts = make([]int32, len(idea_itself.vertexes_counts))
	copy(object.vertexes_counts, idea_itself.vertexes_counts)

	return object
}

func (idea_itself *idea) Draw() error {
	gl.BindVertexArray(idea_itself.vao)

	gl.LineWidth(1.5)
	gl.PointSize(5)

	if len(idea_itself.vertexes_firsts) == 1 {
		gl.DrawArrays(idea_itself.primitive_type, 0, int32(len(idea_itself.vertexes)))
	} else {
		gl.MultiDrawArrays(idea_itself.primitive_type,
								(*int32)(gl.Ptr(idea_itself.vertexes_firsts)),
								(*int32)(gl.Ptr(idea_itself.vertexes_counts)),
								int32(len(idea_itself.vertexes_firsts) + 1))
	}

	return nil
}

func (idea_itself *idea) SetPrimitiveType(primitive_type uint32) {
	idea_itself.primitive_type = primitive_type
}

func (idea_itself *idea) AddVertexes(vertexes []mgl32.Vec3) {
	idea_itself.vertexes = append(idea_itself.vertexes, vertexes...)

	last_address := idea_itself.vertexes_firsts[len(idea_itself.vertexes_firsts) - 1]
	idea_itself.vertexes_firsts = append(idea_itself.vertexes_firsts,
					last_address + int32(len(vertexes)))

	idea_itself.vertexes_counts = append(idea_itself.vertexes_counts, int32(len(vertexes)))

	idea_itself.replaceBuffer()
}

// 渡されたIdeaとvertexesとかfirsts,countsをくっつける。
func (idea_itself *idea) MergeIdea(o Idea) {
	idea_itself.vertexes = append(idea_itself.Vertexes(), o.Vertexes()...)

	last_address := idea_itself.vertexes_firsts[len(idea_itself.vertexes_firsts) - 1]
	last_count := idea_itself.vertexes_counts[len(idea_itself.vertexes_counts) - 1]
	firsts := o.vertexesFirsts()
	for i, v := range firsts {
		firsts[i] = v + last_address + last_count
	}

	idea_itself.vertexes_firsts = append(idea_itself.vertexesFirsts(), firsts...)
	idea_itself.vertexes_counts = append(idea_itself.vertexesCounts(), o.vertexesCounts()...)

	idea_itself.replaceBuffer()
}

func (idea_itself *idea) Program() uint32 {
	return idea_itself.program
}

// ideaに登録されているvertexesのバッファを更新する関数
// privateなため、スライスはアドレスのコピーで良いものとする。
func (idea_itself *idea) rebindBuffer() {
	gl.BindVertexArray(idea_itself.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, idea_itself.vbo)

	float_array := point.VectorToArrayOfFloat32(idea_itself.vertexes)
	gl.BufferSubData(gl.ARRAY_BUFFER,
						0,
						len(float_array) * 4,
						gl.Ptr(float_array))
}

func (idea_itself *idea) replaceBuffer() {
	gl.BindVertexArray(idea_itself.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, idea_itself.vbo)

	float_array := point.VectorToArrayOfFloat32(idea_itself.vertexes)
	gl.BufferData(gl.ARRAY_BUFFER,
						len(float_array) * 4,
						gl.Ptr(float_array),
						gl.DYNAMIC_DRAW)
}
func (idea_itself *idea) Vertexes() []mgl32.Vec3 {
	return idea_itself.vertexes
}

func (idea_itself *idea) vertexesFirsts() []int32 {
	return idea_itself.vertexes_firsts
}

func (idea_itself *idea) vertexesCounts() []int32 {
	return idea_itself.vertexes_counts
}
