package wordcount

import (
	"flag"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gwcoffey/otis/shared/cli"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

const maxWidth = 40
const indentSize = "  "

type node struct {
	Path     string
	Children []node
	Words    int
}

func wordCount(s string) int {
	return len(strings.Fields(s))
}

func readDir(path string) node {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	var children []node
	totalWords := 0

	for _, entry := range entries {
		subPath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			thisNode := readDir(subPath)
			totalWords += thisNode.Words
			children = append(children, thisNode)
		} else {
			if strings.HasSuffix(entry.Name(), ".md") {
				data, err := os.ReadFile(subPath)
				if err != nil {
					panic(err)
				}
				words := wordCount(string(data))
				totalWords += words
				children = append(children, node{Path: subPath, Children: []node{}, Words: words})
			}
		}
	}

	return node{Path: path, Children: children, Words: totalWords}
}

func prettify(str string) string {
	clean := str
	if filepath.Ext(str) == ".md" {
		re := regexp.MustCompile(`^(\d+)-(.*).md$`)
		matches := re.FindStringSubmatch(str)
		if len(matches) == 3 {
			num, err := strconv.Atoi(matches[1])
			if err != nil {
				panic("no scene number on " + str)
			}

			clean = fmt.Sprintf("%02d. %s", num, matches[2])
		}
	}
	words := strings.Split(clean, "-")
	if len(words) > 0 {
		words[0] = cases.Title(language.English).String(words[0])
	}
	return strings.Join(words, " ")
}

func truncate(str string) string {
	result := str
	if utf8.RuneCountInString(result) > maxWidth {
		result = result[0:maxWidth-1] + "…"
	}

	return result
}

func printDir(out *message.Printer, aNode *node, indent string) {

	name := truncate(indent + prettify(filepath.Base(aNode.Path)))
	format := fmt.Sprintf("%%-%d.%ds : %%7d\n", maxWidth, maxWidth)
	_, err := out.Printf(format, name, aNode.Words)
	if err != nil {
		panic(err)
	}

	for _, content := range aNode.Children {
		printDir(out, &content, indent+indentSize)
	}
}

func WordCount(args []string) {
	fs := flag.NewFlagSet("wordcount", flag.ExitOnError)
	fs.Usage = cli.UsageFn("wordcount [path]")
	cli.MustParse(fs, args)

	// use the full manuscript by default...
	rootPath := "manuscript/"

	// ...unless a path is specified
	if fs.NArg() == 1 {
		rootPath = fs.Arg(0)
	}

	out := message.NewPrinter(language.English)

	rootNode := readDir(rootPath)

	_, _ = out.Println()
	printDir(out, &rootNode, "")
	_, _ = out.Println()
}
