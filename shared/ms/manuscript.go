package ms

import "fmt"

type manuscript struct {
	node *node
}

type Manuscript interface {
	fmt.Stringer
	Works() []Work
}

func (m *manuscript) String() string {
	// path in practice will just be "manuscript/" but in tests it is more useful
	return fmt.Sprintf("Manuscript{%s}", m.node.path)
}

func (m *manuscript) Works() (works []Work) {
	// if the root is a work, add it
	if m.node.workMeta != nil {
		works = append(works, &work{node: m.node, manuscript: m})
	}

	// if the first-level children are works, add them
	for _, child := range m.node.children {
		if child.workMeta != nil {
			works = append(works, &work{node: child, manuscript: m})
		}
	}

	return
}
