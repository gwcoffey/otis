// Package echo provides a test command that simply echos back its arguments
package echo

import (
	"fmt"
)

type Args struct {
	Output []string `arg:"positional" help:"any string to echo back"`
}

func Echo(args *Args) {
	values := make([]any, len(args.Output))
	for i, arg := range args.Output {
		values[i] = arg
	}
	fmt.Println(values...)
}
