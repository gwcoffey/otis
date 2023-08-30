package main

import (
	"errors"
	"fmt"
	"github.com/alexflint/go-arg"
	"gwcoffey/otis/commands/compile"
	"gwcoffey/otis/commands/echo"
	"gwcoffey/otis/commands/wordcount"
	"os"
	"path/filepath"
)

var args struct {
	MsPath    *string         `arg:"--manuscript" help:"path to the manuscript (when not specified use current project)"`
	Echo      *echo.Args      `arg:"subcommand:echo"`
	WordCount *wordcount.Args `arg:"subcommand:wordcount"`
	Compile   *compile.Args   `arg:"subcommand:compile"`
}

func findRoot() (path string, err error) {
	path, err = os.Getwd()
	if err != nil {
		return
	}

	for path != "/" {
		otisPath := filepath.Join(path, "otis.yml")
		if _, err = os.Stat(otisPath); err == nil || !os.IsNotExist(err) {
			return
		}
		path = filepath.Dir(path)
	}

	return "", errors.New("this is not an otis project")
}

func getMsPath() (msPath string, err error) {
	if args.MsPath != nil {
		return *args.MsPath, nil
	} else {
		return findRoot()
	}
}

func main() {
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.Fail("missing subcommand")
	}

	msPath, err := getMsPath()
	if err != nil {
		panic(err)
	}

	switch {
	case args.Echo != nil:
		echo.Echo(args.Echo)
	case args.WordCount != nil:
		wordcount.WordCount(msPath, args.WordCount)
	case args.Compile != nil:
		compile.Compile(args.Compile)
	default:
		panic(fmt.Sprintf("unexpected and unhandled command"))
	}
}
