package fuzzy

import (
	"fmt"
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
	return system
}

// Evaluate - perform one step in a system (one rules activation and assignments)
func (system *System) Evaluate() {
	topActivationsGrouped := make(map[*fVariable](map[*fRange]centroid))

	for _, rule := range system.rules {
		topActivationsGrouped[rule.affected[0].variable] = make(map[*fRange]centroid)
	}

	for i := 0; i < len(system.rules); i++ {
		variableRange := system.rules[i].affected[0].variableRange
		variable := system.variables[variableRange.parentName]

		value, present := topActivationsGrouped[variable][variableRange]
		centroidData := system.rules[i].getCentroid()

		if centroidData.area != 0 {
			// if there is no centroid stored, store current
			if !present {
				topActivationsGrouped[variable][variableRange] = centroidData
			}

			// if current centroid represents stronger activation, overwrite stored
			if value.area < centroidData.area {
				topActivationsGrouped[variable][variableRange] = centroidData
			}
		}
	}

	// for every group, calculate final centroid and assign new variable
	for variable, topActivations := range topActivationsGrouped {
		centroids := make([]centroid, 0, 0)
		for _, centroidData := range topActivations {
			centroids = append(centroids, centroidData)
		}
		variable.value = centroidOfCentroids(centroids).x
		fmt.Printf("Assigned value %f to %s\n", variable.value, variable.name)

	}

}
