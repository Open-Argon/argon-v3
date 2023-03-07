package main

import (
	"fmt"
)

var operations = [][]string{
	{"-"},
	{"+"},
	{"/"},
	{"*"},
	{"%"},
	{"**", "^"},
	{"=="},
	{"!=", "≠"},
	{"<=", "≤"},
	{">=", "≥"},
	{"<"},
	{">"},
	{"&&", " and "},
	{"||", " or "},
}

type operationType struct {
	operation int
	values    []any
	line      int
	code      string
	path      string
}

func parseOperations(code UNPARSEcode, index int, codelines []UNPARSEcode) (operationType, bool, ArErr, int) {
	for i := 0; i < len(operations); i++ {
		values := []any{}
		current := 0
		for l := 0; l < len(code.code); l++ {
			for j := 0; j < len(operations[i]); j++ {
				if len(code.code[l:]) >= len(operations[i][j]) && code.code[l:l+len(operations[i][j])] == operations[i][j] {

					resp, success, _, respindex := translateVal(
						UNPARSEcode{
							code:     code.code[current:l],
							realcode: code.realcode,
							line:     code.line,
							path:     code.path,
						}, index, codelines, false)

					if success {
						index += respindex - 1
						values = append(values, resp)
						current = l + len(operations[i][j])
					}
				}
			}
		}
		if len(values) > 0 {
			resp, success, err, respindex := translateVal(
				UNPARSEcode{
					code:     code.code[current:],
					realcode: code.realcode,
					line:     code.line,
					path:     code.path,
				}, index, codelines, false)
			if success {
				index += respindex - 1
				values = append(values, resp)
				return operationType{
					i,
					values,
					code.line,
					code.realcode,
					code.path,
				}, true, err, index
			}
			return operationType{}, false, err, index
		}
	}
	return operationType{}, false, ArErr{}, index
}

func calcNegative(o operationType, stack stack) (number, ArErr) {

	resp, err := runVal(
		o.values[0],
		stack,
	)
	resp = classVal(resp)
	if err.EXISTS {
		return nil, err
	}
	if !isAnyNumber(resp) {
		return nil, ArErr{
			"Runtime Error",
			"Cannot subtract from type '" + typeof(resp) + "'",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
	output := resp.(number)
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
		)
		resp = classVal(resp)
		if err.EXISTS {
			return nil, err
		}
		if isAnyNumber(resp) {
			output = output.Sub(output, resp.(number))
		} else {
			return nil, ArErr{
				"Runtime Error",
				"Cannot subtract type '" + typeof(resp) + "'",
				o.line,
				o.path,
				o.code,
				true,
			}
		}
	}
	return output, ArErr{}
}

func runOperation(o operationType, stack stack) (any, ArErr) {
	switch o.operation {
	case 0:
		resp, err := calcNegative(o, stack)
		if err.EXISTS {
			return resp, err
		}
		return resp, ArErr{}

	}
	panic("Unknown operation: " + fmt.Sprint(o.operation))
}
