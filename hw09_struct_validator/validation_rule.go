package main

import (
	"errors"
	"fmt"
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
	return fmt.Sprintf("%s:%s", vr.Type, vr.Rule)
}

func ValidationRuleFromStringRule(strRule string) (ValidationRule, error) {
	var vr ValidationRule
	parsedRule := strings.Split(strRule, ":")
	if len(parsedRule) < 2 {
		return vr, ErrInvalidValidationRule
	}
	switch parsedRule[0] {
	case RuleTypeLen, RuleTypeRegexp, RuleTypeMax, RuleTypeMin, RuleTypeIn:
		vr.Type = parsedRule[0]
		vr.Rule = strings.Join(parsedRule[1:], ":") // join if validation rule contains ":" symbol
		return vr, nil
	case RuleTypeNested:
		vr.Type = RuleTypeNested
	}
	return vr, ErrInvalidValidationRule
}

func ValidationRulesFromTagValue(tagValue string) (ValidationRules, error) {
	var rules ValidationRules
	for _, strRule := range strings.Split(tagValue, "|") {
		vr, err := ValidationRuleFromStringRule(strRule)
		if err != nil {
			return nil, err
		}
		rules = append(rules, vr)
	}
	return rules, nil
}
