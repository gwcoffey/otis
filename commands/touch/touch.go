package touch

import (
	"fmt"
	ms2 "gwcoffey/otis/ms"
	"gwcoffey/otis/text"
	"os"
	"path/filepath"
)

type Args struct {
	Path string `arg:"positional,required" help:"where to put the scene"`
	Name string `arg:"positional,required" help:"the name of the scene"`
	At   *int   `arg:"--at,-a" help:"the scene number at which to insert"`
}

func targetSceneNumber(args *Args, folder ms2.Folder) int {
	var sceneNumber int
	if args.At != nil {
		sceneNumber = *args.At
	} else {
		sceneNumber = ms2.FindNextSceneNumber(folder)
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

func Touch(args *Args) {
	m, err := ms2.LoadHere()
	if err != nil {
		panic(err)
	}

	path, err := filepath.Abs(args.Path)
	if err != nil {
		panic(err)
	}

	sceneContainer, err := m.ResolveFolder(path)
	if err != nil {
		panic(err)
	}

	// target either the end of the scene list or the provided scene number
	sceneNumber := targetSceneNumber(args, sceneContainer)

	// if the target scene number is already in use, move things to make room for it
	err = ms2.MakeRoomForScene(sceneContainer.Scenes(), sceneNumber)

	// add the new scene file
	err = createScene(path, sceneNumber, args.Name)
	if err != nil {
		panic(err)
	}
}
