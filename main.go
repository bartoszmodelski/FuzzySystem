package main

import (
	"./fuzzy"
)

func main() {
	input := "tippingRulebase" + "\n" +
		"Rule 1 If driving is good and journey_time is short then tip will be big" + "\n" +
		"Rule 2 If driving is good or journey_time is short then tip will be moderate" + "\n" +
		"Rule 3 If driving is average or journey_time is medium then tip will be moderate" + "\n" +
		"Rule 4 If driving is bad and journey_time is long then tip will be small" + "\n" +
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
	fuzzy.NewSystem(input)
}
