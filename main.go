package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"gwcoffey/otis/commands/compile"
	"gwcoffey/otis/commands/initcmd"
	"gwcoffey/otis/commands/mv"
	"gwcoffey/otis/commands/touch"
	"gwcoffey/otis/commands/wordcount"
)

var args struct {
	Init      *initcmd.Args   `arg:"subcommand:init" help:"initialize a new otis project"`
	Touch     *touch.Args     `arg:"subcommand:touch" help:"add a new scene"`
	Move      *mv.Args        `arg:"subcommand:mv" help:"move a scene or folder"`
	WordCount *wordcount.Args `arg:"subcommand:wc" help:"count words in your manuscript"`
	Compile   *compile.Args   `arg:"subcommand:compile" help:"compile the manuscript for submission"`
}

func main() {
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.Fail("missing subcommand")
	}

	switch {
	case args.Init != nil:
		initcmd.Init(args.Init)
	case args.WordCount != nil:
		wordcount.WordCount(args.WordCount)
	case args.Compile != nil:
		compile.Compile(args.Compile)
	case args.Touch != nil:
		touch.Touch(args.Touch)
	case args.Move != nil:
		mv.Mv(args.Move)
	default:
		panic(fmt.Sprintf("unexpected and unhandled command"))
	}
}
