package touch

import (
	"fmt"
	"gwcoffey/otis/shared/ms"
	"gwcoffey/otis/shared/o"
	"gwcoffey/otis/shared/text"
	"os"
	"path/filepath"
)

type Args struct {
	Path string `arg:"positional,required" help:"where to put the scene"`
	Name string `arg:"positional,required" help:"the name of the scene"`
	At   *int   `arg:"--at,-a" help:"the scene number at which to insert"`
}

func targetSceneNumber(args *Args, sceneContainer ms.SceneContainer) int {
	var sceneNumber int
	if args.At != nil {
		sceneNumber = *args.At
	} else {
		sceneNumber = ms.FindNextSceneNumber(sceneContainer)
	}
	return sceneNumber
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
	err = ms.MakeRoomForScene(sceneContainer.Scenes(), sceneNumber)

	// add the new scene file
	err = createScene(path, sceneNumber, args.Name)
	if err != nil {
		panic(err)
	}
}
