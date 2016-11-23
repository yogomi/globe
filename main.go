package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"./misc"
	"./shader"
	"./idea"
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


	var line_vertexes = []mgl32.Vec3 {
		mgl32.Vec3{1.0, 1.0, 0.0},
		mgl32.Vec3{1.5, 1.5, 0.0},
	}

	fmt.Println(program)

	line := idea.NewIdea()
	line.Initialize(program)
	line.BindVertexes(line_vertexes)

	hexagon := idea.NewIdea()
	hexagon.Initialize(program)
	hexagon.BindVertexes(point.HexagonVertex())

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

		line.Draw()
		hexagon.Draw()

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
