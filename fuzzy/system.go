package fuzzy

import (
	"strings"
)

const (
	variableDeclaration = iota
	variableAssignment  = iota
	rule                = iota
	empty               = iota
)

// System - main object representing fuzzy system
type System struct {
	rules     []*fRule
	variables map[string]*fVariable
}

// NewSystem - construct new fuzzy system
func NewSystem(input string) *System {
	system := new(System)

	system.variables = make(map[string]*fVariable)
	system.rules = make([]*fRule, 0, 0)

	for i := 0; i < 3; i++ {
		offset := 0
		data, id := identifyNextChunk(input)
		for id != empty {
			if int(id) == i {
				switch id {
				case variableDeclaration:
					variable, key := newFVariable(data)
					system.variables[key] = variable
				case variableAssignment:
					split := strings.Split(data, "=")
					variableName := strings.TrimSpace(split[0])
					variableValue := parseInt(strings.TrimSpace(split[1]))
					system.variables[variableName].UpdateValue(float32(variableValue))
				case rule:
					buffer := parseRules(data, system.variables)
					system.rules = append(system.rules, buffer...)
				}

			}
			offset += len(data)
			data, id = identifyNextChunk(input[offset:])
		}
	}
	//spew.Dump(system.rules)

	return system
}

// Evaluate - perform one step in a system (one rules activation and assignments)
func (system *System) Evaluate() {
	topActivations := make(map[*fRange]xAndArea)

	for i := 0; i < len(system.rules); i++ {
		variableRange := system.rules[i].affected[0].variableRange
		value, present := topActivations[variableRange]

		activation := system.rules[i].getXandAreaOfCentroid()

		if present {
			if value.area < activation.area {
				//if current activation is stronger than
				topActivations[variableRange] = activation
			}
		} else {
			topActivations[variableRange] = activation
		}
	}

	topActivationsGrouped := make(map[*fVariable]([]*fRange))

	for key := range topActivations {
		parentVariable := system.variables[key.parentName]

		topActivationsGrouped[parentVariable] = append(topActivationsGrouped[parentVariable],
			key)
	}

}
