package main

import "strings"

var negativeCompile = makeRegex(`( *)-(.|\n)+( *)`)

type negative struct {
	VAL  any
	line int
	code string
	path string
}

func isNegative(code UNPARSEcode) bool {
	return negativeCompile.MatchString(code.code)
}

func parseNegative(code UNPARSEcode, index int, codeline []UNPARSEcode) (any, bool, ArErr, int) {
	trimmed := strings.TrimSpace(code.code)
	trimmednegative := strings.TrimLeft(trimmed, "-")
	difference := len(trimmed) - len(trimmednegative)
	resp, worked, err, i := translateVal(UNPARSEcode{
		code:     trimmednegative,
		realcode: code.realcode,
		line:     code.line,
		path:     code.path,
	}, index, codeline, 0)

	if difference%2 == 0 {
		return resp, worked, err, i
	}
	return negative{
		VAL:  resp,
		line: code.line,
		code: code.realcode,
		path: code.path,
	}, worked, err, i
}
