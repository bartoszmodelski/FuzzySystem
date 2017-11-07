package fuzzy

import (
	"strings"
)

const (
	none = iota
	and  = iota
	or   = iota
)

type fRule struct {
	name       string
	sources    []*fVariableRangePair
	affected   []*fVariableRangePair
	connective int8 //AND or OR constants
}

func (rule *fRule) getXandAreaOfCentroid() xAndArea {
	activation := rule.GetActivation()

	if len(rule.affected) > 1 {
		panic("Processing rules with multiple result not implemented yet.")
	}

	return rule.affected[0].variableRange.getXandAreaOfCenterOfMass(activation)
}

// GetActivation Returns rule's activation (based on values containe in fVariable objects)
func (rule *fRule) GetActivation() float32 {
	variablesCount := len(rule.sources)
	activations := make([]float32, variablesCount)

	for i := 0; i < variablesCount; i++ {
		activations[i] = rule.sources[i].variableRange.GetActivation()
	}

	return rule.minMax(activations)
}

func (rule *fRule) minMax(activations []float32) float32 {
	buffer := activations[0]
	if rule.connective == and { // min, find smallest activation and store in buffer
		for i := 1; i < len(activations); i++ {
			if buffer > activations[i] {
				buffer = activations[i]
			}
		}
	} else if rule.connective == or { // max, find biggest activation and store in buffer
		for i := 1; i < len(activations); i++ {
			if buffer < activations[i] {
				buffer = activations[i]
			}
		}
	} else {
		panic("Internal error: unrecognized value of connective (should be one of and/or constants).")
	}

	return buffer
}

func newFRule(input string, variables map[string]*fVariable) *fRule {
	// init empty structure representing rule
	rule := new(fRule)
	rule.affected = make([]*fVariableRangePair, 0, 0)
	rule.sources = make([]*fVariableRangePair, 0, 0)

	// filter out junk from input
	data := strings.NewReplacer(" is ", " ", " the ", " ", " will be ", " ").Replace(input)

	// split into <rule name> + <rule>
	split := strings.Split(data, " If ")

	// save name
	rule.name = split[0]
	ruleContent := split[1]

	// split <rule> into <conditional> + <result>
	split = strings.Split(ruleContent, " then ")

	conditional := split[0]
	result := split[1]

	rule.parseConditionalClause(conditional, variables)
	rule.parseResultClause(result, variables)
	return rule
}

func (rule *fRule) parseResultClause(result string, variables map[string]*fVariable) {
	resultParts := strings.Split(result, " ")

	for i := 0; i < len(resultParts); i++ {
		resultParts[i] = strings.TrimSpace(resultParts[i])
	}

	for i := 0; i < len(resultParts); i += 3 {
		rule.affected = append(rule.affected,
			findVariableAndRange(variables, resultParts[i], resultParts[i+1]))
	}
}

func (rule *fRule) parseConditionalClause(conditional string, variables map[string]*fVariable) {

	if strings.Contains(conditional, " or ") {
		rule.connective = or
	} else if strings.Contains(conditional, " and ") {
		rule.connective = and
	} else {
		rule.connective = none
	}

	conditionalParts := strings.Split(conditional, " ")

	for i := 0; i < len(conditionalParts); i++ {
		conditionalParts[i] = strings.TrimSpace(conditionalParts[i])
	}

	for i := 0; i < len(conditionalParts); i += 3 {

		rule.sources = append(rule.sources,
			findVariableAndRange(variables, conditionalParts[i], conditionalParts[i+1]))
	}
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

	for i := 1; i < len(split); i++ {
		trimmed := strings.TrimSpace(split[i])
		if len(trimmed) > 0 {
			rules = append(rules, newFRule(trimmed, variables))
		}
	}

	return rules
}
