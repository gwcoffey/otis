package ms

import "fmt"

type scene struct {
	node       *node
	folder     *folder
	manuscript *manuscript
	chapter    *chapter
}

type Scene interface {
	fmt.Stringer
	FileSystemObject
	Folder() Folder
	Number() int
	Text() (string, error)
}

func (s *scene) String() string {
	return fmt.Sprintf("Scene{number=%d, name=%s} of %s", s.Number(), s.PrettyFileName(), s.folder)
}

func (s *scene) Path() string {
	return s.node.path
}

func (s *scene) Folder() Folder {
	return s.folder
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
