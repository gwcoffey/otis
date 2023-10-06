package msfs

import (
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

func LastSceneNumber(dir string) (int, error) {
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

func NextSceneNumber(dir string) (int, error) {
	num, err := LastSceneNumber(dir)
	return num + 1, err
}
