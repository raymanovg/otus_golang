package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	RuleTypeLen    = "len"
	RuleTypeRegexp = "regexp"
	RuleTypeMax    = "max"
	RuleTypeMin    = "min"
	RuleTypeIn     = "in"
	RuleTypeNested = "nested"
)

var ErrInvalidValidationRule = errors.New("invalid validation rule")

type ValidationRules []ValidationRule

type ValidationRule struct {
	Type string
	Rule string
}

func (vr ValidationRule) String() string {
	if vr.Rule != "" {
		return fmt.Sprintf("%s:%s", vr.Type, vr.Rule)
	}
	return vr.Type
}

func RuleFromString(strRule string) (ValidationRule, error) {
	var vr ValidationRule
	parsedRule := strings.SplitN(strRule, ":", 2)
	switch parsedRule[0] {
	case RuleTypeLen, RuleTypeRegexp, RuleTypeMax, RuleTypeMin, RuleTypeIn:
		if len(parsedRule) < 2 || parsedRule[1] == "" {
			return vr, ErrInvalidValidationRule
		}
		vr.Type = parsedRule[0]
		vr.Rule = parsedRule[1]
		return vr, nil
	case RuleTypeNested:
		vr.Type = RuleTypeNested
		return vr, nil
	}
	return vr, ErrInvalidValidationRule
}

func RulesFromTag(tag reflect.StructTag) (ValidationRules, error) {
	var rules ValidationRules
	if tagValue, ok := tag.Lookup("validate"); ok {
		for _, strRule := range strings.Split(tagValue, "|") {
			vr, err := RuleFromString(strRule)
			if err != nil {
				return nil, err
			}
			rules = append(rules, vr)
		}
	}

	return rules, nil
}
