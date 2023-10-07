package touch

import (
	"gwcoffey/otis/ms"
	"gwcoffey/otis/msfs"
	"gwcoffey/otis/work"
	"path/filepath"
)

type Args struct {
	Path  string `arg:"positional,required" help:"where to put the scene"`
	Name  string `arg:"positional,required" help:"the name of the scene"`
	At    *int   `arg:"--at,-a" help:"the scene number at which to insert"`
	Force bool   `arg:"--force,-f" help:"move other files around without confirmation"`
}

func targetSceneNumber(args *Args) (num int, err error) {
	if args.At != nil {
		num = *args.At
	} else {
		num, err = msfs.NextIndex(args.Path)
	}
	return
}

func Touch(args *Args) (err error) {
	_, err = ms.LoadContaining(args.Path)
	if err != nil {
		return
	}

	// target either the end of the scene list or the provided scene number
	sceneNumber, err := targetSceneNumber(args)
	if err != nil {
		return
	}

	// make a work list for this add
	workList, err := msfs.MakeRoom(args.Path, sceneNumber)
	workList = work.AddFile(workList, filepath.Join(args.Path, msfs.MakeFilename(args.Name, sceneNumber)))

	err = work.Execute(workList, args.Force)
	if err != nil {
		return
	}

	return nil
}
