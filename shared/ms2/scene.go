package ms2

type Scene interface {
	FileSystemObject
	Number() int
	Text() (string, error)
}

func (n *node) Number() int {
	if n.sceneCfg != nil {
		return n.sceneCfg.number
	}
	panic("attempt to read scene number from non-scene node")
}

func (n *node) Text() (string, error) {
	err := n.loadContent()
	if err != nil {
		return "", err
	}
	text := string(n.content)
	return text, nil
}
