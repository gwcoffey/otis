package ms2

type FileSystemObject interface {
	Path() string
}

func (n *node) Path() string {
	return n.path
}
