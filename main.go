package main

import (
	"errors"
	"fmt"
	"github.com/alexflint/go-arg"
	"gwcoffey/otis/commands/compile"
	"gwcoffey/otis/commands/initcmd"
	"gwcoffey/otis/commands/mkdir"
	"gwcoffey/otis/commands/mv"
	"gwcoffey/otis/commands/touch"
	"gwcoffey/otis/commands/wordcount"
	"gwcoffey/otis/oerr"
	"os"
)

var args struct {
	Init      *initcmd.Args   `arg:"subcommand:init" help:"initialize a new otis project"`
	Touch     *touch.Args     `arg:"subcommand:touch" help:"add a new scene"`
	MkDir     *mkdir.Args     `arg:"subcommand:mkdir" help:"add a new folder"`
	Move      *mv.Args        `arg:"subcommand:mv" help:"move a scene or folder"`
	WordCount *wordcount.Args `arg:"subcommand:wc" help:"count words in your manuscript"`
	Compile   *compile.Args   `arg:"subcommand:compile" help:"compile the manuscript for submission"`
}

func reportErrorAndExit(err error) {
	_, perr := fmt.Fprintf(os.Stderr, "otis: (error) %s\n", err.Error())
	if perr != nil {
		panic(perr)
	}

	var otisErr *oerr.OtisError
	ok := errors.As(err, &otisErr)
	if ok {
		os.Exit(int(otisErr.Code))
	} else {
		os.Exit(-1)
	}
}

func main() {
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.Fail("missing subcommand")
	}

	var err error

	switch {
	case args.Init != nil:
		err = initcmd.Init(args.Init)
	case args.WordCount != nil:
		err = wordcount.WordCount(args.WordCount)
	case args.Compile != nil:
		err = compile.Compile(args.Compile)
	case args.Touch != nil:
		err = touch.Touch(args.Touch)
	case args.MkDir != nil:
		err = mkdir.MkDir(args.MkDir)
	case args.Move != nil:
		err = mv.Mv(args.Move)
	}

	if err != nil {
		reportErrorAndExit(err)
	}
}
