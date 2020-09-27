package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ErrInvalidString is returned when input string is in a wrong format.
var ErrInvalidString = errors.New("invalid string")

// Unpack unpacks string duplicating runes. Example: "a4bc2d5e" => "aaaabccddddde".
func Unpack(input string) (string, error) {
	var result strings.Builder
	lastLetter := ""
	i := 0
	inputRunes := []rune(input)
	for _, rune := range input {
		switch {
		case unicode.IsLetter(rune):
			lastLetter = string(rune)
			if i+1 == utf8.RuneCountInString(input) {
				result.WriteRune(rune)
				continue
			}
			if i+1 > utf8.RuneCountInString(input) {
				continue
			}
			nextSymbol := inputRunes[i+1]
			if !unicode.IsDigit(nextSymbol) {
				result.WriteRune(rune)
			}

		case unicode.IsDigit(rune):
			if lastLetter == "" {
				return "", ErrInvalidString
			}
			digit, _ := strconv.Atoi(string(rune))
			if digit > 1 {
				s := strings.Repeat(lastLetter, digit)
				result.WriteString(s)
			}
			lastLetter = ""

		default:
			return "", ErrInvalidString
		}
		i++
	}
	return result.String(), nil
}
