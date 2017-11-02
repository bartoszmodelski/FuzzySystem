package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	rule                = iota
	variableDeclaration = iota
	variableAssignment  = iota
	empty               = iota
)

//
//}

func main() {
	input := "tippingRulebase" + "\n" +
		"Rule 1 If the driving is good and the journey_time is short then the tip will be big" + "\n" +
		"Rule 2 If the driving is good or the journey_time is short then the tip will be moderate" + "\n" +
		"Rule 3 If the driving is average or the journey_time is medium then the tip will be moderate" + "\n" +
		"Rule 4 If the driving is bad and the journey_time is long then the tip will be small" + "\n" +
		"driving " + "\n" +
		"bad 0 30 0 20 " + "\n" +
		"average 50 50 20 20" + "\n" +
		"good 80 100 20 0" + "\n" +
		"journey_time " + "\n" +
		"short 0 0 0 10" + "\n" +
		"medium 10 10 5 5" + "\n" +
		"long 20 20 10 0" + "\n" +
		"tip" + "\n" +
		"small 50 50 50 50" + "\n" +
		"moderate 100 100 50 50" + "\n" +
		"big 150 150 50 50" + "\n" +
		"journey_time = 9" + "\n" +
		"driving = 65"

	data, id := identifyNextChunk(input)
	offset := 0

	for id != empty {
		fmt.Println(strconv.Itoa(int(id)) + ": " + strings.TrimSpace(data))
		offset += len(data)
		data, id = identifyNextChunk(input[offset:])
	}
}

type FuzzyVariable struct {
	i int
}

type FuzzyRange struct {
	name                        string
	leftPeak, rightPeak         int
	leftBoundary, rightBoundary int
}

func identifyNextChunk(data string) (string, int8) {
	var id int8

	regexps := map[int8](*regexp.Regexp){
		variableDeclaration: regexp.MustCompile(`\A\s*\w+\s*(\s*\n\s*(\w+(\s+-?\d+){4}))+`),
		variableAssignment:  regexp.MustCompile(`\A\s*\w+\s*\=\s*\d+\s*`),
		rule:                regexp.MustCompile(`\A\s*\w+\s*(Rule \d+ If [^\n]* then [^\n]*\n)+`),
		empty:               regexp.MustCompile(`\s*`)}

	if regexps[variableDeclaration].MatchString(data) {
		id = variableDeclaration
	} else if regexps[variableAssignment].MatchString(data) {
		id = variableAssignment
	} else if regexps[rule].MatchString(data) {
		id = rule
	} else if regexps[empty].MatchString(data) {
		return "", empty
	} else {
		panic("Unrecognized data in configuration file: " + data)
	}

	chunk := regexps[id].FindString(data)
	return chunk, id
}

func NewFuzzyRange(data string) *FuzzyRange {
	fz := new(FuzzyRange)

	notWhitespace := regexp.MustCompile("[^\\s]+")
	parts := notWhitespace.FindAllString(data, -1)

	fz.name = parts[0]
	fz.leftPeak = parseInt(parts[1])
	fz.rightPeak = parseInt(parts[2])
	fz.leftBoundary = fz.leftPeak - parseInt(parts[3])
	fz.rightBoundary = fz.rightPeak + parseInt(parts[4])

	return fz
}

func parseInt(number string) int {
	value, err := strconv.Atoi(number)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return value
}
