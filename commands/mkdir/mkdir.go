package mkdir

import (
	"gwcoffey/otis/ms"
	"gwcoffey/otis/msfs"
	"gwcoffey/otis/work"
	"path/filepath"
)

type Args struct {
	Path  string `arg:"positional,required" help:"where to put the folder"`
	Name  string `arg:"positional,required" help:"the name of the folder"`
	At    *int   `arg:"--at,-a" help:"the index at which to insert"`
	Force bool   `arg:"--force,-f" help:"move other files around without confirmation"`
}

func MkDir(args *Args) (err error) {
	_, err = ms.LoadContaining(args.Path)
	if err != nil {
		return
	}

	// if no --at is provided, go to the end of the list
	var index int
	if args.At != nil {
		index = *args.At
	} else {
		index, err = msfs.NextIndex(args.Path)
		if err != nil {
			return
		}
	}

	// make a work list for this add
	workList, err := msfs.MakeRoom(args.Path, index)
	workList = work.AddDir(workList, filepath.Join(args.Path, msfs.MakeDirname(args.Name, index)))

	err = work.Execute(workList, args.Force)
	if err != nil {
		return
	}

	return nil
}
