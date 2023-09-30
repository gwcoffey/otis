package tex

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

// escapeOptionalArg textEscapes text for an argument to a command
func escapeOptionalArg(text string) string {
	return optionalCommandEscapes.Replace(text)
}

// escapeRequiredArg textEscapes text for an argument to a command
func escapeRequiredArg(text string) string {
	return requiredCommandEscapes.Replace(text)
}

// escapeText textEscapes characters that would otherwise be interpreted by latex, making the
// value safe to insert into a latex file as content.
func escapeText(text string) string {
	return textEscapes.Replace(text)
}

// formatMarkdown converts basic markdown into latex-formatted text. This very simple
// implementation only supports a small subset of markdown, namely:
//
//	*emphasis* 		-> \emph
//	> blockquotes 	-> \begin{quotation}...\end{quotation}
func formatMarkdown(text string) string {
	text = emphasisPattern.ReplaceAllString(text, `\emph{$1}`)
	text = blockquotePattern.ReplaceAllStringFunc(text, func(match string) string {
		clean := blockquoteCleanerPattern.ReplaceAllString(match, "")
		return fmt.Sprintf("\\begin{quotation}\n%s\n\\end{quotation}\n", clean)
	})
	return text
}

// command outputs a latex command with (optional) arguments
func command(command string, options []string, args []string) string {
	builder := strings.Builder{}
	builder.WriteString("\\")
	builder.WriteString(command)

	if len(options) > 0 {
		builder.WriteString("[")
		for i, option := range options {
			if i > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(escapeOptionalArg(option))
		}
		builder.WriteString("]")
	}

	if len(args) > 0 {
		builder.WriteString("{")
		for i, arg := range args {
			if i > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(escapeRequiredArg(arg))
		}
		builder.WriteString("}")
	}

	builder.WriteString("\n")

	return builder.String()
}

// wrap word-wraps text so lines aren't so long that errors are hard to track down
func wrap(text string) (wrapped string) {
	var result strings.Builder
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if len(line) <= 80 {
			result.WriteString(line + " \n")
			continue
		}
		words := strings.Fields(line)
		currentLine := ""
		for _, word := range words {
			if len(currentLine)+len(word) > 80 {
				result.WriteString(currentLine + " \n")
				currentLine = ""
			}
			if len(currentLine) > 0 {
				currentLine += " "
			}
			currentLine += word
		}
		if currentLine != "" {
			result.WriteString(currentLine + " \n")
		}
	}
	return result.String()
}
