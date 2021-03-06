package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "aaa0b", expected: "aab"},
		{input: "AAbc3", expected: "AAbccc"},
		{input: "ББbc3", expected: "ББbccc"},
		{input: "бб3bc0", expected: "ббббb"},
		{input: "世界2", expected: "世界界"},
		{input: "世БаG2", expected: "世БаGG"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{
		" ",
		` `,
		" a",
		"a ",
		` a`,
		`a `,
		"3abc",
		"45",
		"aaa10b",
		"aab12",
		"🐋abc",
		"2ava",
		`\`,
		"\n",
		`\n`,
		"\t",
		`\t`,
		"*",
		"-",
		"_",
	}

	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
