package wordcount

import (
	"errors"
	"flag"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gwcoffey/otis/shared/cli"
	"gwcoffey/otis/shared/ms"
	"os"
	"strings"
	"unicode/utf8"
)

const maxWidth = 40
const indentSize = "  "

func sceneWordCount(scene ms.Scene) int {
	text, err := scene.Text()
	if err != nil {
		panic(err)
	}
	return len(strings.Fields(*text))
}

func dirWordCount(dir ms.Dir) int {
	total := 0
	for _, scene := range dir.Scenes() {
		total += sceneWordCount(scene)
	}
	for _, subdir := range dir.SubDirs() {
		total += dirWordCount(subdir)
	}
	return total
}

func truncate(str string) string {
	result := str
	if utf8.RuneCountInString(result) > maxWidth {
		result = result[0:maxWidth-1] + "…"
	}
	return result
}

func printDir(dir ms.Dir, indent string) {
	printLine(truncate(indent+dir.Name()), dirWordCount(dir), true)
	for _, subdir := range dir.SubDirs() {
		printDir(subdir, indent+indentSize)
	}
	for _, scene := range dir.Scenes() {
		printScene(scene, indent+indentSize)
	}
}

func printScene(scene ms.Scene, indent string) {
	printLine(truncate(indent+scene.Name()), sceneWordCount(scene), false)
}

func printLine(label string, count int, emphasize bool) {
	out := message.NewPrinter(language.English)
	format := fmt.Sprintf("%%-%d.%ds : %%7d\n", maxWidth, maxWidth)
	if emphasize {
		format = fmt.Sprintf("\033[94m%s\033[0m", format)
	}
	_, err := out.Printf(format, label, count)
	if err != nil {
		panic(err)
	}
}

func readParams(fs *flag.FlagSet) (*string, error) {
	if fs.NArg() == 0 {
		return nil, nil
	} else if fs.NArg() == 1 {
		result := fs.Arg(0)
		return &result, nil
	} else {
		return nil, errors.New("unexpected parameters")
	}
}

func findMsRoot(path *string) ms.Dir {
	root := ms.LoadMs()

	// if a path is specified dig for it
	if path != nil {
		for _, name := range strings.Split(*path, "/") {
			newRoot := root.SubDir(name)
			if newRoot == nil {
				panic(fmt.Sprintf("no such manuscript directory: %s", path))
			}
			root = *newRoot
		}
	}

	return root
}

func WordCount(args []string) {
	fs := flag.NewFlagSet("wordcount", flag.ExitOnError)
	fs.Usage = cli.UsageFn("wordcount [path]")
	cli.MustParse(fs, args)

	path, err := readParams(fs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	root := findMsRoot(path)
	printDir(root, "")
}
