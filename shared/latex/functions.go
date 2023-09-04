package latex

import (
	"strings"
)

func Command(command string, options []string, args []string) string {
	builder := strings.Builder{}
	builder.WriteString("\\")
	builder.WriteString(command)

	if len(options) > 0 {
		builder.WriteString("[")
		for i, option := range options {
			if i > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(EscapeOptionalArg(option))
		}
		builder.WriteString("]")
	}

	if len(args) > 0 {
		builder.WriteString("{")
		for i, arg := range args {
			if i > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(EscapeRequiredArg(arg))
		}
		builder.WriteString("}")
	}

	builder.WriteString("\n")

	return builder.String()
}

func Markdown(md string) string {
	return FormatMarkdown(EscapeText(md))
}

func Wrap(text string) (wrapped string) {
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
