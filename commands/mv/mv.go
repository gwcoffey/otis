package mv

import (
	"fmt"
	ms2 "gwcoffey/otis/ms"
	"gwcoffey/otis/oerr"
	"os"
	"path/filepath"
	"regexp"
)

type Args struct {
	Path       string  `arg:"positional,required" help:"the scene to move"`
	TargetPath *string `arg:"positional" help:"the target path to move to"`
	At         *int    `arg:"--at,-a" help:"the scene number at which to insert"`
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

func Mv(args *Args) {
	m, err := ms2.LoadHere()
	if err != nil {
		panic(err)
	}

	path, err := filepath.Abs(args.Path)
	if err != nil {
		panic(err)
	}

	scene, err := m.ResolveScene(path)
	if err != nil {
		panic(err)
	}

	if args.TargetPath != nil {
		target, err := m.ResolveFolder(*args.TargetPath)
		if err != nil {
			panic(err)
		}

		sceneNumber := targetSceneNumber(args, target)

		err = moveToPath(scene, target, sceneNumber)
	} else if args.At != nil {
		err = moveToSceneNumber(scene, *args.At)
	} else {
		panic(oerr.PathOrAtRequired())
	}

}

var namePrefixRegex = regexp.MustCompile(`^\d+`)

func moveToSceneNumber(scene ms2.Scene, sceneNumber int) (err error) {
	if sceneNumber < scene.Number() {
		err = ms2.MakeRoomForScene(scene.Folder().Scenes(), sceneNumber)
		if err != nil {
			return
		}

		// TODO: New plan:
		// 1. copy the container to a temp location
		// 2. in the copy, move the operand scene to a temp location
		// 3. re-number the remaining scenes to make room
		// 4. move the operand scene into place
		// 5. atomically replace the original container with the new one
		// ... if at any point we fail, we clean up the temp files and the original
		// ... container is un-touched
		// re-order things, and then atomically move it back

		// when moving backwards, by the time we've made the space, the scene we're moving has
		// moved up by one
		fromScenePath := namePrefixRegex.ReplaceAllString(scene.Path(), fmt.Sprintf("%02d", scene.Number()+1))
		toScenePath := namePrefixRegex.ReplaceAllString(scene.Path(), fmt.Sprintf("%02d", sceneNumber))
		err = os.Rename(fromScenePath, toScenePath)
		return
	} else if sceneNumber > scene.Number() {
		// when moving forward, we have to fill in the space we left behind so the actual
		// final destination is one less than requested
		err = ms2.MakeRoomForScene(scene.Folder().Scenes(), sceneNumber)
		if err != nil {
			return
		}

		err = os.Rename(scene.Path(), namePrefixRegex.ReplaceAllString(scene.Path(), fmt.Sprintf("%02d", sceneNumber)))
		return
	} else {
		// move to self -- nothing to do
		return
	}
}

func moveToPath(scene ms2.Scene, target ms2.Folder, sceneNumber int) (err error) {
	err = ms2.MakeRoomForScene(target.Scenes(), sceneNumber)
	if err != nil {
		return
	}

	err = os.Rename(scene.Path(), target.Path())
	return
}
