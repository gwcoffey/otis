// Package echo provides a test command that simply echos back its arguments
package echo

import (
	"flag"
	"fmt"
	"gwcoffey/otis/shared/cli"
)

func Echo(args []string) {
	fs := flag.NewFlagSet("echo", flag.ExitOnError)
	fs.Usage = cli.UsageFn("echo <anything>...")
	cli.MustParse(fs, args)

	values := make([]any, fs.NArg())
	for i, arg := range fs.Args() {
		values[i] = arg
	}
	fmt.Println(values...)
}
