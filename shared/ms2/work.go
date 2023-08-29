package ms2

type work struct {
	*node
	*manuscript
}

type Work interface {
	Scener
	Folderer
	AllScenes() []Scene
	Title() string
	RunningTitle() string
	Author() string
	AuthorSurname() string
	Chapters() []Chapter
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

	w.node.walk(func(node *node) {
		if node.chapterMeta != nil {
			chapters = append(chapters, &chapter{node: node, work: w})
		}
	})

	return
}

func (w *work) Folders() []Folder {
	return w.node.folders()
}
