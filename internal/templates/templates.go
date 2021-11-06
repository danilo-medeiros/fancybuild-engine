package templates

import (
	"fmt"
	"strings"
)

func Capitalize(text string) string {
	splitted := strings.Split(text, "")
	splitted[0] = strings.ToUpper(splitted[0])
	return strings.Join(splitted, "")
}

func Camelize(text ...string) string {
	for index, part := range text {
		if index == 0 {
			continue
		}

		text[index] = fmt.Sprintf("%s%s", strings.ToUpper(string(part[0])), part[1:])
	}

	return strings.Join(text, "")
}

func Slice(text string, start int, end int) string {
	return text[start:end]
}

func Pluralize(text string) string {
	// TODO: Implement pluralization rules: https://www.grammarly.com/blog/plural-nouns/
	return fmt.Sprintf("%ss", text)
}

func Empty(text string) bool {
	return len(text) == 0
}
