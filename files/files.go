package files

import (
	"bufio"
	"bytes"
	"gwcoffey/otis/ms"
	"os"
	"path/filepath"
	"strings"
)

// ReadFileWithImport reads and returns the contents of file at path. If the file starts with
// a line of this form:
//
//	#import a/file/path
//
// The line is replaced with the contents of the referenced file.
func ReadFileWithImport(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	if !scanner.Scan() {
		return data, nil
	}

	firstLine := scanner.Text()
	if !strings.HasPrefix(firstLine, "#import ") {
		return data, nil
	}

	importPath := strings.TrimSpace(strings.TrimPrefix(firstLine, "#import "))
	importData, err := os.ReadFile(filepath.Join(filepath.Base(path), importPath))
	if err != nil {
		return nil, err
	}

	var result bytes.Buffer
	result.Write(importData)
	for scanner.Scan() {
		result.Write(scanner.Bytes())
		result.WriteByte('\n')
	}

	return result.Bytes(), scanner.Err()
}

// TmpDir returns the path to the temporary build directory of a given manuscript, creating
// it if necessary
func TmpDir(manuscript ms.Manuscript) (path string, err error) {
	path = filepath.Join(manuscript.Path(), ".build")
	// make sure it exists
	err = os.MkdirAll(path, os.ModePerm)
	return
}

// DistDir returns the path to the distribution directory of a given manuscript, creating
// it if necessary
func DistDir(manuscript ms.Manuscript) (path string, err error) {
	path = filepath.Join(manuscript.Path(), "dist")
	// make sure it exists
	err = os.MkdirAll(path, os.ModePerm)
	return
}
