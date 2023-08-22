package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"gwcoffey/otis/commands/echo"
	"gwcoffey/otis/commands/wordcount"
)

var args struct {
	Echo      *echo.Args      `arg:"subcommand:echo"`
	WordCount *wordcount.Args `arg:"subcommand:wordcount"`
}

func main() {
	p := arg.MustParse(&args)
	if p.Subcommand() == nil {
		p.Fail("missing subcommand")
	}

	switch {
	case args.Echo != nil:
		echo.Echo(args.Echo)
	case args.WordCount != nil:
		wordcount.WordCount(args.WordCount)
	default:
		panic(fmt.Sprintf("unexpected and unhandled command"))
	}
}
