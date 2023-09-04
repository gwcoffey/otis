package html

import (
	_ "embed"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gwcoffey/otis/shared/cfg"
	"gwcoffey/otis/shared/ms"
	"html/template"
	"strings"
)

type templateData struct {
	Work      ms.Work
	Config    cfg.Config
	WordCount string
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

func WorkToHtml(config cfg.Config, work ms.Work) (html string, err error) {
	htemplate, err := loadTemplate()
	out := strings.Builder{}

	wordcount, err := work.MsWordCount()
	if err != nil {
		return
	}

	err = htemplate.Execute(&out, templateData{Work: work, Config: config, WordCount: wordcount})
	if err != nil {
		return
	}

	html = out.String()
	return
}
