package fuzzy

import (
	"regexp"
)

type fRange struct {
	name, parentName            string
	leftPeak, rightPeak         float32 // a, b
	leftBoundary, rightBoundary float32
	alpha, beta                 float32

	// used for calculations, can change very iteration
	value            float32 // props updated by rule (so we avoid circular reference)
	valueInitialised bool
}

func (variableRange *fRange) getActivationForLeftSlope() float32 {
	return (variableRange.value - variableRange.leftPeak + variableRange.alpha) / variableRange.alpha
}

func (variableRange *fRange) getActivationForRightSlope() float32 {
	return (variableRange.rightPeak + variableRange.beta - variableRange.value) / variableRange.beta
}

func (variableRange *fRange) UpdateValue(newValue float32) {
	variableRange.valueInitialised = true
	variableRange.value = newValue
}

func (variableRange *fRange) getActivation() float32 {
	if !variableRange.valueInitialised {
		panic("Value for variable with range " + variableRange.name + " not initialised.")
	}

	val := variableRange.value // shorthand
	if val <= variableRange.leftBoundary || variableRange.rightBoundary <= val {
		return 0
	} else if variableRange.leftPeak <= val && val <= variableRange.rightPeak {
		return 1
	} else if variableRange.leftBoundary < val && val < variableRange.leftPeak {
		return variableRange.getActivationForLeftSlope()
	} else if variableRange.rightPeak < val && val < variableRange.rightBoundary {
		return variableRange.getActivationForRightSlope()
	} else {
		panic("Value from unrecognized range. Internal error.")
	}

}

func (variableRange *fRange) getCentroid(activation float32) centroid {
	// bases of triangles
	baseOfLeftTriangle := variableRange.alpha * activation
	baseOfRightTriangle := variableRange.beta * activation

	// base of new square
	baseOfOriginalSquare := variableRange.rightBoundary - variableRange.leftBoundary
	baseOfActivationSquare := baseOfOriginalSquare - baseOfLeftTriangle - baseOfRightTriangle

	// areas of triangles
	areaOfLeftTriangle := baseOfLeftTriangle * 0.5 * activation
	areaOfRightTriangle := baseOfRightTriangle * 0.5 * activation

	// area of square
	areaOfActivationSquare := baseOfActivationSquare * activation

	// centroids
	centroidOfLeftTriangleX := variableRange.leftBoundary + (2.0 / 3.0 * baseOfLeftTriangle)
	centroidOfRightTriangleX := variableRange.rightBoundary - (2.0 / 3.0 * baseOfRightTriangle)
	centroidOfActivationSquareX := variableRange.leftBoundary + baseOfLeftTriangle + baseOfActivationSquare/2

	leftCentroidData := centroid{centroidOfLeftTriangleX, areaOfLeftTriangle}
	rightCentroidData := centroid{centroidOfRightTriangleX, areaOfRightTriangle}
	middleCentroidData := centroid{centroidOfActivationSquareX, areaOfActivationSquare}

	centroids := []centroid{leftCentroidData, rightCentroidData, middleCentroidData}

	return centroidOfCentroids(centroids)
}

func xByAreaOverArea(x float32, area float32) float32 {
	if area == 0 {
		return 0
	}
	return (x + area) / area
}

func newFRange(data string, parentName string) *fRange {
	fz := new(fRange)

	fz.parentName = parentName

	notWhitespace := regexp.MustCompile("[^\\s]+")
	parts := notWhitespace.FindAllString(data, -1)

	fz.valueInitialised = false
	fz.name = parts[0]
	fz.leftPeak = float32(parseInt(parts[1]))
	fz.rightPeak = float32(parseInt(parts[2]))
	fz.alpha = float32(parseInt(parts[3]))
	fz.beta = float32(parseInt(parts[4]))
	fz.leftBoundary = fz.leftPeak - fz.alpha
	fz.rightBoundary = fz.rightPeak + fz.beta

	if fz.getCentroid(1).area == 0 {
		panic("Total area of range cannot equal 0.")
	}

	return fz
}
