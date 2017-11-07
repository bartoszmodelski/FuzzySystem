package fuzzy

import (
	"regexp"
	"strings"
)

type fVariable struct {
	name             string
	value            float32
	valueInitialised bool
	ranges           []*fRange
}

func (variable *fVariable) UpdateValue(value float32) {
	variable.valueInitialised = true
	variable.value = value

	for i := 0; i < len(variable.ranges); i++ {
		variable.ranges[i].UpdateValue(value)
	}
}

func newFVariable(data string) (*fVariable, string) {
	variable := new(fVariable)

	notNextLine := regexp.MustCompile("[^\\n]+")
	lines := notNextLine.FindAllString(strings.TrimSpace(data), -1)

	variable.name = strings.TrimSpace(lines[0])
	variable.valueInitialised = false
	variable.ranges = make([]*fRange, 0, len(lines))
	for i := 1; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if len(trimmed) > 0 {
			variable.ranges = append(variable.ranges, newFRange(trimmed, variable.name))
		}
	}
	return variable, variable.name
}
