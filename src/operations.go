package main

import (
	"fmt"
	"reflect"
	"strings"
)

var operations = [][]string{
	{
		"&&",
		" and ",
	}, {
		"||",
		" or ",
	}, {
		" not in ",
	}, {
		" in ",
	}, {
		"<=",
	}, {
		">=",
	}, {
		"<",
	}, {
		">",
	}, {
		"!=",
	}, {
		"==",
	}, {
		"+",
	}, {
		"-",
	}, {
		"*",
	}, {
		"%",
	}, {
		"//",
	}, {
		"/",
	}, {
		"^",
		"**",
	}}

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

func compareValues(o operationType, stack stack) (bool, ArErr) {
	if len(o.values) != 2 {
		return false, ArErr{
			"Runtime Error",
			"Invalid number of values for comparison",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
	resp, err := runVal(
		o.values[0],
		stack,
	)
	resp = classVal(resp)
	if err.EXISTS {
		return false, err
	}

	resp2, err := runVal(
		o.values[1],
		stack,
	)
	resp2 = classVal(resp2)
	if err.EXISTS {
		return false, err
	}
	switch o.operation {
	case 4:
		if isAnyNumber(resp) && isAnyNumber(resp2) {
			return resp.(number).Cmp(resp2.(number)) <= 0, ArErr{}
		}
		return false, ArErr{
			"Runtime Error",
			"Cannot compare type '" + typeof(resp) + "' with type '" + typeof(resp2) + "' with opperation '<='",
			o.line,
			o.path,
			o.code,
			true,
		}
	case 5:
		if isAnyNumber(resp) && isAnyNumber(resp2) {
			return resp.(number).Cmp(resp2.(number)) >= 0, ArErr{}
		}
		return false, ArErr{
			"Runtime Error",
			"Cannot compare type '" + typeof(resp) + "' with type '" + typeof(resp2) + "' with opperation '>='",
			o.line,
			o.path,
			o.code,
			true,
		}
	case 6:
		if isAnyNumber(resp) && isAnyNumber(resp2) {
			return resp.(number).Cmp(resp2.(number)) < 0, ArErr{}
		}
		return false, ArErr{
			"Runtime Error",
			"Cannot compare type '" + typeof(resp) + "' with type '" + typeof(resp2) + "' with opperation '<'",
			o.line,
			o.path,
			o.code,
			true,
		}
	case 7:
		if isAnyNumber(resp) && isAnyNumber(resp2) {
			return resp.(number).Cmp(resp2.(number)) > 0, ArErr{}
		}
		return false, ArErr{
			"Runtime Error",
			"Cannot compare type '" + typeof(resp) + "' with type '" + typeof(resp2) + "' with opperation '>'",
			o.line,
			o.path,
			o.code,
			true,
		}
	case 8:
		return !equals(resp, resp2), ArErr{}
	case 9:
		return equals(resp, resp2), ArErr{}
	default:
		return false, ArErr{
			"Runtime Error",
			"Invalid comparison operation",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
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
		if typeof(resp) == "number" {
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

func calcAdd(o operationType, stack stack) (any, ArErr) {

	resp, err := runVal(
		o.values[0],
		stack,
	)
	resp = classVal(resp)
	if err.EXISTS {
		return nil, err
	}
	var output any = resp
	if typeof(output) != "number" {
		output = anyToArgon(resp, false, true, 3, 0)
	}
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
		)
		resp = classVal(resp)
		if err.EXISTS {
			return nil, err
		}
		if typeof(output) == "number" && typeof(resp) == "string" {
			output = anyToArgon(output, false, true, 3, 0)
		}
		if typeof(output) == "number" {
			output = newNumber().Add(output.(number), resp.(number))
		} else {
			output = output.(string) + anyToArgon(resp, false, true, 3, 0)
		}

	}
	return output, ArErr{}
}

func calcAnd(o operationType, stack stack) (any, ArErr) {
	var output any = false
	for i := 0; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
		)
		resp = classVal(resp)
		if err.EXISTS {
			return nil, err
		}
		if !anyToBool(resp) {
			return resp, ArErr{}
		}
		output = resp
	}
	return output, ArErr{}
}

func calcOr(o operationType, stack stack) (any, ArErr) {
	var output any = false
	for i := 0; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
		)
		resp = classVal(resp)
		if err.EXISTS {
			return nil, err
		}
		if anyToBool(resp) {
			return resp, ArErr{}
		}
		output = resp
	}
	return output, ArErr{}
}

func stringInSlice(a any, list []any) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func calcIn(o operationType, stack stack) (bool, ArErr) {
	if len(o.values) != 2 {
		return false, ArErr{
			"Runtime Error",
			"Invalid number of arguments for 'not in'",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
	resp, err := runVal(
		o.values[0],
		stack,
	)
	resp = classVal(resp)
	if err.EXISTS {
		return false, err
	}

	resp2, err := runVal(
		o.values[1],
		stack,
	)
	resp2 = classVal(resp2)
	if err.EXISTS {
		return false, err
	}

	switch x := resp2.(type) {
	case string:
		check := anyToArgon(resp, false, true, 3, 0)
		return strings.Contains(x, check), ArErr{}
	case []any:
		return stringInSlice(resp, x), ArErr{}
	case map[any]any:
		_, ok := x[resp]
		return ok, ArErr{}
	default:
		return false, ArErr{
			"Runtime Error",
			"Cannot check if type '" + typeof(resp) + "' is in type '" + typeof(resp2) + "'",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
}

func equals(a any, b any) bool {
	if typeof(a) == "number" && typeof(b) == "number" {
		return a.(number).Cmp(b.(number)) == 0
	} else if typeof(a) == "string" || typeof(b) == "string" {
		return anyToArgon(a, false, true, 3, 0) == anyToArgon(b, false, true, 3, 0)
	}
	return reflect.DeepEqual(a, b)
}

func runOperation(o operationType, stack stack) (any, ArErr) {
	switch o.operation {
	case 0:
		return calcAnd(o, stack)
	case 1:
		return calcOr(o, stack)
	case 2:
		resp, err := calcIn(o, stack)
		resp = !resp
		return resp, err
	case 3:
		return calcIn(o, stack)
	case 4:
		return compareValues(o, stack)
	case 5:
		return compareValues(o, stack)
	case 6:
		return compareValues(o, stack)
	case 7:
		return compareValues(o, stack)
	case 8:
		return compareValues(o, stack)
	case 9:
		return compareValues(o, stack)
	case 10:
		return calcAdd(o, stack)
	case 11:
		return calcNegative(o, stack)

	}
	panic("Unknown operation: " + fmt.Sprint(o.operation))
}
