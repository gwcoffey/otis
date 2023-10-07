package msfs

import (
	"gwcoffey/otis/commands/work"
	"os"
	"path/filepath"
)

// TmpDir returns the path to the temporary build directory of a given manuscript, creating
// it if necessary
func TmpDir(msPath string) (path string, err error) {
	path = filepath.Join(msPath, ".build")
	// make sure it exists
	err = os.MkdirAll(path, os.ModePerm)
	return
}

// DistDir returns the path to the distribution directory of a given manuscript, creating
// it if necessary
func DistDir(msPath string) (path string, err error) {
	path = filepath.Join(msPath, "dist")
	// make sure it exists
	err = os.MkdirAll(path, os.ModePerm)
	return
}

func LastIndex(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	num := 0

	for _, entry := range entries {
		var n int
		n, err = FileNumber(entry.Name())
		if err != nil {
			// ignore unnumbered files
			continue
		}
		if n >= num {
			num = n
		}
	}

	return num, nil
}

func NextIndex(dir string) (int, error) {
	num, err := LastIndex(dir)
	return num + 1, err
}

// MakeRoom makes room in the given directory for a new item with the given index by moving
// existing items with later indices up one spot
func MakeRoom(path string, index int) (workList work.List, err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	workList = work.List{}

	for _, entry := range entries {
		var n int
		n, err = FileNumber(entry.Name())
		if err != nil {
			return
		}
		if n >= index {
			workList = work.AppendRename(workList, filepath.Join(path, entry.Name()), RenumberFilename(entry.Name(), n+1))
		}
	}

	return
}
