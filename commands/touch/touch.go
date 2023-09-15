package touch

import (
	"fmt"
	"gwcoffey/otis/shared/ms"
	"gwcoffey/otis/shared/o"
	"gwcoffey/otis/shared/text"
	"os"
	"path/filepath"
	"regexp"
)

type Args struct {
	Path string `arg:"positional,required" help:"where to put the scene"`
	Name string `arg:"positional,required" help:"the name of the scene"`
	At   *int   `arg:"--at,-a" help:"the scene number at which to insert"`
}

func findNextSceneNumber(f ms.SceneContainer) int {
	sceneNumber := 0
	for _, scene := range f.Scenes() {
		if sceneNumber <= scene.Number() {
			sceneNumber = scene.Number() + 1
		}
	}
	return sceneNumber
}

func targetSceneNumber(args *Args, sceneContainer ms.SceneContainer) int {
	var sceneNumber int
	if args.At != nil {
		sceneNumber = *args.At
	} else {
		sceneNumber = findNextSceneNumber(sceneContainer)
	}
	return sceneNumber
}

var nameReplaceRegex = regexp.MustCompile(`^\d\d`)

func makeRoomForScene(scenes []ms.Scene, sceneNumber int) (err error) {
	// iterate the scenes and move the specified scene number forward one spot
	// (and recursively make room for that move if needed)
	for i, scene := range scenes {
		if scene.Number() == sceneNumber {
			if len(scenes) > 1 {
				// recursively make room for the move we're about to make
				err = makeRoomForScene(scenes[i+1:], scene.Number()+1)
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

func createScene(path string, sceneNumber int, name string) (err error) {
	fileName := fmt.Sprintf("%02d-%s.md", sceneNumber, text.ToKebab(name))
	file, err := os.OpenFile(filepath.Join(path, fileName), os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return
	}

	err = file.Close()
	return
}

func Touch(otis o.Otis, args *Args) {
	m, err := otis.Manuscript()
	if err != nil {
		panic(err)
	}

	path, err := filepath.Abs(args.Path)
	if err != nil {
		panic(err)
	}

	sceneContainer, err := m.ResolveSceneContainer(path)
	if err != nil {
		panic(err)
	}

	// target either the end of the scene list or the provided scene number
	sceneNumber := targetSceneNumber(args, sceneContainer)

	// if the target scene number is already in use, move things to make room for it
	err = makeRoomForScene(sceneContainer.Scenes(), sceneNumber)

	// add the new scene file
	err = createScene(path, sceneNumber, args.Name)
	if err != nil {
		panic(err)
	}
}
