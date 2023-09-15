package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"gwcoffey/otis/commands/compile"
	"gwcoffey/otis/commands/echo"
	"gwcoffey/otis/commands/initcmd"
	"gwcoffey/otis/commands/touch"
	"gwcoffey/otis/commands/wordcount"
	"gwcoffey/otis/shared/o"
)

var args struct {
	ProjectPath *string `arg:"--project" help:"path to the project directory (when not specified use current project)"`

	Init      *initcmd.Args   `arg:"subcommand:init"`
	Echo      *echo.Args      `arg:"subcommand:echo"`
	WordCount *wordcount.Args `arg:"subcommand:wordcount"`
	Touch     *touch.Args     `arg:"subcommand:touch"`
	Compile   *compile.Args   `arg:"subcommand:compile"`
}

func getProjectRoot() (msPath string, err error) {
	if args.ProjectPath != nil {
		return *args.ProjectPath, nil
	} else {
		return o.FindProjectRoot()
	}
}

func doInit() {
	projectPath := "."
	if args.ProjectPath != nil {
		projectPath = *args.ProjectPath
	}

	initcmd.Init(projectPath, args.Init)
}

func doProjectCommand() {
	projectPath, err := getProjectRoot()
	if err != nil {
		panic(err)
	}

	otis, err := o.Load(projectPath)
	if err != nil {
		panic(err)
	}

	switch {
	case args.Echo != nil:
		echo.Echo(args.Echo)
	case args.WordCount != nil:
		wordcount.WordCount(otis, args.WordCount)
	case args.Compile != nil:
		compile.Compile(otis, args.Compile)
	case args.Touch != nil:
		touch.Touch(otis, args.Touch)
	default:
		panic(fmt.Sprintf("unexpected and unhandled command"))
	}
}

func main() {
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.Fail("missing subcommand")
	}

	// init is the only command that doesn't need a project
	if args.Init != nil {
		doInit()
	} else {
		doProjectCommand()
	}
}
