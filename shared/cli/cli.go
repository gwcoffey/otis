// Package cli provides helpers for handling cli arguments
package cli

import (
	"flag"
	"fmt"
	"os"
)

func PrintUsage(usageString string) {
	_, err := fmt.Fprintf(os.Stderr, "Usage: %s %s\n\n", os.Args[0], usageString)
	if err != nil {
		panic(err)
	}
}

func UsageFn(usageString string) func() {
	return func() {
		PrintUsage(usageString)
	}
}

func MustParse(fs *flag.FlagSet, args []string) {
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}
}
