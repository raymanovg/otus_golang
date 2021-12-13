package hw09structvalidator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatingRuleFromString(t *testing.T) {
	cases := []struct {
		name       string
		strRule    string
		expectRule ValidationRule
		expectErr  error
	}{
		{"len", "len:10", ValidationRule{RuleTypeLen, "10"}, nil},
		{"min", "min:1", ValidationRule{RuleTypeMin, "1"}, nil},
		{"max", "max:2", ValidationRule{RuleTypeMax, "2"}, nil},
		{"regexp", "regexp:^http:/\\w+.\\w+$", ValidationRule{RuleTypeRegexp, "^http:/\\w+.\\w+$"}, nil},
		{"in", "in:admin,stuff", ValidationRule{RuleTypeIn, "admin,stuff"}, nil},
		{"nested", "nested", ValidationRule{RuleTypeNested, ""}, nil},
		{"invalid nested rule value", "nested:1", ValidationRule{RuleTypeNested, ""}, nil},
		{"invalid len rule value", "len:", ValidationRule{}, ErrInvalidValidationRule},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualRule, actualErr := RuleFromString(c.strRule)
			require.ErrorIs(t, c.expectErr, actualErr)
			require.Equal(t, c.expectRule, actualRule)
		})
	}
}

func TestValidationRuleFromTag(t *testing.T) {
	cases := []struct {
		name        string
		tag         reflect.StructTag
		expectRules ValidationRules
		expectErr   error
	}{
		{
			"len|regexp",
			`validate:"len:10|regexp:\\d+"`,
			ValidationRules{
				{RuleTypeLen, "10"},
				{RuleTypeRegexp, "\\d+"},
			},
			nil,
		},
		{
			"min|max",
			`validate:"min:1|max:10"`,
			ValidationRules{
				{RuleTypeMin, "1"},
				{RuleTypeMax, "10"},
			},
			nil,
		},
		{
			"len,regexp,in",
			`validate:"len:3|regexp:\\d+|in:123,321,432"`,
			ValidationRules{
				{RuleTypeLen, "3"},
				{RuleTypeRegexp, "\\d+"},
				{RuleTypeIn, "123,321,432"},
			},
			nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualRules, actualErr := RulesFromTag(c.tag)
			require.ErrorIs(t, c.expectErr, actualErr)
			require.Equal(t, c.expectRules, actualRules)
		})
	}
}
