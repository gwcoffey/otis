package rtf

import (
	"fmt"
	ms2 "gwcoffey/otis/ms"
	"strings"
)

// toRtfText prepares text for insertion into RTF; it:
// - reduces consecutive newlines to a single newline and escapes it
// - converts *emphasis* to underlines
// - escapes non-7bit-ascii characters,
func toRtfText(text string) string {
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

// writeScene writes a scene break (if needed) and then the scene itself
func writeScene(scidx int, scene ms2.Scene, out *strings.Builder) (err error) {
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
	out.WriteString(toRtfText(text))
	out.WriteString(`\par}`)
	out.WriteString("\n")
	return
}

func ManuscriptToHtml(m ms2.Manuscript) (rtf string, err error) {
	wcount, err := ms2.ApproximateWordCount(m)
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
	out.WriteString(m.AuthorRealName())
	out.WriteString("\t")
	out.WriteString(wcount + " words\\\n")

	// paragraph with address lines
	out.WriteString("\\pard\n")
	out.WriteString(strings.ReplaceAll(m.AuthorAddress(), "\n", "\\\n"))
	out.WriteString("\\\n")

	// paragraph double-spaced and centered
	out.WriteString(`\pard\sl480\slmult1\qc `)

	// output title and byline
	out.WriteString("\\\n\\\n\\\n\\\n\\\n\\\n\\\n\\\n" + strings.ToUpper(m.Title()))
	out.WriteString("\\\n")
	out.WriteString("By " + m.AuthorName())

	// start a new section with header
	out.WriteString(`\sect\sectd\sbknone\page`)
	out.WriteString(`{\header\pard\f0\fs24\qr `)
	out.WriteString(m.AuthorSurname())
	out.WriteString(" / ")
	out.WriteString(strings.ToUpper(m.RunningTitle()))
	out.WriteString(` / \chpgn`)
	out.WriteString(` \par}`)

	// content
	if len(m.Chapters()) > 0 {
		for chidx, chapter := range m.Chapters() {
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
		for scidx, scene := range m.Scenes() {
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
