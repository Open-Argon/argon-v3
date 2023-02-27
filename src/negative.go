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

func parseNegative(code UNPARSEcode, index int, codeline []UNPARSEcode) (negative, bool, ArErr, int) {
	resp, worked, err, i := translateVal(UNPARSEcode{
		code:     strings.TrimSpace(code.code)[1:],
		realcode: code.realcode,
		line:     code.line,
		path:     code.path,
	}, index, codeline, false)
	return negative{
		VAL:  resp,
		line: code.line,
		code: code.realcode,
		path: code.path,
	}, worked, err, i
}
