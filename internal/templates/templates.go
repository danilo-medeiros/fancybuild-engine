package templates

import (
	"fmt"
	"regexp"
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
	var nounMap = map[string]string{
		"foot":       "feet",
		"tooth":      "teeth",
		"goose":      "geese",
		"man":        "men",
		"woman":      "women",
		"mouse":      "mice",
		"die":        "dice",
		"ox":         "oxen",
		"child":      "children",
		"person":     "people",
		"sheep":      "sheep",
		"fish":       "fish",
		"deer":       "deer",
		"moose":      "moose",
		"swine":      "swine",
		"buffalo":    "buffalo",
		"shrimp":     "shrimp",
		"trout":      "trout",
		"aircraft":   "aircraft",
		"watercraft": "watercraft",
		"hovercraft": "hovercraft",
		"spacecraft": "spacecraft",
		"photo":      "photos",
		"piano":      "pianos",
		"halo":       "halos",
		"cactus":     "cacti",
		"focus":      "foci",
		"phenomenon": "phenomena",
		"criterion":  "criteria",
	}

	// Exceptions
	if val, ok := nounMap[text]; ok {
		return val
	}

	// If the singular noun ends in ‑is, the plural ending is ‑es.
	{
		pattern := regexp.MustCompile("is$")

		if pattern.MatchString(text) {
			return pattern.ReplaceAllString(text, "es")
		}
	}

	// If the singular noun ends in ‑s, -ss, -sh, -ch, -x, or -z,
	// add ‑es to the end to make it plural.
	{
		pattern := regexp.MustCompile("s$|ss$|sh$|ch$|x$|z$|o$")

		if pattern.MatchString(text) {
			return text + "es"
		}
	}

	// If the noun ends with ‑f or ‑fe,
	// the f is often changed to ‑ve before adding the -s to form the plural version.
	{
		pattern := regexp.MustCompile("f$|fe$")

		if pattern.MatchString(text) {
			return pattern.ReplaceAllString(text, "ves")
		}
	}

	// If a singular noun ends in ‑y and the letter before the -y is a consonant,
	// change the ending to ‑ies to make the noun plural.
	{
		pattern1 := regexp.MustCompile("y$")
		pattern2 := regexp.MustCompile("ay$|ey$|iy$|oy$|uy$")

		if pattern1.MatchString(text) && !pattern2.MatchString(text) {
			return pattern1.ReplaceAllString(text, "ies")
		}
	}

	// If the singular noun ends in ‑o, add ‑es to make it plural.
	{
		pattern := regexp.MustCompile("o$")

		if pattern.MatchString(text) {
			return text + "es"
		}
	}

	//  To make regular nouns plural, add ‑s to the end.
	return text + "s"
}

func Empty(text string) bool {
	return len(text) == 0
}

// SimpleFormat - Runs a simple formatting in a string
func SimpleFormat(text string) string {
	result := text

	{
		pattern := regexp.MustCompile("\n+")
		result = pattern.ReplaceAllString(result, "\n")
	}

	lines := strings.Split(result, "\n")
	previousLineHasBreak := false

	//lint:ignore S1007 we want to define a regular expression with a conditional
	openingBlockPattern := regexp.MustCompile("{$|\\($")
	closingBlockPattern := regexp.MustCompile("}$|^\t+\\)$|^\\)$")
	commentPattern := regexp.MustCompile("^\t*//.*")

	for i, line := range lines {
		isPreviousLineComment := false
		isPreviousLineOpeningBlock := false
		isCurrentLineOpeningBlock := openingBlockPattern.MatchString(line)

		if i-1 >= 0 {
			isPreviousLineOpeningBlock = openingBlockPattern.MatchString(lines[i-1])
			isPreviousLineComment = commentPattern.MatchString(lines[i-1])
		}

		if isCurrentLineOpeningBlock && !previousLineHasBreak && !isPreviousLineOpeningBlock && !isPreviousLineComment {
			lines[i] = fmt.Sprintf("\n%s", line)
			continue
		}

		if i+1 >= len(lines) {
			continue
		}

		isNextLineClosingBlock := closingBlockPattern.MatchString(lines[i+1])
		isCurrentLineClosingBlock := closingBlockPattern.MatchString(line)

		if isCurrentLineClosingBlock && !isNextLineClosingBlock {
			lines[i] = fmt.Sprintf("%s\n", line)
			previousLineHasBreak = true
			continue
		}

		previousLineHasBreak = false
	}

	result = strings.Join(lines, "\n")
	return result
}
