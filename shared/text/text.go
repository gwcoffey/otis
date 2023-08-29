package text

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func KebabToSentence(str string) string {
	words := strings.Split(str, "-")
	words[0] = cases.Title(language.English).String(words[0])
	return strings.Join(words, " ")
}
