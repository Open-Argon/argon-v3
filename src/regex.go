package main

import (
	"log"
	"regexp"
)

func makeRegex(str string) *regexp.Regexp {
	Compile, err := regexp.Compile("^(" + str + ")$")
	if err != nil {
		log.Fatal(err)
	}
	return Compile
}
