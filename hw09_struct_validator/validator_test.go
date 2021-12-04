package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Request struct {
		Application App   `validate:"nested"`
		Auth        Token `validate:"nested"`
	}
)

func TestSuccessValidate(t *testing.T) {
	cases := []struct {
		in                interface{}
		expectedErr       error
		expectedErrString string
	}{
		{
			User{
				"0123456789",
				"",
				50,
				"test@test.io",
				"stuff",
				[]string{"79123212121"},
				json.RawMessage{},
			},
			ValidationErrors{
				{"ID", errors.New("value len must be 36")},
			},
			"ID: value len must be 36",
		},
		{
			User{
				"0123456789",
				"",
				55,
				"test@test.io",
				"stuff",
				[]string{"79123212121"},
				json.RawMessage{},
			},
			ValidationErrors{
				{"ID", errors.New("value len must be 36")},
				{"Age", errors.New("value max must be 50")},
			},
			"ID: value len must be 36" +
				"\nAge: value max must be 50",
		},
		{
			User{
				"0123456789",
				"",
				55,
				"test.io",
				"stuff",
				[]string{"79123212121"},
				json.RawMessage{},
			},
			ValidationErrors{
				{"ID", errors.New("value len must be 36")},
				{"Age", errors.New("value max must be 50")},
				{"Email", errors.New("value is not matched by expression ^\\w+@\\w+\\.\\w+$")},
			},
			"ID: value len must be 36" +
				"\nAge: value max must be 50" +
				"\nEmail: value is not matched by expression ^\\w+@\\w+\\.\\w+$",
		},
		{
			User{
				"0123456789",
				"",
				55,
				"test.io",
				"support",
				[]string{"911"},
				json.RawMessage{},
			},
			ValidationErrors{
				{"ID", errors.New("value len must be 36")},
				{"Age", errors.New("value max must be 50")},
				{"Email", errors.New("value is not matched by expression ^\\w+@\\w+\\.\\w+$")},
				{"Role", errors.New("value must be in admin,stuff")},
				{"Phones", errors.New("value len must be 11")},
			},
			"ID: value len must be 36\nAge: value max must be 50" +
				"\nEmail: value is not matched by expression ^\\w+@\\w+\\.\\w+$" +
				"\nRole: value must be in admin,stuff" +
				"\nPhones: value len must be 11",
		},
		{
			Response{200, "body content"},
			nil,
			"",
		},
		{
			Request{App{"1.0"}, Token{}},
			ValidationErrors{
				{"Application", ValidationErrors{{"Version", errors.New("value len must be 5")}}},
			},
			"Application.Version: value len must be 5",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			c := c
			t.Parallel()

			actualErr := Validate(c.in)

			require.Equal(t, c.expectedErr, actualErr)
			if actualErr != nil {
				require.Equal(t, c.expectedErrString, actualErr.Error())
			}
			_ = c
		})
	}
}

func TestInvalidValidationRuleForType(t *testing.T) {
	cases := []struct {
		name              string
		in                interface{}
		expectedErr       error
		expectedErrString string
	}{
		{
			"invalid 'len' validate for int",
			struct {
				ID int `validate:"len:10"`
			}{1},
			ValidationErrors{
				ValidationError{"ID", errors.New("unknown validation rule 'len:10' for int value")},
			},
			"ID: unknown validation rule 'len:10' for int value",
		},
		{
			"invalid 'regexp' validate for int",
			struct {
				ID int `validate:"regexp:\\d+"`
			}{1},
			ValidationErrors{
				ValidationError{"ID", errors.New("unknown validation rule 'regexp:\\d+' for int value")},
			},
			"ID: unknown validation rule 'regexp:\\d+' for int value",
		},
		{
			name: "invalid 'in' validate for int",
			in: struct {
				ID int `validate:"in:foo,bar"`
			}{1},
			expectedErr: ValidationErrors{
				ValidationError{
					"ID",
					fmt.Errorf(
						"invalid validation rule 'in:foo,bar': %w",
						&strconv.NumError{Func: "Atoi", Num: "foo", Err: strconv.ErrSyntax},
					),
				},
			},
			expectedErrString: "ID: invalid validation rule 'in:foo,bar': strconv.Atoi: parsing \"foo\": invalid syntax",
		},
		{
			name: "invalid 'max' validate for string",
			in: struct {
				Name string `validate:"max:15"`
			}{"Foo"},
			expectedErr: ValidationErrors{
				ValidationError{
					"Name",
					errors.New("unknown validation rule max:15 for string value"),
				},
			},
			expectedErrString: "Name: unknown validation rule max:15 for string value",
		},
		{
			name: "invalid 'min' validate for string",
			in: struct {
				Name string `validate:"min:10"`
			}{"Foo"},
			expectedErr: ValidationErrors{
				ValidationError{
					"Name",
					errors.New("unknown validation rule min:10 for string value"),
				},
			},
			expectedErrString: "Name: unknown validation rule min:10 for string value",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c := c
			t.Parallel()

			actualErr := Validate(c.in)

			require.Equal(t, c.expectedErr, actualErr)
			if actualErr != nil {
				require.Equal(t, c.expectedErrString, actualErr.Error())
			}
			_ = c
		})
	}
}
