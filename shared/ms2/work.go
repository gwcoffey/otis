package ms2

type Work interface {
	Title() string
	RunningTitle() string
	Author() string
	AuthorSurname() string
	Chapters() []Chapter
	Scenes() []Scene
	Folders() []Folder
}

func (n *node) Title() string {
	return n.workCfg.Title
}

func (n *node) RunningTitle() string {
	return n.workCfg.RunningTitle
}

func (n *node) Author() string {
	return n.workCfg.Author
}

func (n *node) AuthorSurname() string {
	return n.workCfg.AuthorSurname
}

func (n *node) Scenes() []Scene {
	var result []Scene
	for _, child := range n.children {
		if child.sceneCfg != nil {
			result = append(result, child)
		}
	}
	return result
}

func (n *node) Chapters() []Chapter {
	var result []Chapter
	return result
}

func (n *node) Folders() []Folder {
	var result []Folder
	return result
}
