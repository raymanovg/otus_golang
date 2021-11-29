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
		errStr += fmt.Sprintf("%s: %s \n", err.Field, err.Err.Error())
	}
	return errStr
}

func Validate(iv interface{}) error {
	v := reflect.ValueOf(iv)
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, but received %T", iv)
	}

	var verrors ValidationErrors
	for i := 0; i < t.NumField(); i++ {
		structType := t.Field(i)
		structValue := v.Field(i)
		validateTagValue, ok := structType.Tag.Lookup("validate")
		if !ok {
			continue
		}

		validationRules, err := ValidationRulesFromTagValue(validateTagValue)
		if err != nil {
			verrors = append(verrors, ValidationError{structType.Name, err})
			continue
		}

		//nolint:exhaustive
		switch structValue.Kind() {
		case reflect.String:
			if err := ValidateString(structValue.String(), validationRules); err != nil {
				verrors = append(verrors, ValidationError{structType.Name, err})
			}
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			if err := ValidateInt(structValue.Int(), validationRules); err != nil {
				verrors = append(verrors, ValidationError{structType.Name, err})
			}
		case reflect.Slice:
			sliceValue := structValue.Index(0)
			for i := 0; i < structValue.Len(); i++ {
				switch sliceValue.Kind() {
				case reflect.String:
					if err := ValidateString(sliceValue.String(), validationRules); err != nil {
						verrors = append(verrors, ValidationError{structType.Name, err})
					}
				case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
					if err := ValidateInt(sliceValue.Int(), validationRules); err != nil {
						verrors = append(verrors, ValidationError{structType.Name, err})
					}
				}
			}
		case reflect.Struct:
			// TODO
		}
	}

	return verrors
}

func ValidateString(strValue string, rules ValidationRules) error {
	for _, vr := range rules {
		//nolint: exhaustive
		switch vr.Type {
		case RuleTypeLen:
			i, err := strconv.Atoi(vr.Rule)
			if err != nil {
				return fmt.Errorf("invalid validation rule %s", vr)
			}
			if i != len(strValue) {
				return fmt.Errorf("value len must be %s", vr.Rule)
			}
		case RuleTypeIn:
			var ok bool
			for _, expect := range strings.Split(vr.Rule, ",") {
				ok = strValue == strings.TrimSpace(expect)
			}
			if !ok {
				return fmt.Errorf("value must be in %s", vr.Rule)
			}
		case RuleTypeRegexp:
			rgxp, err := regexp.Compile(vr.Rule)
			if err != nil {
				return fmt.Errorf("invalid validation rule %s: %w", vr, err)
			}
			if !rgxp.Match([]byte(strValue)) {
				return fmt.Errorf("value is not matched by expression %s", vr.Rule)
			}
		default:
			return fmt.Errorf("unknow validation rule %s for string", vr)
		}
	}

	return nil
}

func ValidateInt(intValue int64, rules ValidationRules) error {
	for _, vr := range rules {
		//nolint: exhaustive
		switch vr.Type {
		case RuleTypeMin:
			i, err := strconv.Atoi(vr.Rule)
			if err != nil {
				return fmt.Errorf("invalid validation rule '%s': %w", vr, err)
			}
			if int64(i) > intValue {
				return fmt.Errorf("value min must be %s", vr.Rule)
			}
		case RuleTypeMax:
			max, err := strconv.Atoi(vr.Rule)
			if err != nil {
				return fmt.Errorf("invalid validation rule '%s': %w", vr, err)
			}
			if int64(max) < intValue {
				return fmt.Errorf("value max must be %s", vr.Rule)
			}
		case RuleTypeIn:
			var ok bool
			for _, expect := range strings.Split(vr.Rule, ",") {
				i, err := strconv.Atoi(expect)
				if err != nil {
					return fmt.Errorf("invalid validation rule '%s': %w", vr, err)
				}
				ok = intValue == int64(i)
			}
			if !ok {
				return fmt.Errorf("value must be in '%s'", vr.Rule)
			}
		default:
			return fmt.Errorf("unknown validation rule '%s' for int value", vr)
		}
	}

	return nil
}
