package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"gwcoffey/otis/commands/compile"
	"gwcoffey/otis/commands/echo"
	"gwcoffey/otis/commands/wordcount"
	"gwcoffey/otis/shared/cfg"
)

var args struct {
	ProjectPath *string         `arg:"--project" help:"path to the project directory (when not specified use current project)"`
	Echo        *echo.Args      `arg:"subcommand:echo"`
	WordCount   *wordcount.Args `arg:"subcommand:wordcount"`
	Compile     *compile.Args   `arg:"subcommand:compile"`
}

func getProjectRoot() (msPath string, err error) {
	if args.ProjectPath != nil {
		return *args.ProjectPath, nil
	} else {
		return cfg.FindProjectRoot()
	}
}

func main() {
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.Fail("missing subcommand")
	}

	projectPath, err := getProjectRoot()
	if err != nil {
		panic(err)
	}

	config, err := cfg.Load(projectPath)
	if err != nil {
		panic(err)
	}

	switch {
	case args.Echo != nil:
		echo.Echo(args.Echo)
	case args.WordCount != nil:
		wordcount.WordCount(config, args.WordCount)
	case args.Compile != nil:
		compile.Compile(config, args.Compile)
	default:
		panic(fmt.Sprintf("unexpected and unhandled command"))
	}
}
