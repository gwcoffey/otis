package ms

import (
	"fmt"
	"strings"
)

type authorMeta struct {
	Name     string  `yaml:"name"`
	Surname  *string `yaml:"surname"`
	RealName *string `yaml:"realname"`
}

type manuscriptMeta struct {
	Title        string     `yaml:"title"`
	RunningTitle *string    `yaml:"runningTitle"`
	Author       authorMeta `yaml:"author"`
	AddressLines string     `yaml:"address"`
}

type manuscript struct {
	path string
	meta manuscriptMeta
	node *node
}

type Manuscript interface {
	fmt.Stringer
	Title() string
	RunningTitle() string
	AuthorName() string
	AuthorSurname() string
	AuthorRealName() string
	AuthorAddress() string
	Path() string
	Folders() []Folder
	Chapters() []Chapter
	Scenes() []Scene
}

func (m *manuscript) String() string {
	// path in practice will just be "manuscript/" but in tests it is more useful
	return fmt.Sprintf("Manuscript{%s}", m.path)
}

func (m *manuscript) Title() string {
	return m.meta.Title
}

func (m *manuscript) RunningTitle() string {
	if m.meta.RunningTitle == nil {
		return m.Title()
	} else {
		return *m.meta.RunningTitle
	}
}

func (m *manuscript) AuthorName() string {
	return m.meta.Author.Name
}

func (m *manuscript) AuthorSurname() string {
	if m.meta.Author.Surname == nil {
		nameParts := strings.Split(m.AuthorName(), " ")
		return nameParts[len(nameParts)-1]
	} else {
		return *m.meta.Author.Surname
	}
}

func (m *manuscript) AuthorRealName() string {
	if m.meta.Author.RealName == nil {
		return m.AuthorName()
	} else {
		return *m.meta.Author.RealName
	}
}

func (m *manuscript) AuthorAddress() string {
	return m.meta.AddressLines
}

func (m *manuscript) Path() string {
	return m.path
}

func (m *manuscript) Folders() []Folder {
	return m.node.folders()
}

func (m *manuscript) Chapters() (chapters []Chapter) {

	count := 1
	m.node.walk(func(node *node) {
		if node.chapterMeta != nil {
			var number *int
			if node.chapterMeta.Numbered == nil || *node.chapterMeta.Numbered {
				newNumber := count
				number = &newNumber
				count++
			}
			chapters = append(chapters, &chapter{node: node, manuscript: m, number: number})
		}
	})

	return
}

func (m *manuscript) Scenes() (scenes []Scene) {
	m.node.walk(func(node *node) {
		if !node.isDir {
			scenes = append(scenes, &scene{node: node, manuscript: m})
		}
	})

	return
}
