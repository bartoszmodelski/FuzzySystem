package fuzzy

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
)

func parseInt(number string) int {
	value, err := strconv.Atoi(number)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return value
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

	if id == rule {
		fmt.Println(chunk)
	}

	return chunk, id
}
