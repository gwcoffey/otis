package ms2

type Manuscript interface {
	FileSystemObject
	Works() []Work
}

func (n *node) Works() (work []Work) {
	// if the root is a work, add it
	if n.workCfg != nil {
		work = append(work, n)
	}

	// if the first-level children are works, add them
	for _, child := range n.children {
		if child.workCfg != nil {
			work = append(work, child)
		}
	}

	return
}
