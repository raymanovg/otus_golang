package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if err := validate(input); err != nil {
		return "", err
	}

	length := len(input)
	var b strings.Builder
	var last rune

	for i, r := range input {
		repeat := 1
		isDigit := unicode.IsDigit(r)
		if isDigit {
			repeat = encodeDigitRune(r)
		}
		if unicode.IsLetter(last) {
			b.WriteString(strings.Repeat(encodeLetterRune(last), repeat))
		}
		if !isDigit && length-1 == i {
			b.WriteString(encodeLetterRune(r))
		}
		last = r
	}

	return b.String(), nil
}

func validate(input string) error {
	var last rune
	for _, r := range input {
		if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
			return ErrInvalidString
		}
		if unicode.IsDigit(r) && !unicode.IsLetter(last) {
			return ErrInvalidString
		}
		last = r
	}
	return nil
}

func encodeDigitRune(r rune) int {
	d, _ := strconv.Atoi(encodeRune(r))
	return d
}

func encodeLetterRune(r rune) string {
	return encodeRune(r)
}

func encodeRune(r rune) string {
	buf := make([]byte, utf8.RuneLen(r))
	_ = utf8.EncodeRune(buf, r)
	return string(buf)
}
