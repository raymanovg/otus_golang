package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (vErrors ValidationErrors) Error() string {
	b := strings.Builder{}
	for _, vErr := range vErrors {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		if errors.As(vErr.Err, &ValidationErrors{}) {
			b.WriteString(fmt.Sprintf("%s.%s", vErr.Field, vErr.Err))
		} else {
			b.WriteString(fmt.Sprintf("%s: %s", vErr.Field, vErr.Err))
		}
	}
	return b.String()
}

func Validate(iv interface{}) error {
	v := reflect.ValueOf(iv)
	if v.Type().Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, but received %T", iv)
	}
	return ValidateStruct(v)
}

func ValidateStruct(v reflect.Value) error {
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, but received %T", v)
	}

	var vErrors ValidationErrors
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		if !isPublic(structField) {
			continue
		}

		validationRules, err := RulesFromTag(structField.Tag)
		if err != nil {
			return err
		}

		structValue := v.Field(i)

		var vErr error
		switch structValue.Kind() { //nolint:exhaustive
		case reflect.String:
			vErr = validateString(structValue, validationRules)
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			vErr = validateInt(structValue, validationRules)
		case reflect.Slice:
			vErr = validateSlice(structValue, validationRules)
		case reflect.Struct:
			if validationRules[0].Type == RuleTypeNested {
				vErr = ValidateStruct(structValue)
			}
		}

		if vErr != nil {
			vErrors = append(vErrors, ValidationError{structField.Name, vErr})
		}
	}

	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}

func isPublic(sf reflect.StructField) bool {
	return unicode.IsUpper([]rune(sf.Name)[0])
}

func validateSlice(v reflect.Value, rules ValidationRules) error {
	if v.Type().Kind() != reflect.Slice {
		return fmt.Errorf("expected a slice, but received %T", v.Interface())
	}
	if v.Len() == 0 {
		return nil
	}

	kind := v.Index(0).Kind()
	for i := 0; i < v.Len(); i++ {
		switch kind { //nolint:exhaustive
		case reflect.String:
			if err := validateString(v.Index(i), rules); err != nil {
				return err
			}
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			if err := validateInt(v.Index(i), rules); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateString(v reflect.Value, rules ValidationRules) error {
	if v.Type().Kind() != reflect.String {
		return fmt.Errorf("expected a struct, but received %T", v.Interface())
	}

	strVal := v.String()
	for _, vr := range rules {
		switch vr.Type {
		case RuleTypeLen:
			i, err := strconv.Atoi(vr.Rule)
			if err != nil {
				return fmt.Errorf("invalid validation rule %s", vr)
			}
			if i != len(strVal) {
				return fmt.Errorf("value len must be %s", vr.Rule)
			}
		case RuleTypeIn:
			var ok bool
			for _, expect := range strings.Split(vr.Rule, ",") {
				ok = strVal == strings.TrimSpace(expect)
				if ok {
					break
				}
			}
			if !ok {
				return fmt.Errorf("value must be in %s", vr.Rule)
			}
		case RuleTypeRegexp:
			rgxp, err := regexp.Compile(vr.Rule)
			if err != nil {
				return fmt.Errorf("invalid validation rule %s: %w", vr, err)
			}
			if !rgxp.Match([]byte(strVal)) {
				return fmt.Errorf("value is not matched by expression %s", vr.Rule)
			}
		default:
			return fmt.Errorf("unknown validation rule %s for string value", vr)
		}
	}

	return nil
}

func validateInt(v reflect.Value, rules ValidationRules) error {
	switch v.Type().Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
	default:
		return fmt.Errorf("expected a struct, but received %T", v.Interface())
	}

	intVal := v.Int()
	for _, vr := range rules {
		switch vr.Type {
		case RuleTypeMin:
			i, err := strconv.Atoi(vr.Rule)
			if err != nil {
				return fmt.Errorf("invalid validation rule '%s': %w", vr, err)
			}
			if int64(i) > intVal {
				return fmt.Errorf("value min must be %s", vr.Rule)
			}
		case RuleTypeMax:
			max, err := strconv.Atoi(vr.Rule)
			if err != nil {
				return fmt.Errorf("invalid validation rule '%s': %w", vr, err)
			}
			if int64(max) < intVal {
				return fmt.Errorf("value max must be %s", vr.Rule)
			}
		case RuleTypeIn:
			var ok bool
			for _, expect := range strings.Split(vr.Rule, ",") {
				i, err := strconv.Atoi(expect)
				if err != nil {
					return fmt.Errorf("invalid validation rule '%s': %w", vr, err)
				}
				if intVal == int64(i) {
					ok = true
					break
				}
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
