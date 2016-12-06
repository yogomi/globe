package hexagonshell

import (
	"fmt"

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
}

type hexagonshell struct {
	belonging_shells [7]Hexagonshell
	tile idea.Idea
	base_vertexes [6]mgl32.Vec3
	stage uint
	angle float32
}

func CreateBaseShell(size, angle float32, program uint32) Hexagonshell {
	base_shell := &hexagonshell{[7]Hexagonshell{}, idea.NewIdea(), [6]mgl32.Vec3{}, 0, 0.0}

	vertexes := HexagonVertex(size)

	base_shell.tile.Initialize(vertexes, program)

	s := base_shell.base_vertexes[:]
	copy(s, vertexes[:6])

	return base_shell
}

func (shell *hexagonshell) Transport(v mgl32.Vec3) {
	if shell.stage == 0 {
		shell.tile.Transport(v)
	} else {
		for _, s := range shell.belonging_shells {
			if s != nil {
				s.Transport(v)
			}
		}
	}
	for i, vertex := range shell.base_vertexes {
		shell.base_vertexes[i] = vertex.Add(v)
	}
	shell.tile.RebindVertexes(append(shell.base_vertexes[:], shell.base_vertexes[0]))
}

func (shell *hexagonshell) Rotate(start, end mgl32.Vec3, radian float32) {
	if shell.stage == 0 {
		shell.tile.Rotate(start, end, radian)
	} else {
		for _, s := range shell.belonging_shells {
			if s != nil {
				s.Rotate(start, end, radian)
			}
		}
	}

	axis := end.Sub(start).Normalize()
	roter := mgl32.QuatRotate(radian, axis)
	for i, v := range shell.base_vertexes {
		shell.base_vertexes[i] = roter.Rotate(v.Sub(start)).Add(start)
	}
	shell.tile.RebindVertexes(append(shell.base_vertexes[:], shell.base_vertexes[0]))
}

func (shell *hexagonshell) Copy() Hexagonshell {
	new_shell := &hexagonshell{}

	for i, s := range shell.belonging_shells {
		if s != nil {
			new_shell.belonging_shells[i] = s.Copy()
		}
	}
	new_shell.tile = shell.tile.Copy()
	new_shell.base_vertexes = shell.base_vertexes
	new_shell.stage = shell.stage
	new_shell.angle = shell.angle

	return new_shell
}

func (shell *hexagonshell) Draw() error {
	return shell.DrawStage(0)
}

func (shell *hexagonshell) DrawStage(stage uint) error {
	var err error = nil
	if stage == shell.stage {
		err = shell.tile.Draw()
		if err != nil {
			return err
		}
	} else {
		for _, s := range shell.belonging_shells {
			if s != nil {
				err = s.DrawStage(stage)
			}
			if err != nil {
				return err
			}
		}
	}

	return err
}

// TODO Rotateに関して
func (shell *hexagonshell) GrowUp() {
	fmt.Println("GrowUp start")
	center_shell := shell.Copy()
	shell.belonging_shells[0] = center_shell
	for i, _ := range shell.base_vertexes {
		pare := i + 1
		if pare == len(shell.base_vertexes) {
			pare = 0
		}
		new_shell := center_shell.Copy()

		transport_vector := shell.base_vertexes[i].Add(shell.base_vertexes[pare])
		new_shell.Transport(transport_vector)

		shell.belonging_shells[i + 1] = new_shell
	}

	shell.base_vertexes = nextBaseVertexes(shell.base_vertexes, shell.angle)
	shell.tile.RebindVertexes(append(shell.base_vertexes[:], shell.base_vertexes[0]))

	shell.stage++
	fmt.Println("GrowUp end")
}

// TODO Rotateに関して
func nextBaseVertexes(base_vertexes [6]mgl32.Vec3, angle float32) [6]mgl32.Vec3 {
	vertexes := [6]mgl32.Vec3{}
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
