package hw02unpackstring

import (
	"errors"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result []rune
	runes := []rune(input)

	if len(runes) == 0 {
		return "", nil
	}

	if unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		if !unicode.IsDigit(char) {
			result = append(result, char)
			continue
		}

		if i == 0 || unicode.IsDigit(runes[i-1]) {
			return "", ErrInvalidString
		}

		count := int(char - '0')
		if count == 0 {
			result = removeLastSymbol(result)
			continue
		}
		result = append(result, getRepeatedRune(runes[i-1], count-1)...)
		continue
	}

	return string(result), nil
}

func removeLastSymbol(result []rune) []rune {
	result = result[:len(result)-1]
	return result
}

func getRepeatedRune(r rune, count int) []rune {
	repeated := make([]rune, count)
	for i := range repeated {
		repeated[i] = r
	}
	return repeated
}
