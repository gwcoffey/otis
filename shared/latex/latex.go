package latex

import (
	"fmt"
	"regexp"
	"strings"
)

var escapes = strings.NewReplacer(
	`#`, `\#`,
	`$`, `\$`,
	`%`, `\%`,
	`&`, `\&`,
	`~`, `$\sim$`,
	`_`, `\_`,
	`^`, `$\textasciicircum$`,
	`\`, `$\backslash$`,
	`{`, `\{`,
	`}`, `\}`,
)

var sceneBreakPattern = regexp.MustCompile(`(?m)\s*^\s*\*\*\*\s*$\s*`)
var emphasisPattern = regexp.MustCompile(`\*(.+)\*`)
var blockquotePattern = regexp.MustCompile(`^>.*(?:\n>.*)*`)
var blockquoteCleanerPattern = regexp.MustCompile(`^>\s+`)

// Escape escapes characters that would otherwise be interpreted by latex, making the
// value safe to insert into a latex file as content.
func Escape(text string) string {
	return escapes.Replace(text)
}

func EscapeNewlines(text string) string {
	return strings.ReplaceAll(text, "\n", "\\\\\n")
}

// FormatMarkdown converts basic markdown into latex-formatted text. This very simple
// implementation only supports a small subset of markdown, namely:
//
//	*emphasis* 		-> \emph
//	> blockquotes 	-> \begin{quotation}...\end{quotation}
//	***				-> \newscene
func FormatMarkdown(text string) string {
	text = sceneBreakPattern.ReplaceAllString(text, "\n\n\\newscene\n\n")
	text = emphasisPattern.ReplaceAllString(text, `\emph{$1}`)
	text = blockquotePattern.ReplaceAllStringFunc(text, func(match string) string {
		clean := blockquoteCleanerPattern.ReplaceAllString(match, "")
		return fmt.Sprintf("\\begin{quotation}\n%s\n\\end{quotation}\n", clean)
	})
	return text
}
