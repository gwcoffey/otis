package ms

import "fmt"

type work struct {
	node       *node
	manuscript *manuscript
}

type Work interface {
	fmt.Stringer
	FileSystemObject
	Scener
	Folderer
	AllScenes() []Scene
	Title() string
	RunningTitle() string
	Author() string
	AuthorSurname() string
	Chapters() []Chapter
}

func (w *work) String() string {
	return fmt.Sprintf("Work{%s}", w.node.path)
}

func (w *work) Path() string {
	return w.node.path
}

func (w *work) PrettyFileName() string {
	return w.node.prettyFileName()
}

func (w *work) Number() int {
	return w.node.fileNumber
}

func (w *work) AllScenes() (scenes []Scene) {
	w.node.walk(func(node *node) {
		if !node.isDir {
			scenes = append(scenes, &scene{node: node, work: w})
		}
	})

	return
}

func (w *work) Title() string {
	return w.node.workMeta.Title
}

func (w *work) RunningTitle() string {
	return w.node.workMeta.RunningTitle
}

func (w *work) Author() string {
	return w.node.workMeta.Author
}

func (w *work) AuthorSurname() string {
	return w.node.workMeta.AuthorSurname
}

func (w *work) Scenes() (scenes []Scene) {
	for _, child := range w.node.children {
		if !child.isDir {
			scenes = append(scenes, &scene{node: child, work: w})
		}
	}
	return
}

func (w *work) Chapters() (chapters []Chapter) {

	count := 1
	w.node.walk(func(node *node) {
		if node.chapterMeta != nil {
			var number *int
			if node.chapterMeta.Numbered == nil || *node.chapterMeta.Numbered {
				newNumber := count
				number = &newNumber
				count++
			}
			chapters = append(chapters, &chapter{node: node, work: w, number: number})
		}
	})

	return
}

func (w *work) Folders() []Folder {
	return w.node.folders()
}
