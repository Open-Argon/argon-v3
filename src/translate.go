package main

import (
	"fmt"
	"log"
)

// returns (translateNumber | translateString), success, error
func translateVal(code UNPARSEcode, index int, codelines []UNPARSEcode, isLine bool) (any, bool, string) {
	if isLine {
		if isComment(code) {
			return nil, true, ""
		}
	}

	if isNumber(code) {
		return parseNumber(code)
	} else if isString(code) {
		return parseString(code)
	}
	if isLine {
		return nil, false, "Syntax Error: invalid code on line " + fmt.Sprint(code.line) + ": " + code.code
	}
	return nil, false, ""
}

// returns [](translateNumber | translateString), error
func translate(codelines []UNPARSEcode) ([]any, string) {
	translated := []any{}
	for i, code := range codelines {
		val, _, err := translateVal(code, i, codelines, true)

		if err != "" {
			log.Fatal(err)
			return nil, err
		}
		if val == nil {
			continue
		}
		translated = append(translated, val)
	}
	return translated, ""
}
