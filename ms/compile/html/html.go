package html

import (
	_ "embed"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	ms2 "gwcoffey/otis/ms"
	"html/template"
	"strings"
)

type templateData struct {
	Manuscript ms2.Manuscript
	WordCount  string
}

//go:embed output.html.tmpl
var templateText string

func loadTemplate() (tmpl *template.Template, err error) {
	tmpl, err = template.New("document").
		Funcs(template.FuncMap{
			"breaks": func(s string) template.HTML {
				return template.HTML(strings.Replace(template.HTMLEscapeString(s), "\n", "<br>", -1))
			},
			"markdown": func(s string) template.HTML {
				extensions := parser.CommonExtensions
				p := parser.NewWithExtensions(extensions)
				doc := p.Parse([]byte(s))

				htmlFlags := html.CommonFlags
				opts := html.RendererOptions{Flags: htmlFlags}
				renderer := html.NewRenderer(opts)

				return template.HTML(markdown.Render(doc, renderer))
			},
		}).
		Parse(templateText)
	return
}

func ManuscriptToHtml(m ms2.Manuscript) (html string, err error) {
	htemplate, err := loadTemplate()
	out := strings.Builder{}

	wordcount, err := ms2.ApproximateWordCount(m)
	if err != nil {
		return
	}

	err = htemplate.Execute(&out, templateData{Manuscript: m, WordCount: wordcount})
	if err != nil {
		return
	}

	html = out.String()
	return
}
