package s

import (
	"regexp"
	"unicode"
)

func ucFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// Make a Regex to say we only want letters and numbers
var reg = regexp.MustCompilePOSIX("[^a-zA-Z0-9]+")

func toName(in string) string {
	return ucFirst(reg.ReplaceAllString(in, ""))
}
