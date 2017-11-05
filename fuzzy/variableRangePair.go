package fuzzy

type fVariableRangePair struct {
	variable      *fVariable
	variableRange *fRange
}

func newFVariableRangePair(variable *fVariable, variableRange *fRange) *fVariableRangePair {
	rulePair := new(fVariableRangePair)
	rulePair.variableRange = variableRange
	rulePair.variable = variable

	return rulePair
}
