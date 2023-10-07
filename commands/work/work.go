package work

import (
	"fmt"
	"gwcoffey/otis/cli"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type action int

const (
	rename action = iota
	addFile
	addDir
	move
)

type Work struct {
	action action
	path   string
	arg    string
}

type List []Work

var manuscriptPrefixRegex = regexp.MustCompile("^manuscript/")

func AppendRename(list List, path string, newName string) List {
	return append(list, Work{action: rename, path: path, arg: newName})
}

func AddFile(list List, path string) List {
	return append(list, Work{action: addFile, path: path})
}

func AddDir(list List, path string) List {
	return append(list, Work{action: addDir, path: path})
}

func AppendMove(list List, from string, to string) List {
	return append(list, Work{action: move, path: from, arg: to})
}

func PrintableString(items List) string {
	builder := strings.Builder{}
	for _, w := range items {
		builder.WriteString("  ")
		switch w.action {
		case rename:
			builder.WriteString("RENAME ")
			builder.WriteString(manuscriptPrefixRegex.ReplaceAllString(w.path, ""))
			builder.WriteString(" → ")
			builder.WriteString(w.arg)
		case addFile, addDir:
			builder.WriteString("   ADD ")
			builder.WriteString(manuscriptPrefixRegex.ReplaceAllString(w.path, ""))
		case move:
			builder.WriteString("  MOVE ")
			builder.WriteString(manuscriptPrefixRegex.ReplaceAllString(w.path, ""))
			builder.WriteString(" → ")
			builder.WriteString(manuscriptPrefixRegex.ReplaceAllString(w.arg, ""))
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func Execute(items List, force bool) (err error) {
	proceed := force || len(items) <= 1
	if !proceed {
		prompt := fmt.Sprintf("About to change:\n\n%s\nOK to proceed?", PrintableString(items))
		proceed = cli.Confirm(prompt)
	}

	if proceed {
		for _, workItem := range items {
			switch workItem.action {
			case rename:
				err = os.Rename(workItem.path, filepath.Join(filepath.Dir(workItem.path), workItem.arg))
				if err != nil {
					return
				}
			case addFile:
				var file *os.File
				file, err = os.OpenFile(workItem.path, os.O_CREATE|os.O_EXCL, 0666)
				if err != nil {
					return
				}
				err = file.Close()
				if err != nil {
					return
				}
			case addDir:
				err = os.Mkdir(workItem.path, 0777)
				if err != nil {
					return
				}
			case move:
				err = os.Rename(workItem.path, workItem.arg)
				if err != nil {
					return
				}
			}
		}
	}

	return
}
