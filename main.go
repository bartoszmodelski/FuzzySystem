package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"./fuzzy"
)

func main() {
	system := fuzzy.NewSystem(loadFile("current.txt"))
	system.Evaluate()
}

func loadFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s", content)
}
