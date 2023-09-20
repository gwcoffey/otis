package ms

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type SceneContainer interface {
	Scenes() []Scene
	Path() string
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

func FindNextSceneNumber(sceneContainer SceneContainer) int {
	sceneNumber := 0
	for _, scene := range sceneContainer.Scenes() {
		if sceneNumber <= scene.Number() {
			sceneNumber = scene.Number() + 1
		}
	}
	return sceneNumber
}

var nameReplaceRegex = regexp.MustCompile(`^\d\d`)

func MakeRoomForScene(scenes []Scene, sceneNumber int) (err error) {
	// iterate the scenes and move the specified scene number forward one spot
	// (and recursively make room for that move if needed)
	for i, scene := range scenes {
		if scene.Number() == sceneNumber {
			if len(scenes) > 1 {
				// recursively make room for the move we're about to make
				err = MakeRoomForScene(scenes[i+1:], scene.Number()+1)
				if err != nil {
					return
				}
			}

			// determine new filename by incrementing scene number prefix on existing name
			newName := nameReplaceRegex.ReplaceAllString(filepath.Base(scene.Path()), fmt.Sprintf("%02d", sceneNumber+1))

			// move the file
			err = os.Rename(scene.Path(), filepath.Join(filepath.Dir(scene.Path()), newName))
			if err != nil {
				return
			}
		}
	}
	return
}
