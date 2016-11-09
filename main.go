package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"./misc"
	"./shader"
	"./point"
)

const window_width = 800
const window_height = 600

var angle float64 = 0.0
var rotate_span float64 = 0.0

func main() {
	fmt.Println("start")

	misc.CheckError(glfw.Init())
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	window, err := glfw.CreateWindow(window_width,
										window_height,
										"Globe",
										nil,
										nil)
	misc.CheckError(err)
	window.SetKeyCallback(keyCallback)
	window.MakeContextCurrent()

	misc.CheckError(gl.Init())

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	program, err := shader.CreateProgram()
	misc.CheckError(err)

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0),
									float32(window_width)/window_height,
									0.1,
									10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{0, 0, 5}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// Configure the vertex data
	var vao_hexagon uint32
	gl.GenVertexArrays(1, &vao_hexagon)
	gl.BindVertexArray(vao_hexagon)

	var vbo_hexagon uint32
	gl.GenBuffers(1, &vbo_hexagon)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_hexagon)

	var hexagon_vertex = point.HexagonVertex()

	gl.BufferData(gl.ARRAY_BUFFER,
						len(hexagon_vertex) * 4,
						gl.Ptr(hexagon_vertex),
						gl.DYNAMIC_DRAW)

	vert_attrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vert_attrib)
	gl.VertexAttribPointer(vert_attrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	previous_time := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		now := glfw.GetTime()
		elapsed := now - previous_time
		previous_time = now
		angle += elapsed * rotate_span
		// Render
		// rotation
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(vao_hexagon)

		gl.LineWidth(1.5);
		gl.DrawArrays(gl.LINE_STRIP, 0, int32(len(hexagon_vertex) / 3))

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}

	fmt.Println("end")
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyUp && action != glfw.Release {
	} else if key == glfw.KeyRight && action != glfw.Release {
		rotate_span += 0.1
	} else if key == glfw.KeyLeft && action != glfw.Release {
		rotate_span -= 0.1
	}
}
