package tex

import (
	_ "embed"
	"gwcoffey/otis/shared/ms"
	"gwcoffey/otis/shared/o"
	"strings"
)

func writeScene(scidx int, scene ms.Scene, out *strings.Builder) (err error) {
	if scidx > 0 {
		out.WriteString(command("newscene", nil, nil))
	}
	text, err := scene.Text()
	if err != nil {
		return
	}
	out.WriteString(wrap(formatMarkdown(escapeText(text))))
	if !strings.HasSuffix(text, "\n") {
		out.WriteString("\n")
	}
	return
}

func WorkToTex(work ms.Work, otis o.Otis) (tex string, err error) {

	out := strings.Builder{}
	out.WriteString(command("documentclass", []string{"novel", "courier"}, []string{"sffms"}))
	out.WriteString(command("frenchspacing", nil, nil))
	out.WriteString(command("author", nil, []string{otis.AuthorName()}))

	if name := otis.AuthorRealName(); name != nil {
		out.WriteString(command("authorname", nil, []string{*name}))
	}
	if name := otis.AuthorSurname(); name != nil {
		out.WriteString(command("surname", nil, []string{*name}))
	}
	out.WriteString(command("address", nil, []string{otis.Address()}))

	out.WriteString(command("title", nil, []string{work.Title()}))
	out.WriteString(command("runningtitle", nil, []string{work.RunningTitle()}))

	wcount, err := work.MsWordCount()
	if err != nil {
		return
	}
	out.WriteString(command("wordcount", nil, []string{wcount}))

	out.WriteString(command("begin", nil, []string{"document"}))

	if len(work.Chapters()) > 0 {
		for _, chapter := range work.Chapters() {
			out.WriteString("\n") // blank line before each chap for better readability
			if chapter.Number() == nil {
				out.WriteString(command("chapter*", nil, []string{chapter.Title()}))
			} else {
				out.WriteString(command("chapter", nil, []string{chapter.Title()}))
			}
			for i, scene := range chapter.Scenes() {
				err = writeScene(i, scene, &out)
				if err != nil {
					return
				}
			}
		}
	} else { // no chapters
		for i, scene := range work.Scenes() {
			err = writeScene(i, scene, &out)
			if err != nil {
				return
			}
		}
	}

	out.WriteString("\n")
	out.WriteString(command("end", nil, []string{"document"}))

	tex = out.String()
	return
}
