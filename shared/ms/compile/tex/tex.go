package tex

import (
	_ "embed"
	"gwcoffey/otis/shared/cfg"
	"gwcoffey/otis/shared/latex"
	"gwcoffey/otis/shared/ms"
	"strings"
)

func writeNewScene(i int, out *strings.Builder) {
	if i > 0 {
		out.WriteString(latex.Command("newscene", nil, nil))
	}
}

func writeScene(scene ms.Scene, out *strings.Builder) (err error) {
	text, err := scene.Text()
	if err != nil {
		return
	}
	out.WriteString(latex.Wrap(latex.Markdown(text)))
	if !strings.HasSuffix(text, "\n") {
		out.WriteString("\n")
	}
	return
}

func WorkToTex(work ms.Work, config cfg.Config) (tex string, err error) {

	out := strings.Builder{}
	out.WriteString(latex.Command("documentclass", []string{"novel", "courier"}, []string{"sffms"}))
	out.WriteString(latex.Command("frenchspacing", nil, nil))
	out.WriteString(latex.Command("author", nil, []string{config.Author.Name}))

	if name := config.Author.RealName; name != nil {
		out.WriteString(latex.Command("authorname", nil, []string{*name}))
	}
	if name := config.Author.Surname; name != nil {
		out.WriteString(latex.Command("surname", nil, []string{*name}))
	}
	out.WriteString(latex.Command("address", nil, []string{config.Address}))

	out.WriteString(latex.Command("title", nil, []string{work.Title()}))
	out.WriteString(latex.Command("runningtitle", nil, []string{work.RunningTitle()}))

	wcount, err := work.MsWordCount()
	if err != nil {
		return
	}
	out.WriteString(latex.Command("wordcount", nil, []string{wcount}))

	out.WriteString(latex.Command("begin", nil, []string{"document"}))

	if len(work.Chapters()) > 0 {
		for _, chapter := range work.Chapters() {
			out.WriteString("\n") // blank line before each chap for better readability
			if chapter.Number() == nil {
				out.WriteString(latex.Command("chapter*", nil, []string{chapter.Title()}))
			} else {
				out.WriteString(latex.Command("chapter", nil, []string{chapter.Title()}))
			}
			for i, scene := range chapter.Scenes() {
				writeNewScene(i, &out)
				err = writeScene(scene, &out)
				if err != nil {
					return
				}
			}
		}
	} else { // no chapters
		for i, scene := range work.Scenes() {
			writeNewScene(i, &out)
			err = writeScene(scene, &out)
			if err != nil {
				return
			}
		}
	}

	out.WriteString("\n")
	out.WriteString(latex.Command("end", nil, []string{"document"}))

	tex = out.String()
	return
}
