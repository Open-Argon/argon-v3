package main

import "strings"

var arrayCompile = makeRegex(`( *)\[(.|\n)*\]( *)`)

type CreateArray struct {
	value ArArray
	line  int
	code  string
	path  string
}

func isArray(code UNPARSEcode) bool {
	return arrayCompile.MatchString(code.code)
}

func parseArray(code UNPARSEcode, index int, codelines []UNPARSEcode) (any, bool, ArErr, int) {
	trimmed := strings.TrimSpace(code.code)
	trimmed = trimmed[1 : len(trimmed)-1]
	arguments, worked, err := getValuesFromLetter(trimmed, ",", index, codelines, true)
	return CreateArray{
		value: arguments,
		line:  code.line,
		code:  code.realcode,
		path:  code.path,
	}, worked, err, 1
}

func runArray(a CreateArray, stack stack, stacklevel int) ([]any, ArErr) {
	var array ArArray
	for _, val := range a.value {
		val, err := runVal(val, stack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		array = append(array, val)
	}
	return array, ArErr{}
}
