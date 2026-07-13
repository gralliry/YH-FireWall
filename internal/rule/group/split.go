package group

import (
	"slices"
	"strings"
)

const StandardSplitter = ","

var splitters = []rune{',', '\n', '\r', ';', '\t', ' '}

func split(raw string) []string {
	return strings.FieldsFunc(raw, func(r rune) bool {
		return slices.Contains(splitters, r)
	})
}
