package compile

import (
	_ "embed"
	"fmt"
	"gwcoffey/otis/shared/cfg"
	"gwcoffey/otis/shared/latex"
	"gwcoffey/otis/shared/ms"
	"os"
	"strings"
	"text/template"
)

type Args struct {
	Submission bool `help:"compile for submission"`
}

type templParams struct {
	Config     cfg.Config
	Manuscript ms.Dir
}

//go:embed output.tex.tmpl
var templateText string

func Command(command string, value string) string {
	escaped := latex.EscapeNewlines(latex.Escape(value))
	return fmt.Sprintf("\\%s{%s}", command, escaped)
}

func Multiline(lines []string) string {
	return strings.Join(lines, "\\\n")
}

func Markdown(md string) string {
	return latex.Escape(latex.FormatMarkdown(md))
}

func Compile(args *Args) {
	config := cfg.FindAndLoad()
	manuscript := ms.Load()

	tmpl, err := template.
		New("document").
		Funcs(template.FuncMap{
			"command":   Command,
			"multiline": Multiline,
			"markdown":  latex.Escape,
		}).
		Parse(templateText)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, templParams{Config: config, Manuscript: manuscript})
	if err != nil {
		panic(err)
	}
}
