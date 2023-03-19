package main

import (
	"fmt"
	"math"
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
		totalindex := 1
		for l := 0; l < len(code.code); l++ {
			for j := 0; j < len(operations[i]); j++ {
				if len(code.code[l:]) >= len(operations[i][j]) && code.code[l:l+len(operations[i][j])] == operations[i][j] {

					resp, success, _, respindex := translateVal(
						UNPARSEcode{
							code:     code.code[current:l],
							realcode: code.realcode,
							line:     code.line,
							path:     code.path,
						}, index, codelines, 0)

					if success {
						totalindex += respindex - 1
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
				}, index, codelines, 0)
			if success {
				totalindex += respindex - 1
				values = append(values, resp)
				return operationType{
					i,
					values,
					code.line,
					code.realcode,
					code.path,
				}, true, err, totalindex
			}
			return operationType{}, false, err, totalindex
		}
	}
	return operationType{}, false, ArErr{}, 0
}

func compareValues(o operationType, stack stack, stacklevel int) (bool, ArErr) {
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
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
	if err.EXISTS {
		return false, err
	}

	resp2, err := runVal(
		o.values[1],
		stack,
		stacklevel+1,
	)
	resp2 = ArValidToAny(resp2)
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

func calcNegative(o operationType, stack stack, stacklevel int) (number, ArErr) {

	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
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
	output := newNumber().Set(resp.(number))
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		resp = ArValidToAny(resp)
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

func calcDivide(o operationType, stack stack, stacklevel int) (number, ArErr) {

	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
	if err.EXISTS {
		return nil, err
	}
	if !isAnyNumber(resp) {
		return nil, ArErr{
			"Runtime Error",
			"Cannot divide from type '" + typeof(resp) + "'",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
	output := newNumber().Set(resp.(number))
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		resp = ArValidToAny(resp)
		if err.EXISTS {
			return nil, err
		}
		if typeof(resp) == "number" {
			output = output.Quo(output, resp.(number))
		} else {
			return nil, ArErr{
				"Runtime Error",
				"Cannot divide type '" + typeof(resp) + "'",
				o.line,
				o.path,
				o.code,
				true,
			}
		}
	}
	return output, ArErr{}
}

func calcAdd(o operationType, stack stack, stacklevel int) (any, ArErr) {

	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
	if err.EXISTS {
		return nil, err
	}
	var output any = resp
	if typeof(output) != "number" {
		output = anyToArgon(resp, false, true, 3, 0, false, 0)
	} else {
		output = newNumber().Set(output.(number))
	}
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		resp = ArValidToAny(resp)
		if err.EXISTS {
			return nil, err
		}
		if typeof(output) == "number" && typeof(resp) == "string" {
			output = anyToArgon(output, false, true, 3, 0, false, 0)
		}
		if typeof(output) == "number" {
			output = output.(number).Add(output.(number), resp.(number))
		} else {
			output = output.(string) + anyToArgon(resp, false, true, 3, 0, false, 0)
		}

	}
	return AnyToArValid(output), ArErr{}
}

func calcMul(o operationType, stack stack, stacklevel int) (any, ArErr) {

	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
	if err.EXISTS {
		return nil, err
	}
	var output any = resp
	if typeof(output) != "number" {
		output = anyToArgon(resp, false, true, 3, 0, false, 0)
	} else {
		output = newNumber().Set(output.(number))
	}
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		resp = ArValidToAny(resp)
		if err.EXISTS {
			return nil, err
		}
		if typeof(output) == "number" && typeof(resp) == "string" {
			output = anyToArgon(output, false, true, 3, 0, false, 0)
		}
		if typeof(output) == "number" {
			output = output.(number).Mul(output.(number), resp.(number))
		} else if typeof(resp) == "number" {
			n, _ := resp.(number).Float64()
			output = strings.Repeat(output.(string), int(n))
		} else {
			return nil, ArErr{
				"Runtime Error",
				"Cannot multiply type '" + typeof(resp) + "'",
				o.line,
				o.path,
				o.code,
				true,
			}
		}
	}
	return output, ArErr{}
}

func calcAnd(o operationType, stack stack, stacklevel int) (any, ArErr) {
	var output any = false
	for i := 0; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		resp = ArValidToAny(resp)
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

func calcOr(o operationType, stack stack, stacklevel int) (any, ArErr) {
	var output any = false
	for i := 0; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		resp = ArValidToAny(resp)
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

func calcIn(o operationType, stack stack, stacklevel int) (bool, ArErr) {
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
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
	if err.EXISTS {
		return false, err
	}

	resp2, err := runVal(
		o.values[1],
		stack,
		stacklevel+1,
	)
	resp2 = ArValidToAny(resp2)
	if err.EXISTS {
		return false, err
	}

	switch x := resp2.(type) {
	case string:
		check := anyToArgon(resp, false, true, 3, 0, false, 0)
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
		return anyToArgon(a, false, false, 3, 0, false, 0) == anyToArgon(b, false, false, 3, 0, false, 0)
	}
	return reflect.DeepEqual(a, b)
}

func calcMod(o operationType, stack stack, stacklevel int) (number, ArErr) {
	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
	if err.EXISTS {
		return nil, err
	}
	if !isAnyNumber(resp) {
		return nil, ArErr{
			"Runtime Error",
			"Cannot calculate modulus from type '" + typeof(resp) + "'",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
	output := newNumber().Set(resp.(number))
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		resp = ArValidToAny(resp)
		if err.EXISTS {
			return nil, err
		}
		if typeof(resp) == "number" {
			n1, _ := output.Float64()
			n2, _ := resp.(number).Float64()
			output = newNumber().SetFloat64(math.Mod(n1, n2))
		} else {
			return nil, ArErr{
				"Runtime Error",
				"Cannot calculate modulus of type '" + typeof(resp) + "'",
				o.line,
				o.path,
				o.code,
				true,
			}
		}
	}
	return output, ArErr{}
}

func calcIntDiv(o operationType, stack stack, stacklevel int) (number, ArErr) {
	resp, err := calcDivide(o, stack, stacklevel+1)
	x, _ := resp.Float64()
	resp = newNumber().SetFloat64(math.Trunc(x))
	return resp, err
}

func calcPower(o operationType, stack stack, stacklevel int) (number, ArErr) {
	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
	if err.EXISTS {
		return nil, err
	}
	if typeof(resp) != "number" {
		return nil, ArErr{
			"Runtime Error",
			"Cannot calculate power of type '" + typeof(resp) + "'",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
	output := newNumber().Set(resp.(number))
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		resp = ArValidToAny(resp)
		if err.EXISTS {
			return nil, err
		}
		if typeof(resp) == "number" {
			n1, _ := output.Float64()
			n2, _ := resp.(number).Float64()
			output = newNumber().SetFloat64(math.Pow(n1, n2))
			if output == nil {
				output = infinity
			}
		} else {
			return nil, ArErr{
				"Runtime Error",
				"Cannot calculate power of type '" + typeof(resp) + "'",
				o.line,
				o.path,
				o.code,
				true,
			}
		}
	}
	return output, ArErr{}
}

func runOperation(o operationType, stack stack, stacklevel int) (any, ArErr) {
	switch o.operation {
	case 0:
		return calcAnd(o, stack, stacklevel+1)
	case 1:
		return calcOr(o, stack, stacklevel+1)
	case 2:
		resp, err := calcIn(o, stack, stacklevel+1)
		resp = !resp
		return resp, err
	case 3:
		return calcIn(o, stack, stacklevel+1)
	case 4:
		return compareValues(o, stack, stacklevel+1)
	case 5:
		return compareValues(o, stack, stacklevel+1)
	case 6:
		return compareValues(o, stack, stacklevel+1)
	case 7:
		return compareValues(o, stack, stacklevel+1)
	case 8:
		return compareValues(o, stack, stacklevel+1)
	case 9:
		return compareValues(o, stack, stacklevel+1)
	case 10:
		return calcAdd(o, stack, stacklevel+1)
	case 11:
		return calcNegative(o, stack, stacklevel+1)
	case 12:
		return calcMul(o, stack, stacklevel+1)
	case 13:
		return calcMod(o, stack, stacklevel+1)
	case 14:
		return calcIntDiv(o, stack, stacklevel+1)
	case 15:
		return calcDivide(o, stack, stacklevel+1)
	case 16:
		return calcPower(o, stack, stacklevel+1)

	}
	panic("Unknown operation: " + fmt.Sprint(o.operation))
}
