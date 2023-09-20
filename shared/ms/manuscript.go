package ms

import (
	"fmt"
	"gwcoffey/otis/shared/o/oerr"
	"path/filepath"
	"strings"
)

type manuscript struct {
	node *node
}

type Manuscript interface {
	fmt.Stringer
	Works() []Work
	ResolveSceneContainer(path string) (SceneContainer, error)
	ResolveScene(path string) (Scene, error)
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

func walkFolder(folder Folder, fn func(Folder) bool) {
	if fn(folder) {
		return
	}
	for _, f := range folder.Folders() {
		stop := fn(f)
		if stop {
			break
		}
		walkFolder(f, fn)
	}
}

func (m *manuscript) ResolveSceneContainer(path string) (result SceneContainer, err error) {
	var work Work
	for _, w := range m.Works() {
		if strings.HasPrefix(path, w.Path()) {
			work = w
			break
		}
	}

	if work == nil {
		err = oerr.FolderPathNotFound(path)
		return
	}

	if work.Path() == path {
		return work, nil
	}

	for _, folder := range work.Folders() {
		walkFolder(folder, func(f Folder) bool {
			if f.Path() == path {
				result = f
				return true
			} else {
				return false
			}
		})
	}

	if result == nil {
		err = oerr.FolderPathNotFound(path)
		return
	}

	return
}

func (m *manuscript) ResolveScene(path string) (result Scene, err error) {
	sceneContainer, err := m.ResolveSceneContainer(filepath.Dir(path))
	if err != nil {
		return
	}

	// search the container for a matching scene
	for _, scene := range sceneContainer.Scenes() {
		if scene.Path() == path {
			result = scene
			break
		}
	}

	if result == nil {
		err = oerr.ScenePathNotFound(path)
		return
	}

	return
}
