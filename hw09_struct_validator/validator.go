package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errStr string
	for _, err := range v {
		errStr += fmt.Sprintf("%s: %s", err.Field, err.Err.Error())
	}
	return errStr
}

func Validate(iv interface{}) error {
	v := reflect.ValueOf(iv)
	t := reflect.TypeOf(iv)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, but received %T", iv)
	}

	var verrors ValidationErrors
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)
		params, ok := ft.Tag.Lookup("validate")
		if !ok {
			continue
		}

		for _, p := range strings.Split(params, "|") {
			pt := strings.Split(p, ":")
			if len(pt) != 2 {
				verrors = append(verrors, ValidationError{ft.Name, fmt.Errorf("invalid validation params format %s", p)})
				continue
			}

			switch fv.Kind() {
			case reflect.String:
				if err := ValidateString(fv.String(), pt); err != nil {
					verrors = append(verrors, ValidationError{ft.Name, err})
				}
			case reflect.Int:
				if err := ValidateInt(fv.Int(), pt); err != nil {
					verrors = append(verrors, ValidationError{ft.Name, err})
				}
			case reflect.Slice:
				// TODO
			}
		}
	}

	return verrors
}

func ValidateString(strValue string, pt []string) error {
	switch pt[0] {
	case "len":
		i, err := strconv.Atoi(pt[1])
		if err != nil {
			return fmt.Errorf("invalid validation param %s:%s", pt[0], pt[1])
		}
		if i != len(strValue) {
			return fmt.Errorf("value len must be %s", pt[1])
		}
	case "in":
		var ok bool
		for _, expect := range strings.Split(pt[1], ",") {
			ok = strValue == strings.TrimSpace(expect)
		}
		if !ok {
			return fmt.Errorf("value must be in %s", pt[1])
		}
	case "regexp":
		expr := strings.Join(pt[1:], "")
		rgxp, err := regexp.Compile(expr)
		if err != nil {
			return fmt.Errorf("invalid validation param %s:%s: %w", pt[0], pt[1], err)
		}
		if !rgxp.Match([]byte(strValue)) {
			return fmt.Errorf("value is not matched by expression %s", expr)
		}
	default:
		return fmt.Errorf("unknow validation params %s", pt[0])
	}

	return nil
}

func ValidateInt(intValue int64, pt []string) error {
	switch pt[0] {
	case "min":
		i, err := strconv.Atoi(pt[1])
		if err != nil {
			return fmt.Errorf("invalid validation param %s:%s", pt[0], pt[1])
		}
		if int64(i) >= intValue {
			return fmt.Errorf("value min must be %s", pt[1])
		}
	case "max":
		i, err := strconv.Atoi(pt[1])
		if err != nil {
			return fmt.Errorf("invalid validation param %s:%s", pt[0], pt[1])
		}
		if int64(i) <= intValue {
			return fmt.Errorf("value min must be %s", pt[1])
		}
	case "in":
		var ok bool
		for _, expect := range strings.Split(pt[1], ",") {
			i, err := strconv.Atoi(expect)
			if err != nil {
				return fmt.Errorf("invalid validation param %s:%s", pt[0], pt[1])
			}
			ok = intValue == int64(i)
		}
		if !ok {
			return fmt.Errorf("value must be in %s", pt[1])
		}
	default:
		return fmt.Errorf("unknow validation params %s", pt[0])
	}

	return nil
}
