package main

import (
	"strings"
)

var commentCompile = makeRegex("(.)*#(.)*")

func isComment(code UNPARSEcode) bool {
	return commentCompile.MatchString(code.code)
}

func isBlank(code UNPARSEcode) bool {
	return strings.TrimSpace(code.code) == ""
}

func parseComment(code UNPARSEcode, index int, codelines []UNPARSEcode) (any, bool, ArErr) {
	split := strings.Split(code.code, "#")
	temp := []string{}
	for i := 0; i < len(split)-1; i++ {
		temp = append(temp, split[i])
		joined := strings.Join(temp, "#")
		resp, worked, _, _ := translateVal(UNPARSEcode{code: joined, realcode: code.realcode, line: code.line, path: code.path}, index, codelines, true)
		if worked {
			return resp, true, ArErr{}
		}
	}
	return nil, false, ArErr{"Syntax Error", "invalid comment", code.line, code.path, code.realcode, true}
}
