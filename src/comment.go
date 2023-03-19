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

func parseComment(code UNPARSEcode, index int, codelines []UNPARSEcode) (any, bool, ArErr, int) {
	split := strings.Split(code.code, "#")
	temp := []string{}
	step := 1
	for i := 0; i < len(split)-1; i++ {
		temp = append(temp, split[i])
		joined := strings.Join(temp, "#")
		if isBlank(UNPARSEcode{code: joined, realcode: code.realcode, line: code.line, path: code.path}) {
			return nil, true, ArErr{}, step
		}
		resp, worked, _, s := translateVal(UNPARSEcode{code: joined, realcode: code.realcode, line: code.line, path: code.path}, index, codelines, 2)
		step += s - 1
		if worked {
			return resp, true, ArErr{}, step
		}
	}
	return nil, false, ArErr{"Syntax Error", "invalid comment", code.line, code.path, code.realcode, true}, step
}
