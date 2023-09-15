package ms

import (
	"errors"
	"fmt"
)

type SceneContainer interface {
	Scenes() []Scene
}

func validateManuscript(manuscript Manuscript) (err error) {
	if len(manuscript.Works()) == 0 {
		err = errors.New("manuscript has no works")
		return
	}
	for _, work := range manuscript.Works() {
		if err = validateWork(work); err != nil {
			return
		}
	}
	return
}

func validateWork(work Work) (err error) {
	if len(work.Chapters()) > 0 && work.AllScenes()[0].Path() != work.Chapters()[0].Scenes()[0].Path() {
		err = errors.New(fmt.Sprintf("work %s has scenes before the first chapter", work.Title()))
	}

	return
}

func MustLoad(path string) Manuscript {
	manuscript, err := Load(path)
	if err != nil {
		panic(err)
	}
	return manuscript
}

func Load(path string) (ms Manuscript, err error) {
	node, err := newRootNode(path)
	if err != nil {
		return
	}
	ms = &manuscript{node: node}
	if err = validateManuscript(ms); err != nil {
		return
	}
	return
}
