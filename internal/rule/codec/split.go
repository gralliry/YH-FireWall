package codec

import (
	"slices"
	"strings"
)

const StandardSplitter = ","

var splitters = []rune{',', '\n', '\r', ';', '\t', ' '}

func Split(raw string) []string {
	return strings.FieldsFunc(raw, func(r rune) bool {
		return slices.Contains(splitters, r)
	})
}
