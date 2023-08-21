package echo

import(
	"flag"
	"fmt"
	"os"
)

func Echo(args []string) {
	fs := flag.NewFlagSet("echo", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s echo <str>\n\n", os.Args[0])
	}
	fs.Parse(args)

	vals := make([]any, fs.NArg())
	for i, arg := range fs.Args() {
		vals[i] = arg
	}
	fmt.Println(vals...)
}


