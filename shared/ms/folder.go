package ms

import (
	"fmt"
)

type folder struct {
	node         *node
	manuscript   *manuscript
	parentFolder *folder
}

type Folder interface {
	fmt.Stringer
	FileSystemObject
	SceneContainer
	FolderContainer
	AllScenes() []Scene
}

type FolderContainer interface {
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

func (f *folder) AllScenes() (scenes []Scene) {
	f.node.walk(func(node *node) {
		if !node.isDir {
			scenes = append(scenes, &scene{node: node, folder: f})
		}
	})
	return
}
