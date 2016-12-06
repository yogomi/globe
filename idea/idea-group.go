package idea

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

type IdeaGroup interface {
	String() string
	Transport(v mgl32.Vec3)
	Rotate(start, end mgl32.Vec3, radian float32)
	Copy() IdeaGroup
	Draw() error
	AddIdea(string, Idea)
	AddIdeaGroup(string, IdeaGroup)
}

type idea_group struct {
	ideas map[string]Idea
	idea_groups map[string]IdeaGroup
}

func NewIdeaGroup() IdeaGroup {
	group := &idea_group{}
	group.ideas = make(map[string]Idea)
	group.idea_groups = make(map[string]IdeaGroup)
	return group
}

func (group *idea_group) String() string {
	return fmt.Sprintf("ideas: {num:%d map:%s}, idea_groups {num:%d map:%s}",
			len(group.ideas),
			group.ideas,
			len(group.idea_groups),
			group.idea_groups)
}

func (group *idea_group) Transport(vec mgl32.Vec3) {
	for _, g := range group.idea_groups {
		g.Transport(vec)
	}
	for _, o := range group.ideas {
		o.Transport(vec)
	}
}

func (group *idea_group) Rotate(start, end mgl32.Vec3, radian float32) {
	for _, g := range group.idea_groups {
		g.Rotate(start, end, radian)
	}
	for _, o := range group.ideas {
		o.Rotate(start, end, radian)
	}
}

func (group *idea_group) Copy() IdeaGroup {
	new_group := &idea_group{}
	new_group.ideas = make(map[string]Idea)
	new_group.idea_groups = make(map[string]IdeaGroup)

	for name, g := range group.idea_groups {
		new_group.idea_groups[name] = g.Copy()
	}
	for name, o := range group.ideas {
		new_group.ideas[name] = o.Copy()
	}
	return new_group
}

func (group *idea_group) Draw() error {
	for _, g := range group.idea_groups {
		g.Draw()
	}
	for _, o := range group.ideas {
		o.Draw()
	}
	return nil
}

func (group *idea_group) AddIdea(name string, o Idea) {
	fmt.Println("AddIdea")
	group.ideas[name] = o
}

func (group *idea_group) AddIdeaGroup(name string, g IdeaGroup) {
	fmt.Println("AddIdeaGroup")
	group.idea_groups[name] = g
}
