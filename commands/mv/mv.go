package mv

import (
	"gwcoffey/otis/commands/work"
	"gwcoffey/otis/ms"
	"gwcoffey/otis/msfs"
	"gwcoffey/otis/oerr"
	"math"
	"os"
	"path/filepath"
)

type Args struct {
	Path       string  `arg:"positional,required" help:"the scene to move"`
	TargetPath *string `arg:"positional" help:"the target path to move to"`
	At         *int    `arg:"--at,-a" help:"the scene number at which to insert"`
	Force      bool    `arg:"--force,-f" help:"move other files around without confirmation"`
}

func appendMoveToEndOfDir(workList work.List, scene string, dir string) (work.List, error) {
	lastSceneNumber, err := msfs.LastIndex(dir)
	if err != nil {
		return nil, err
	}

	workList = work.AppendMove(workList, scene, filepath.Join(dir, msfs.RenumberFilename(filepath.Base(scene), lastSceneNumber+1)))

	workList, err = closeHole(workList, scene)
	if err != nil {
		return nil, err
	}

	return workList, nil
}

func appendMoveToDirAt(workList work.List, scene string, dir string, sceneNumber int) (work.List, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// move existing scenes up by one to make a hole
	for _, entry := range entries {
		var num int
		num, nerr := msfs.FileNumber(filepath.Join(dir, entry.Name()))
		if nerr != nil {
			continue // ignore files with no file number
		}
		if num > sceneNumber {
			workList = work.AppendRename(workList, filepath.Join(dir, entry.Name()), msfs.RenumberFilename(entry.Name(), num+1))
		}
	}

	// move the scene into the hole
	workList = work.AppendMove(workList, scene, filepath.Join(dir, msfs.RenumberFilename(filepath.Base(scene), sceneNumber)))

	// close the hole left behind
	workList, err = closeHole(workList, scene)
	if err != nil {
		return nil, err
	}

	return workList, nil
}

func closeHole(workList work.List, scene string) (work.List, error) {
	dir := filepath.Dir(scene)
	sceneNumber, err := msfs.FileNumber(scene)
	if err != nil {
		return workList, nil // nothing to do, but not really an error
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return workList, err
	}

	for _, entry := range entries {
		var num int
		num, nerr := msfs.FileNumber(entry.Name())
		if nerr != nil {
			continue // just ignore files with no number
		}
		if num > sceneNumber {
			workList = work.AppendRename(workList, filepath.Join(dir, entry.Name()), msfs.RenumberFilename(entry.Name(), num-1))
		}
	}

	return workList, nil
}

func appendMoveInSameDir(workList work.List, manuscript ms.Manuscript, scene string, sceneNumber int) (work.List, error) {
	tmp, err := msfs.TmpDir(manuscript.Path())
	if err != nil {
		return nil, err
	}

	tmpFile := filepath.Join(tmp, filepath.Base(scene))
	originalSceneNumber, nerr := msfs.FileNumber(scene)
	if nerr != nil {
		// we'll move the scene into place, but there's no hole left behind
		// so setting this huge so nothing will be ahead of it
		originalSceneNumber = math.MaxInt
	}

	lastScene, err := msfs.LastIndex(filepath.Dir(scene))
	if err != nil {
		return nil, err
	}

	if lastScene == originalSceneNumber {
		lastScene = lastScene - 1
	}

	sceneNumber = int(math.Min(float64(sceneNumber), float64(lastScene+1)))

	// if we're not moving anything, short circuit
	if sceneNumber == originalSceneNumber {
		return workList, nil
	}

	// otherwise move the scene to a temp location and then make a space for it and move
	// it back in
	workList = work.AppendMove(workList, scene, tmpFile)

	// close the hole
	entries, err := os.ReadDir(filepath.Dir(scene))
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		num, nerr := msfs.FileNumber(entry.Name())
		if num == originalSceneNumber || nerr != nil {
			// ignore the file we're moving and any unnumbered file
			continue
		} else if num >= originalSceneNumber && num <= sceneNumber {
			// scenes between the scene and its new position need to move down to fill the space
			workList = work.AppendRename(workList, filepath.Join(filepath.Dir(scene), entry.Name()), msfs.RenumberFilename(entry.Name(), num-1))
		} else if num >= sceneNumber && num <= originalSceneNumber {
			// scenes after the target and before the scene need to move up to make space
			workList = work.AppendRename(workList, filepath.Join(filepath.Dir(scene), entry.Name()), msfs.RenumberFilename(entry.Name(), num+1))
		}
	}

	// move the scene to from tmp to the target
	workList = work.AppendMove(workList, tmpFile, filepath.Join(filepath.Dir(scene), msfs.RenumberFilename(filepath.Base(scene), sceneNumber)))

	return workList, nil
}

func Mv(args *Args) {
	manuscript, err := ms.LoadContaining(args.Path)
	if err != nil {
		panic(err)
	}

	workList := work.List{}
	if args.TargetPath != nil {
		if args.At == nil {
			workList, err = appendMoveToEndOfDir(workList, args.Path, *args.TargetPath)
		} else {
			workList, err = appendMoveToDirAt(workList, args.Path, *args.TargetPath, *args.At)
		}
	} else if args.At != nil {
		workList, err = appendMoveInSameDir(workList, manuscript, args.Path, *args.At)
	} else {
		panic(oerr.PathOrAtRequired())
	}

	if err != nil {
		panic(err)
	}

	err = work.Execute(workList, args.Force)
	if err != nil {
		panic(err)
	}
}
