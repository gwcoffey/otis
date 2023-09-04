package latex

import (
	"fmt"
	"regexp"
	"strings"
)

var optionalCommandEscapes = strings.NewReplacer(
	"]", "{]}",
	",", "{,}",
	"\n", "\\\\\n")

var requiredCommandEscapes = strings.NewReplacer(
	",", "{,}",
	"\n", "\\\\\n")

var textEscapes = strings.NewReplacer(
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

var emphasisPattern = regexp.MustCompile(`\*(.+?)\*`)
var blockquotePattern = regexp.MustCompile(`^>.*(?:\n>.*)*`)
var blockquoteCleanerPattern = regexp.MustCompile(`^>\s+`)

// EscapeOptionalArg textEscapes text for an argument to a command
func EscapeOptionalArg(text string) string {
	return optionalCommandEscapes.Replace(text)
}

// EscapeRequiredArg textEscapes text for an argument to a command
func EscapeRequiredArg(text string) string {
	return requiredCommandEscapes.Replace(text)
}

// EscapeText textEscapes characters that would otherwise be interpreted by latex, making the
// value safe to insert into a latex file as content.
func EscapeText(text string) string {
	return textEscapes.Replace(text)
}

func EscapeNewlines(text string) string {
	return strings.ReplaceAll(text, "\n", "\\\\\n")
}

// FormatMarkdown converts basic markdown into latex-formatted text. This very simple
// implementation only supports a small subset of markdown, namely:
//
//	*emphasis* 		-> \emph
//	> blockquotes 	-> \begin{quotation}...\end{quotation}
func FormatMarkdown(text string) string {
	text = emphasisPattern.ReplaceAllString(text, `\emph{$1}`)
	text = blockquotePattern.ReplaceAllStringFunc(text, func(match string) string {
		clean := blockquoteCleanerPattern.ReplaceAllString(match, "")
		return fmt.Sprintf("\\begin{quotation}\n%s\n\\end{quotation}\n", clean)
	})
	return text
}
