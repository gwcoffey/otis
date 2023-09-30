package ms

import "fmt"

type chapter struct {
	node       *node
	manuscript *manuscript
	number     *int
}

type Chapter interface {
	fmt.Stringer
	Scenes() []Scene
	Title() string
	Number() *int
}

func (c *chapter) String() string {
	return fmt.Sprintf("Chapter{%s} of %s", c.Title(), c.manuscript)
}

// Scenes returns the ordered set of scenes in this chapter. Because chapters are really just waypoints
// that can appear at any point in the associated manuscript's filesystem hierarchy, we do this by walking
// the entire manuscript. As soon as we find the node that represents this chapter, we begin gathering scenes
// into the result. And once we encounter the next chapter node, we stop gathering.
func (c *chapter) Scenes() (scenes []Scene) {
	capturing := false
	c.manuscript.node.walk(func(node *node) {
		if node.chapterMeta == c.node.chapterMeta {
			capturing = true
		}
		if capturing {
			if node.chapterMeta != nil && node.chapterMeta != c.node.chapterMeta {
				capturing = false
			} else if !node.isDir {
				scenes = append(scenes, &scene{node: node, chapter: c})
			}
		}
	})

	return
}

func (c *chapter) Title() string {
	return c.node.chapterMeta.Title
}

func (c *chapter) Number() *int {
	return c.number
}

func (c *chapter) Path() string {
	return c.node.path
}
