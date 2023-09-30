package tex

import (
	_ "embed"
	ms2 "gwcoffey/otis/ms"
	"strings"
)

func writeScene(scidx int, scene ms2.Scene, out *strings.Builder) (err error) {
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

func ManuscriptToTex(m ms2.Manuscript) (tex string, err error) {

	out := strings.Builder{}
	out.WriteString(command("documentclass", []string{"novel", "courier"}, []string{"sffms"}))
	out.WriteString(command("frenchspacing", nil, nil))
	out.WriteString(command("author", nil, []string{m.AuthorName()}))

	out.WriteString(command("authorname", nil, []string{m.AuthorRealName()}))
	out.WriteString(command("surname", nil, []string{m.AuthorSurname()}))
	out.WriteString(command("address", nil, []string{m.AuthorAddress()}))

	out.WriteString(command("title", nil, []string{m.Title()}))
	out.WriteString(command("runningtitle", nil, []string{m.RunningTitle()}))

	wcount, err := ms2.ApproximateWordCount(m)
	if err != nil {
		return
	}
	out.WriteString(command("wordcount", nil, []string{wcount}))

	out.WriteString(command("begin", nil, []string{"document"}))

	if len(m.Chapters()) > 0 {
		for _, chapter := range m.Chapters() {
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
		for i, scene := range m.Scenes() {
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
