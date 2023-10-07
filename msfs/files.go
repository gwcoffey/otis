package msfs

import (
	"bufio"
	"bytes"
	"fmt"
	"gwcoffey/otis/oerr"
	"gwcoffey/otis/text"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

var numberPrefixPattern = regexp.MustCompile(`^\d+`)

// FileNumber returns the number of a file following otis naming convention of `##-foo`
func FileNumber(path string) (num int, err error) {
	matches := numberPrefixPattern.FindStringSubmatch(filepath.Base(path))
	if len(matches) == 1 {
		return strconv.Atoi(matches[0])
	} else {
		return 0, oerr.MissingFileNumber(path)
	}
}

func FileNameWithoutNumber(path string) (name string, err error) {
	matches := numberPrefixPattern.FindStringSubmatch(filepath.Base(path))
	if len(matches) == 1 {
		name = numberPrefixPattern.ReplaceAllString(filepath.Base(path), "")
		name = strings.TrimPrefix(name, `-`)
	} else {
		err = oerr.MissingFileNumber(path)
	}
	return
}

func RenumberFilename(name string, newNum int) string {
	newName, err := FileNameWithoutNumber(name)
	if err != nil {
		// ignore unnumbered file error and just number it
		newName = name
	}
	return fmt.Sprintf("%02d-%s", newNum, newName)
}

func MakeFilename(name string, num int) string {
	return fmt.Sprintf("%02d-%s.md", num, text.ToKebab(name))
}

func MakeDirname(name string, num int) string {
	return fmt.Sprintf("%02d-%s", num, text.ToKebab(name))
}
