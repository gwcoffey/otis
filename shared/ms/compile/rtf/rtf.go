package rtf

import (
	"fmt"
	"gwcoffey/otis/shared/cfg"
	"gwcoffey/otis/shared/ms"
	"strings"
)

func escapeText(text string) string {
	builder := strings.Builder{}
	inNewline := false
	inEmphasis := false
	for _, r := range text {
		if r == '\n' {
			if !inNewline {
				builder.WriteString("\\\n")
			}
			inNewline = true
		} else if r == '*' {
			if inEmphasis {
				builder.WriteString(`\ul0}`)
				inEmphasis = false
			} else {
				builder.WriteString(`{\ul `)
				inEmphasis = true
			}
		} else if r <= 127 {
			inNewline = false
			builder.WriteRune(r)
		} else if r <= 256 {
			inNewline = false
			builder.WriteString(fmt.Sprintf("\\'%x", r))
		} else if r <= 32768 {
			inNewline = false
			builder.WriteString(fmt.Sprintf("\\uc1\\u%d*", r))
		} else {
			inNewline = false
			builder.WriteString(fmt.Sprintf("\\uc1\\u%d*", r-65536))
		}
	}
	return builder.String()
}

func writeScene(scidx int, scene ms.Scene, out *strings.Builder) (err error) {
	var text string
	if scidx > 0 {
		// output scene break
		out.WriteString(`{\pard\sl480\slmult1\qc #\par}`)
	}
	out.WriteString(`{\pard\fi720\sl480\slmult1\ql `)
	out.WriteString("\n")
	text, err = scene.Text()
	if err != nil {
		return
	}
	out.WriteString(escapeText(text))
	out.WriteString(`\par}`)
	out.WriteString("\n")
	return
}

func WorkToRtf(work ms.Work, config cfg.Config) (rtf string, err error) {
	wcount, err := work.MsWordCount()
	if err != nil {
		return
	}

	out := strings.Builder{}
	// start doc ansi charset
	out.WriteString(`{\rtf1\ansi`)
	// single font in table, courier new
	out.WriteString(`{\fonttbl\f0\fmodern\fcharset0 CourierNewPSMT;}`)
	// 1 inch margins
	out.WriteString(`\margl1440\margr1440`)
	// courier new 12pt throughout
	out.WriteString(`\f0\fs24`)

	// paragraph with right-aligned tab stop at 9360
	out.WriteString(`\pard\tqr\tx9360`)

	// output author name and wordcount
	if name := config.Author.RealName; name != nil {
		out.WriteString(*config.Author.RealName)
	} else {
		out.WriteString(config.Author.Name)
	}
	out.WriteString("\t")
	out.WriteString(wcount + " words\\\n")

	// paragraph with address lines
	out.WriteString("\\pard\n")
	out.WriteString(strings.ReplaceAll(config.Address, "\n", "\\\n"))
	out.WriteString("\\\n")

	// paragraph double-spaced and centered
	out.WriteString(`\pard\sl480\slmult1\qc `)

	// output title and byline
	out.WriteString("\\\n\\\n\\\n\\\n\\\n\\\n\\\n\\\n" + strings.ToUpper(work.Title()))
	out.WriteString("\\\n")
	out.WriteString("By " + config.Author.Name)

	// start a new section with header
	out.WriteString(`\sect\sectd\sbknone\page`)
	out.WriteString(`{\header\pard\f0\fs24\qr `)
	if config.Author.Surname != nil {
		out.WriteString(*config.Author.Surname)
	} else {
		out.WriteString(config.Author.Name)
	}
	out.WriteString(" / ")
	out.WriteString(strings.ToUpper(work.RunningTitle()))
	out.WriteString(` / \chpgn`)
	out.WriteString(` \par}`)

	// content
	if len(work.Chapters()) > 0 {
		for chidx, chapter := range work.Chapters() {
			if chidx > 0 {
				out.WriteString("\\page\n")
			}
			// paragraph double-spaced centered
			out.WriteString(`\pard\sl480\slmult1\qc `)
			if chapter.Number() != nil {
				// output chapter + number
				out.WriteString(fmt.Sprintf("\\\n\\\n\\\n\\\nChapter %d\\\n", *chapter.Number()))
			}
			// output chapter title
			out.WriteString(chapter.Title() + "\\\n\\\n\\\n")

			for scidx, scene := range chapter.Scenes() {
				err = writeScene(scidx, scene, &out)
				if err != nil {
					return
				}
			}
		}
	} else { // no chapters
		for scidx, scene := range work.Scenes() {
			err = writeScene(scidx, scene, &out)
			if err != nil {
				return
			}
		}
	}

	// output end marker
	out.WriteString(`\pard\sl480\slmult1\qc # # # # #`)

	// terminate RTF
	out.WriteString("}")

	rtf = out.String()
	return
}
