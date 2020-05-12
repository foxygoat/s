// Package str provides string utilities.
package str

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	none = iota
	lower
	upper
	digit
	other
)

// CamelSplit splits a camelcase word and returns a list of words. It
// also supports digits. Both lower camel case and upper camel case are
// supported. As per Go naming conventions multiple Upper case letters
// together are considered an abbreviation and therefore a word.
//
// Examples:
//   ""                    → []
//   "lowercase"           → ["lowercase"]
//   "Class"               → ["Class"]
//   "MyClass"             → ["My", "Class"]
//   "MyC"                 → ["My", "C"]
//   "HTML"                → ["HTML"]
//   "PDFLoader"           → ["PDF", "Loader"]
//   "AString"             → ["A", "String"]
//   "SimpleXMLParser"     → ["Simple", "XML", "Parser"]
//   "vimRPCPlugin"        → ["vim", "RPC", "Plugin"]
//   "GL11Version"         → ["GL", "11", "Version"]
//   "99Bottles"           → ["99", "Bottles"]
//   "May5"                → ["May", "5"]
//   "BFG9000"             → ["BFG", "9000"]
//   "BöseÜberraschung"    → ["Böse", "Überraschung"]
//   "Two  spaces"         → ["Two", "  ", "spaces"]
//   "BadUTF8\xe2\xe2\xa1" → ["BadUTF8\xe2\xe2\xa1"]
//
// This code is inspired by https://github.com/fatih/camelcase (MIT licensed)
func CamelSplit(s string) []string {
	if !utf8.ValidString(s) {
		return []string{s}
	}
	words := []string{}
	runes := []rune{}
	lastClass := none
	for _, r := range s {
		class := classify(r)
		switch {
		case class == lastClass || lastClass == none:
			runes = append(runes, r)
		case lastClass == upper && class == lower:
			last := len(runes) - 1
			words = addWord(words, runes[:last])
			runes = []rune{runes[last], r}
		default:
			words = addWord(words, runes)
			runes = []rune{r}
		}
		lastClass = class
	}
	return addWord(words, runes)
}

func classify(r rune) int {
	switch {
	case unicode.IsLower(r):
		return lower
	case unicode.IsUpper(r):
		return upper
	case unicode.IsDigit(r):
		return digit
	}
	return other
}

func addWord(words []string, runes []rune) []string {
	if len(runes) > 0 {
		words = append(words, string(runes))
	}
	return words
}

// ToCamelCase converts a string to lower camel case (as used for unexported Go
// indentifiers), converting the entire first word to lowercase. The word
// boundaries are located as per CamelSplit. For example, "HelloWorld" becomes
// "helloWorld", "HTTPDir" becomes "httpDir". If the string contains invalid
// utf-8, the result is undefined.
func ToCamelCase(s string) string {
	splits := CamelSplit(s)
	if len(splits) == 0 {
		return ""
	}
	splits[0] = strings.ToLower(splits[0])
	return strings.Join(splits, "")
}
