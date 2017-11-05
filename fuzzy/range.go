package fuzzy

import "regexp"

type fRange struct {
	name                        string
	leftPeak, rightPeak         int
	leftBoundary, rightBoundary int
	alpha, beta                 int
}

func newFRange(data string) *fRange {
	fz := new(fRange)

	notWhitespace := regexp.MustCompile("[^\\s]+")
	parts := notWhitespace.FindAllString(data, -1)

	fz.name = parts[0]
	fz.leftPeak = parseInt(parts[1])
	fz.rightPeak = parseInt(parts[2])
	fz.alpha = parseInt(parts[3])
	fz.beta = parseInt(parts[4])
	fz.leftBoundary = fz.leftPeak - fz.alpha
	fz.rightBoundary = fz.rightPeak + fz.beta
	return fz
}
