package ms2

import (
	"fmt"
)

type folder struct {
	*node
	*manuscript
	parentFolder *folder
}

type Folder interface {
	fmt.Stringer
	FileSystemObject
	Scener
	Folderer
}

type Folderer interface {
	Folders() []Folder
}

func (f *folder) String() string {
	return fmt.Sprintf("Folder{%s}", f.Path())
}

func (f *folder) Path() string {
	return f.node.path
}

func (f *folder) Number() int {
	return f.node.fileNumber
}

func (f *folder) PrettyFileName() string {
	return f.node.prettyFileName()
}

func (f *folder) Folders() []Folder {
	return f.node.folders()
}

func (f *folder) Scenes() (scenes []Scene) {
	for _, child := range f.node.children {
		if !child.isDir {
			scenes = append(scenes, &scene{node: child, folder: f})
		}
	}
	return
}
