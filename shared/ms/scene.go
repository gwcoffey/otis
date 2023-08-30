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
	Number() int
	Text() (string, error)
}

type Scener interface {
	Scenes() []Scene
}

func (s *scene) String() string {
	return fmt.Sprintf("Scene{%s} of %s", "", s.folder)
}

func (s *scene) Path() string {
	return s.node.path
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
