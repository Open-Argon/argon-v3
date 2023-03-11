package main

import (
	"strings"
)

var bracketsCompile = makeRegex(`( *)\((.|\n)+\)( *)`)

type brackets struct {
	VAL  any
	line int
	code string
	path string
}

func isBrackets(code UNPARSEcode) bool {
	return bracketsCompile.MatchString(code.code)
}

func parseBrackets(code UNPARSEcode, index int, codeline []UNPARSEcode) (brackets, bool, ArErr, int) {
	trimmed := strings.TrimSpace(code.code)
	resp, worked, err, i := translateVal(UNPARSEcode{
		code:     trimmed[1 : len(trimmed)-1],
		realcode: code.realcode,
		line:     code.line,
		path:     code.path,
	}, index, codeline, 0)
	return brackets{
		VAL:  resp,
		line: code.line,
		code: code.realcode,
		path: code.path,
	}, worked, err, i
}
