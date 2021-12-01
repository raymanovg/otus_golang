package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
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

func TestValidate(t *testing.T) {
	tests := []struct {
		in              interface{}
		expectedErr     error
		expectErrString string
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

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			actualErr := Validate(tt.in)

			require.Equal(t, actualErr, tt.expectedErr)
			if actualErr != nil {
				require.Equal(t, tt.expectErrString, actualErr.Error())
			}
			_ = tt
		})
	}
}
