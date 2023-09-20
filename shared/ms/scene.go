package ms

import "fmt"

type scene struct {
	node    *node
	folder  *folder
	work    *work
	chapter *chapter
}

type Scene interface {
	fmt.Stringer
	FileSystemObject
	Container() SceneContainer
	Number() int
	Text() (string, error)
}

func (s *scene) String() string {
	return fmt.Sprintf("Scene{number=%d, name=%s} of %s", s.Number(), s.PrettyFileName(), s.folder)
}

func (s *scene) Path() string {
	return s.node.path
}

func (s *scene) Container() SceneContainer {
	if s.folder != nil {
		return s.folder
	} else if s.work != nil {
		return s.work
	} else {
		panic("scene with no container (this should not happen)")
	}
}

func (s *scene) Number() int {
	return s.node.fileNumber
}

func (s *scene) Text() (string, error) {
	err := s.node.loadContent()
	if err != nil {
		return "", err
	}
	text := string(s.node.content)
	return text, nil
}

func (s *scene) PrettyFileName() string {
	return s.node.prettyFileName()
}
