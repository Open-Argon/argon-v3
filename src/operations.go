package main

import (
	"fmt"
	"reflect"
	"strings"
)

var operations = []string{
	"&&",
	"||",
	" not in ",
	" in ",
	"<=",
	">=",
	"<",
	">",
	"!=",
	"==",
	"+",
	"-",
	"*",
	"%",
	"//",
	"/",
	"^",
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
		split := strings.Split(code.code, operations[i])
		if len(split) < 2 {
			continue
		}
		var values []any
		lookingAt := 0
		totalStep := 1
		for j := 1; j < len(split); j++ {
			if split[j-1] == "" {
				continue
			}
			joined := strings.Join(split[lookingAt:j], operations[i])
			resp, success, err, respindex := translateVal(
				UNPARSEcode{
					code:     joined,
					realcode: code.realcode,
					line:     code.line,
					path:     code.path,
				}, index, codelines, 0)
			if !success || err.EXISTS {
				continue
			}
			values = append(values, resp)
			totalStep += respindex - 1
			lookingAt = j
		}
		if len(values) > 0 {
			resp, success, err, respindex := translateVal(
				UNPARSEcode{
					code:     strings.Join(split[lookingAt:], operations[i]),
					realcode: code.realcode,
					line:     code.line,
					path:     code.path,
				}, index, codelines, 0)
			if !success || err.EXISTS {
				continue
			}
			values = append(values, resp)
			totalStep += respindex - 1
			return operationType{
				operation: i,
				values:    values,
				line:      code.line,
				code:      code.realcode,
				path:      code.path,
			}, true, ArErr{}, totalStep
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
	if err.EXISTS {
		return false, err
	}

	resp2, err := runVal(
		o.values[1],
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return false, err
	}
	switch o.operation {
	case 4:
		if x, ok := resp.(ArObject); ok {
			if y, ok := x.obj["__LessThanEqual__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp2},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if err.EXISTS {
					return false, err
				}
				return anyToBool(val), ArErr{}
			}
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
		if x, ok := resp.(ArObject); ok {
			if y, ok := x.obj["__GreaterThanEqual__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp2},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if err.EXISTS {
					return false, err
				}
				return anyToBool(val), ArErr{}
			}
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
		if x, ok := resp.(ArObject); ok {
			if y, ok := x.obj["__LessThan__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp2},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if err.EXISTS {
					return false, err
				}
				return anyToBool(val), ArErr{}
			}
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
		if x, ok := resp.(ArObject); ok {
			if y, ok := x.obj["__GreaterThan__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp2},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if err.EXISTS {
					return false, err
				}
				return anyToBool(val), ArErr{}
			}
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
		return notequals(resp, resp2, o, stack, stacklevel+1)
	case 9:
		return equals(resp, resp2, o, stack, stacklevel+1)
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

func calcNegative(o operationType, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	output := resp
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		if err.EXISTS {
			return nil, err
		}
		if x, ok := output.(ArObject); ok {
			if y, ok := x.obj["__Subtract__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if err.EXISTS {
					return nil, err
				}
				output = val
				continue
			}
		}
		return nil, ArErr{
			"Runtime Error",
			"Cannot subtract type '" + typeof(resp) + "' from type '" + typeof(output) + "'",
			o.line,
			o.path,
			o.code,
			true,
		}

	}
	return output, ArErr{}
}

func calcDivide(o operationType, stack stack, stacklevel int) (any, ArErr) {

	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	output := resp
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)

		if err.EXISTS {
			return nil, err
		}
		if x, ok := output.(ArObject); ok {
			if y, ok := x.obj["__Divide__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if err.EXISTS {
					return nil, err
				}
				output = val
				continue
			}
		}
		return nil, ArErr{
			"Runtime Error",
			"Cannot divide type '" + typeof(resp) + "'",
			o.line,
			o.path,
			o.code,
			true,
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
	if err.EXISTS {
		return nil, err
	}
	var output any = resp
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		if err.EXISTS {
			return nil, err
		}
		if x, ok := output.(ArObject); ok {
			if y, ok := x.obj["__Add__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if !err.EXISTS {
					output = val
					continue
				}
			}
		}
		if x, ok := resp.(ArObject); ok {
			if y, ok := x.obj["__PostAdd__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{output},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if !err.EXISTS {
					output = val
					continue
				}
			}
		}
		return nil, ArErr{
			"Runtime Error",
			"Cannot add type '" + typeof(resp) + "' to type '" + typeof(output) + "'",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
	return output, ArErr{}
}

func calcMul(o operationType, stack stack, stacklevel int) (any, ArErr) {

	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	var output any = resp
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		if err.EXISTS {
			return nil, err
		}
		if x, ok := output.(ArObject); ok {
			if y, ok := x.obj["__Multiply__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if !err.EXISTS {
					output = val
					continue
				}
			}
		}
		if x, ok := resp.(ArObject); ok {
			if y, ok := x.obj["__PostMultiply__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{output},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if !err.EXISTS {
					output = val
					continue
				}
			}
		}
		return nil, ArErr{
			"Runtime Error",
			"Cannot multiply type '" + typeof(resp) + "'",
			o.line,
			o.path,
			o.code,
			true,
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

func InSlice(a any, list []any) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func calcNotIn(o operationType, stack stack, stacklevel int) (any, ArErr) {
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

	if err.EXISTS {
		return false, err
	}

	resp2, err := runVal(
		o.values[1],
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return false, err
	}
	if x, ok := resp2.(ArObject); ok {
		if y, ok := x.obj["__NotContains__"]; ok {
			return runCall(
				call{
					y,
					[]any{resp},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
		}
	}
	return false, ArErr{
		"Runtime Error",
		"Cannot check if type '" + typeof(resp) + "' is not in type '" + typeof(resp2) + "'",
		o.line,
		o.path,
		o.code,
		true,
	}
}

func calcIn(o operationType, stack stack, stacklevel int) (any, ArErr) {
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
	if err.EXISTS {
		return false, err
	}

	resp2, err := runVal(
		o.values[1],
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return false, err
	}
	if x, ok := resp2.(ArObject); ok {
		if y, ok := x.obj["__Contains__"]; ok {
			return runCall(
				call{
					y,
					[]any{resp},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
		}
	}
	return false, ArErr{
		"Runtime Error",
		"Cannot check if type '" + typeof(resp) + "' is in type '" + typeof(resp2) + "'",
		o.line,
		o.path,
		o.code,
		true,
	}
}
func notequals(a any, b any, o operationType, stack stack, stacklevel int) (bool, ArErr) {
	if x, ok := a.(ArObject); ok {
		if y, ok := x.obj["__NotEqual__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{b},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if !err.EXISTS {
				return anyToBool(val), ArErr{}
			}
		}
	}
	if x, ok := b.(ArObject); ok {
		if y, ok := x.obj["__NotEqual__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{a},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if err.EXISTS {
				return false, err
			}
			return anyToBool(val), ArErr{}
		}
	}
	return !reflect.DeepEqual(a, b), ArErr{}
}

func equals(a any, b any, o operationType, stack stack, stacklevel int) (bool, ArErr) {
	debugPrintln("equals", a, b)
	if x, ok := a.(ArObject); ok {
		if y, ok := x.obj["__Equal__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{b},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if !err.EXISTS {
				return anyToBool(val), ArErr{}
			}
		}
	}
	if x, ok := b.(ArObject); ok {
		if y, ok := x.obj["__Equal__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{a},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if err.EXISTS {
				return false, err
			}
			return anyToBool(val), ArErr{}
		}
	}
	return reflect.DeepEqual(a, b), ArErr{}
}

func calcMod(o operationType, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	output := resp
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		if err.EXISTS {
			return nil, err
		}
		if x, ok := output.(ArObject); ok {
			if y, ok := x.obj["__Modulo__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if err.EXISTS {
					return nil, err
				}
				output = val
				continue
			}
		}
		return nil, ArErr{
			"Runtime Error",
			"Cannot calculate modulus of type '" + typeof(resp) + "'",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
	return output, ArErr{}
}

func calcIntDiv(o operationType, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	output := resp
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)

		if err.EXISTS {
			return nil, err
		}
		if x, ok := output.(ArObject); ok {
			if y, ok := x.obj["__IntDivide__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if err.EXISTS {
					return nil, err
				}
				output = val
				continue
			}
		}
		return nil, ArErr{
			"Runtime Error",
			"Cannot int divide type '" + typeof(resp) + "'",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
	return output, ArErr{}
}

// func calcPower(o operationType, stack stack, stacklevel int) (number, ArErr) {
// 	resp, err := runVal(
// 		o.values[0],
// 		stack,
// 		stacklevel+1,
// 	)
//
// 	if err.EXISTS {
// 		return nil, err
// 	}
// 	if typeof(resp) != "number" {
// 		return nil, ArErr{
// 			"Runtime Error",
// 			"Cannot calculate power of type '" + typeof(resp) + "'",
// 			o.line,
// 			o.path,
// 			o.code,
// 			true,
// 		}
// 	}
// 	output := newNumber().Set(resp.(number))
// 	for i := 1; i < len(o.values); i++ {
// 		resp, err := runVal(
// 			o.values[i],
// 			stack,
// 			stacklevel+1,
// 		)
//
// 		if err.EXISTS {
// 			return nil, err
// 		}
// 		if typeof(resp) == "number" {
// 			n := newNumber().Set(resp.(number))
// 			if n.Cmp(newNumber().SetInt64(10)) <= 0 {
// 				toOut := newNumber().SetInt64(1)
// 				clone := newNumber().Set(output)
// 				nAbs := (abs(newNumber().Set(n)))
// 				j := newNumber()
// 				for ; j.Cmp(nAbs) < 0; j.Add(j, one) {
// 					toOut.Mul(toOut, clone)
// 				}

// 				nAbs.Sub(nAbs, j)
// 				if nAbs.Cmp(newNumber()) < 0 {
// 					j.Sub(j, one)
// 					n1, _ := toOut.Float64()
// 					n2, _ := nAbs.Float64()
// 					calculated := newNumber().SetFloat64(math.Pow(n1, n2))
// 					if calculated == nil {
// 						calculated = infinity
// 					}
// 					toOut.Mul(toOut, calculated)
// 				}
// 				if n.Cmp(newNumber()) < 0 {
// 					toOut.Quo(newNumber().SetInt64(1), toOut)
// 				}
// 				output.Set(toOut)
// 			} else if n.Cmp(newNumber().SetInt64(1)) != 0 {
// 				n1, _ := output.Float64()
// 				n2, _ := n.Float64()
// 				calculated := newNumber().SetFloat64(math.Pow(n1, n2))
// 				if calculated == nil {
// 					calculated = infinity
// 				}
// 				output.Mul(output, calculated)
// 			}

// 			/*
// 				n1, _ := output.Float64()
// 				n2, _ := resp.(number).Float64()
// 				output = newNumber().SetFloat64(math.Pow(n1, n2))
// 				if output == nil {
// 					output = infinity
// 				}
// 			*/
// 		} else {
// 			return nil, ArErr{
// 				"Runtime Error",
// 				"Cannot calculate power of type '" + typeof(resp) + "'",
// 				o.line,
// 				o.path,
// 				o.code,
// 				true,
// 			}
// 		}
// 	}
// 	return output, ArErr{}
// }

func calcPower(o operationType, stack stack, stacklevel int) (any, ArErr) {

	resp, err := runVal(
		o.values[0],
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	var output any = resp
	for i := 1; i < len(o.values); i++ {
		resp, err := runVal(
			o.values[i],
			stack,
			stacklevel+1,
		)
		if err.EXISTS {
			return nil, err
		}
		if x, ok := output.(ArObject); ok {
			if y, ok := x.obj["__Power__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if err.EXISTS {
					return nil, err
				}
				output = val
				continue
			}
		}
		return nil, ArErr{
			"Runtime Error",
			"Cannot power type '" + typeof(resp) + "' to type '" + typeof(output) + "'",
			o.line,
			o.path,
			o.code,
			true,
		}
	}
	return (output), ArErr{}
}

func runOperation(o operationType, stack stack, stacklevel int) (any, ArErr) {
	switch o.operation {
	case 0:
		return calcAnd(o, stack, stacklevel+1)
	case 1:
		return calcOr(o, stack, stacklevel+1)
	case 2:
		return calcNotIn(o, stack, stacklevel+1)
	case 3:
		return calcIn(o, stack, stacklevel+1)
	case 4, 5, 6, 7, 8, 9:
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
