package action

import (
	"unicode"
	"unicode/utf8"
)

const eof = -1

func nextChar(input string) (rune, int) {
	if len(input) == 0 {
		return eof, 0
	}

	return utf8.DecodeRuneInString(input)
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
