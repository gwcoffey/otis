package text

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
)

func KebabToSentence(str string) string {
	words := strings.Split(str, "-")
	words[0] = cases.Title(language.English).String(words[0])
	return strings.Join(words, " ")
}

var nonRomanRunsRegex = regexp.MustCompile(`[^a-z]+`)

func ToKebab(str string) (result string) {
	result = strings.ToLower(str)
	result = nonRomanRunsRegex.ReplaceAllString(result, "-")
	result = strings.Trim(result, "-")
	return
}
