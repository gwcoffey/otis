package touch

import (
	"gwcoffey/otis/commands/work"
	"gwcoffey/otis/ms"
	"gwcoffey/otis/msfs"
	"os"
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
		num, err = msfs.NextSceneNumber(args.Path)
	}
	return
}

func makeRenameWorkList(path string, sceneNumber int) (workList work.List, err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	workList = work.List{}

	for _, entry := range entries {
		var n int
		n, err = msfs.FileNumber(entry.Name())
		if err != nil {
			return
		}
		if n >= sceneNumber {
			workList = work.AppendRename(workList, filepath.Join(path, entry.Name()), msfs.RenumberFilename(entry.Name(), n+1))
		}
	}

	return
}

func Touch(args *Args) {
	_, err := ms.LoadContaining(args.Path)
	if err != nil {
		panic(err)
	}

	// target either the end of the scene list or the provided scene number
	sceneNumber, err := targetSceneNumber(args)
	if err != nil {
		panic(err)
	}

	// make a work list for this add
	workList, err := makeRenameWorkList(args.Path, sceneNumber)
	workList = work.AppendAdd(workList, filepath.Join(args.Path, msfs.MakeFilename(args.Name, sceneNumber)))

	err = work.Execute(workList, args.Force)
	if err != nil {
		panic(err)
	}
}
