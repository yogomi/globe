package hexagonshell

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"../idea"
)

type Hexagonshell interface {
	Transport(v mgl32.Vec3)
	Rotate(start, end mgl32.Vec3, radian float32)
	Copy() Hexagonshell
	Draw() error
	DrawStage(stage uint) error

	GrowUp()
	growAndMakeChild(outline_vertexes []mgl32.Vec3, next_center_tile idea.Idea)
	grow(outline_vertexes []mgl32.Vec3)
}

type hexagonshell struct {
	child Hexagonshell
	tile idea.Idea
	stage uint
	angle float32
}

func CreateBaseShell(size, angle float32, program uint32) Hexagonshell {
	base_shell := &hexagonshell{nil, idea.NewIdea(), 0, angle}

	vertexes := HexagonVertex(size)

	base_shell.tile.Initialize(vertexes, program)
	base_shell.tile.SetPrimitiveType(gl.LINE_LOOP)

	return base_shell
}

func (shell *hexagonshell) Transport(v mgl32.Vec3) {
	if shell.child != nil {
		shell.child.Transport(v)
	}
	shell.tile.Transport(v)
}

func (shell *hexagonshell) Rotate(start, end mgl32.Vec3, radian float32) {
	if shell.child != nil {
		shell.child.Rotate(start, end, radian)
	}
	shell.tile.Rotate(start, end, radian)
}

func (shell *hexagonshell) Copy() Hexagonshell {
	new_shell := &hexagonshell{}

	new_shell.child = shell.child.Copy()
	new_shell.tile = shell.tile.Copy()
	new_shell.stage = shell.stage
	new_shell.angle = shell.angle

	return new_shell
}

func (shell *hexagonshell) Draw() error {
	return shell.DrawStage(0)
}

func (shell *hexagonshell) DrawStage(stage uint) error {
	fmt.Println(shell.stage)
	var err error = nil
	if stage == shell.stage {
		err = shell.tile.Draw()
		fmt.Println(len(shell.tile.Vertexes()) / 6)
		if err != nil {
			return err
		}
	} else {
		if shell.child != nil {
			err = shell.child.DrawStage(stage)
		}
		if err != nil {
			return err
		}
	}

	return err
}

// TODO Rotateに関して
func (shell *hexagonshell) GrowUp() {
	if len(shell.tile.Vertexes()) != 6 {
		fmt.Println("Child shell cannot GrouUp.")
		return
	}

	fmt.Println("GrowUp start")
	if shell.stage == 8 {
		fmt.Println("GrowUp is overwork")
		return
	}

	outline_vertexes := shell.tile.Vertexes()[:6]
	fmt.Println(outline_vertexes)

	next_center_tile := shell.tile
	shell.tile = idea.NewIdea()
	shell.tile.Initialize(nextBaseVertexes(outline_vertexes, shell.angle),
					next_center_tile.Program())
	shell.tile.SetPrimitiveType(gl.LINE_LOOP)

	if shell.child == nil {
		child_shell := &hexagonshell{nil, next_center_tile, 0, shell.angle}
		child_shell.grow(outline_vertexes)
		shell.child = child_shell
	} else {
		shell.child.growAndMakeChild(outline_vertexes, next_center_tile)
	}
	shell.stage++
	fmt.Println("GrowUp end")
}

func (shell *hexagonshell) grow(outline_vertexes []mgl32.Vec3) {
	center_tile := shell.tile.Copy()
	for i, _ := range outline_vertexes {
		pare := i + 1
		if pare == len(outline_vertexes) {
			pare = 0
		}
		new_tile := center_tile.Copy()

		transport_vector := outline_vertexes[i].Add(outline_vertexes[pare])
		new_tile.Transport(transport_vector)

		shell.tile.MergeIdea(new_tile)
	}
}

func (shell *hexagonshell) growAndMakeChild(outline_vertexes []mgl32.Vec3,
				next_center_tile idea.Idea) {

	temporary_tile := shell.tile
	shell.tile = next_center_tile
	shell.grow(outline_vertexes)
	if shell.child == nil {
		child_shell := &hexagonshell{nil, temporary_tile, 0, shell.angle}
		child_shell.grow(outline_vertexes)
		shell.child = child_shell
	} else {
		shell.child.growAndMakeChild(outline_vertexes, temporary_tile)
	}
	shell.stage++
}

// TODO Rotateに関して
func nextBaseVertexes(base_vertexes []mgl32.Vec3, angle float32) []mgl32.Vec3 {
	vertexes := make([]mgl32.Vec3, len(base_vertexes))
	for i, _ := range base_vertexes {
		pare := i + 1
		if pare == len(base_vertexes) {
			pare = 0
		}
		transport_vector := base_vertexes[i].Add(base_vertexes[pare])
		vertexes[i] = base_vertexes[i].Add(transport_vector)
	}
	return vertexes
}
