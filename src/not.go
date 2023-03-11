package main

import "strings"

var notCompiled = makeRegex(`( *)not (.|\n)+`)

type not struct {
	value any
	line  int
	code  string
	path  string
}

func isnot(code UNPARSEcode) bool {
	return notCompiled.MatchString(code.code)
}

func parseNot(code UNPARSEcode, index int, codelines []UNPARSEcode, isLine int) (any, bool, ArErr, int) {
	trimmed := strings.TrimSpace(code.code)
	trimmed = trimmed[4:]

	val, worked, err, step := translateVal(UNPARSEcode{
		code:     trimmed,
		realcode: code.realcode,
		line:     code.line,
		path:     code.path,
	}, index, codelines, isLine)
	return not{
		value: val,
		line:  code.line,
		code:  code.realcode,
		path:  code.path,
	}, worked, err, step
}

func runNot(n not, stack stack, stacklevel int) (bool, ArErr) {
	val, err := runVal(n.value, stack, stacklevel+1)
	boolean := !anyToBool(val)
	return boolean, err
}
