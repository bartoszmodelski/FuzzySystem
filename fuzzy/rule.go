package fuzzy

import (
	"fmt"
	"strings"
)

const (
	and = iota
	or  = iota
)

type fRule struct {
	name       string
	sources    []*fVariableRangePair
	affected   []*fVariableRangePair
	connective int8 //AND or OR constants
}

func newFRule(input string, variables map[string]*fVariable) *fRule {
	// init empty structure representing rule
	rule := new(fRule)
	rule.affected = make([]*fVariableRangePair, 0, 0)
	rule.sources = make([]*fVariableRangePair, 0, 0)

	// filter out junk from input
	data := strings.NewReplacer(" is ", " ", " the ", " ", " will ", " ", " be ").Replace(input)

	// split into <rule name> + <rule>
	split := strings.Split(data, " If ")

	// save name
	rule.name = split[0]
	ruleContent := split[1]

	// split <rule> into <conditional> + <result>
	split = strings.Split(ruleContent, " then ")

	conditional := split[0]
	result := split[1]

	fmt.Println(conditional)
	fmt.Println(result)

	return rule
}

func findVariableAndRange(variables map[string]*fVariable, variableName string, rangeName string) *fVariableRangePair {
	for i := 0; i < len(variables[variableName].ranges); i++ {
		variableRange := variables[variableName].ranges[i]

		if variableRange.name == rangeName {
			return newFVariableRangePair(variables[variableName], variableRange)
		}
	}
	return nil
}

func parseRules(data string, variables map[string]*fVariable) []*fRule {
	rules := make([]*fRule, 0, 0)

	split := strings.Split(data, "\n")

	for i := 0; i < len(split); i++ {
		trimmed := strings.TrimSpace(split[i])
		if len(trimmed) > 0 {
			newFRule(trimmed, variables)
		}
	}

	return rules
}
