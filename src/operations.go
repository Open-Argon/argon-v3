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

var one = newNumber().SetInt64(1)

type operationType struct {
	operation int
	value1    any
	value2    any
	line      int
	code      string
	path      string
}

func parseOperations(code UNPARSEcode, index int, codelines []UNPARSEcode) (operationType, bool, ArErr, int) {
	for i := 0; i < len(operations); i++ {
		for j := 0; j < len(operations[i]); j++ {
			split := strings.Split(code.code, operations[i][j])
			if len(split) <= 1 {
				continue
			}
			for k := 0; k < len(split)-1; k++ {
				if (len(strings.TrimSpace(split[k])) == 0 || len(strings.TrimSpace(split[k+1])) == 0) && operations[i][j] != "-" {
					break
				}
				val1, worked, err, step1 := translateVal(UNPARSEcode{
					code:     strings.Join(split[:k+1], operations[i][j]),
					realcode: code.realcode,
					line:     code.line,
					path:     code.path,
				}, index, codelines, 0)
				if !worked || err.EXISTS {
					if k == len(split)-1 {
						return operationType{}, false, err, 0
					} else {
						if len(strings.TrimSpace(split[k])) == 0 || len(strings.TrimSpace(split[k+1])) == 0 {
							break
						}
						continue
					}
				}

				val2, worked, err, step2 := translateVal(UNPARSEcode{
					code:     strings.Join(split[k+1:], operations[i][j]),
					realcode: code.realcode,
					line:     code.line,
					path:     code.path,
				}, index, codelines, 0)
				if !worked || err.EXISTS {
					if k == len(split)-1 {
						return operationType{}, false, err, 0
					} else {
						if len(strings.TrimSpace(split[k])) == 0 || len(strings.TrimSpace(split[k+1])) == 0 {
							break
						}
						continue
					}
				}
				return operationType{
					i,
					val1,
					val2,
					code.line,
					code.code,
					code.path,
				}, true, ArErr{}, step1 + step2 - 1

			}
		}
	}
	return operationType{}, false, ArErr{}, 0
}

func compareValues(o operationType, stack stack, stacklevel int) (bool, ArErr) {
	resp, err := runVal(
		o.value1,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return false, err
	}

	resp2, err := runVal(
		o.value2,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return false, err
	}
	switch o.operation {
	case 4:
		if isAnyNumber(resp) && isAnyNumber(resp2) {
			return resp.(number).Cmp(resp2.(number)) <= 0, ArErr{}
		} else if x, ok := resp.(ArObject); ok {
			if y, ok := x.obj["__LessThanEqual__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp2},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if !err.EXISTS {
					return anyToBool(val), ArErr{}
				}
			}
		}
		if x, ok := resp2.(ArObject); ok {
			if y, ok := x.obj["__GreaterThanEqual__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp},
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
		if isAnyNumber(resp) && isAnyNumber(resp2) {
			return resp.(number).Cmp(resp2.(number)) >= 0, ArErr{}
		} else if x, ok := resp.(ArObject); ok {
			if y, ok := x.obj["__GreaterThanEqual__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp2},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if !err.EXISTS {
					return anyToBool(val), ArErr{}
				}
			}
		}
		if x, ok := resp2.(ArObject); ok {
			if y, ok := x.obj["__LessThanEqual__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp},
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
		if isAnyNumber(resp) && isAnyNumber(resp2) {
			return resp.(number).Cmp(resp2.(number)) < 0, ArErr{}
		} else if x, ok := resp.(ArObject); ok {
			if y, ok := x.obj["__LessThan__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp2},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if !err.EXISTS {
					return anyToBool(val), ArErr{}
				}
			}
			if x, ok := resp2.(ArObject); ok {
				if y, ok := x.obj["__GreaterThan__"]; ok {
					val, err := runCall(
						call{
							y,
							[]any{resp},
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
		} else if x, ok := resp.(ArObject); ok {
			if y, ok := x.obj["__GreaterThan__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp2},
						o.code,
						o.line,
						o.path,
					}, stack, stacklevel+1)
				if !err.EXISTS {
					return anyToBool(val), ArErr{}
				}
			}
		}
		if x, ok := resp2.(ArObject); ok {
			if y, ok := x.obj["__LessThan__"]; ok {
				val, err := runCall(
					call{
						y,
						[]any{resp},
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
		o.value1,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	output := resp
	if isAnyNumber(resp) {
		output = newNumber().Set(resp.(number))
	}
	resp, err = runVal(
		o.value2,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	if typeof(output) == "number" && typeof(resp) == "number" {
		return output.(number).Sub(output.(number), resp.(number)), ArErr{}
	} else if x, ok := output.(ArObject); ok {
		if y, ok := x.obj["__Subtract__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{resp},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if !err.EXISTS {
				return val, ArErr{}
			}
		}
	}
	if x, ok := resp.(ArObject); ok {
		if y, ok := x.obj["__PostSubtract__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{output},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if err.EXISTS {
				return nil, err
			}
			return val, ArErr{}
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

func calcDivide(o operationType, stack stack, stacklevel int) (any, ArErr) {

	resp, err := runVal(
		o.value1,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	output := resp
	if isAnyNumber(resp) {
		output = newNumber().Set(resp.(number))
	}
	resp, err = runVal(
		o.value2,
		stack,
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
	if err.EXISTS {
		return nil, err
	}
	var outputErr ArErr = ArErr{
		"Runtime Error",
		"Cannot divide type '" + typeof(resp) + "'",
		o.line,
		o.path,
		o.code,
		true,
	}
	if typeof(resp) == "number" && typeof(output) == "number" {
		output = output.(number).Quo(output.(number), resp.(number))
		return output, ArErr{}
	} else if x, ok := output.(ArObject); ok {
		if y, ok := x.obj["__Divide__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{resp},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if !err.EXISTS {
				return val, ArErr{}
			}
			outputErr = err
		}
	}

	if x, ok := resp.(ArObject); ok {
		if y, ok := x.obj["__PostDivide__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{output},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if err.EXISTS {
				return nil, err
			}
			return val, ArErr{}
		}
	}
	return nil, outputErr
}

func calcAdd(o operationType, stack stack, stacklevel int) (any, ArErr) {

	resp, err := runVal(
		o.value1,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	var output any = resp
	if typeof(output) == "number" {
		output = newNumber().Set(output.(number))
	}
	resp, err = runVal(
		o.value2,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	if typeof(output) == "number" && typeof(resp) == "number" {
		output = output.(number).Add(output.(number), resp.(number))
		return output, ArErr{}
	} else if x, ok := output.(ArObject); ok {
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
				return output, ArErr{}
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
				return output, ArErr{}
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

func calcMul(o operationType, stack stack, stacklevel int) (any, ArErr) {

	resp, err := runVal(
		o.value1,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	var output any = resp
	if isAnyNumber(resp) {
		output = newNumber().Set(resp.(number))
	}
	resp, err = runVal(
		o.value2,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	if typeof(output) == "number" && typeof(resp) == "number" {
		output = output.(number).Mul(output.(number), resp.(number))
		return output, ArErr{}
	} else if x, ok := output.(ArObject); ok {
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
				return output, ArErr{}
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
			if err.EXISTS {
				return nil, err
			}
			return val, ArErr{}
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

func calcAnd(o operationType, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(
		o.value1,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	if !anyToBool(resp) {
		return resp, ArErr{}
	}
	resp, err = runVal(
		o.value2,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	if !anyToBool(resp) {
		return resp, ArErr{}
	}
	return resp, ArErr{}
}

func calcOr(o operationType, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(
		o.value1,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	if anyToBool(resp) {
		return resp, ArErr{}
	}
	resp, err = runVal(
		o.value2,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	if anyToBool(resp) {
		return resp, ArErr{}
	}
	return resp, ArErr{}
}

// InSlice checks if an element is present in a slice of any type.
// It returns true if the element is found, false otherwise.
func InSlice(a any, list []any) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// calcNotIn is a function that calculates the 'not in' operation between two values.
// It takes in an operationType 'o', a stack 'stack', and a stack level 'stacklevel'.
// It returns an 'any' value and an 'ArErr' error.
// The 'o' parameter contains information about the operation to be performed, including the values to be compared, the line of code, and the file path.
func calcNotIn(o operationType, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(
		o.value1,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return false, err
	}

	resp2, err := runVal(
		o.value2,
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

// calcIn is a function that calculates the 'in' operation between two values.
// It takes in an operationType 'o', a stack 'stack', and a stack level 'stacklevel'.
// It returns an 'any' value and an 'ArErr' error.
// The 'o' parameter contains information about the operation to be performed, including the values to be compared, the line of code, and the file path.
func calcIn(o operationType, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(
		o.value1,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return false, err
	}

	resp2, err := runVal(
		o.value2,
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
	if typeof(a) == "number" && typeof(b) == "number" {
		return a.(number).Cmp(b.(number)) != 0, ArErr{}
	} else if x, ok := a.(ArObject); ok {
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
				return !anyToBool(val), ArErr{}
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
	if typeof(a) == "number" && typeof(b) == "number" {
		return a.(number).Cmp(b.(number)) == 0, ArErr{}
	} else if x, ok := a.(ArObject); ok {
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
		if y, ok := x.obj["__GreaterThanEqual__"]; ok {
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
		o.value1,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	output := resp
	if isAnyNumber(resp) {
		output = newNumber().Set(resp.(number))
	}
	resp, err = runVal(
		o.value2,
		stack,
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
	if err.EXISTS {
		return nil, err
	}
	if typeof(resp) == "number" && typeof(output) == "number" {
		x := newNumber().Set(resp.(number))
		x.Quo(output.(number), x)
		x = floor(x)
		x.Mul(x, resp.(number))
		output.(number).Sub(output.(number), x)
		return output, ArErr{}
	} else if x, ok := output.(ArObject); ok {
		if y, ok := x.obj["__Modulo__"]; ok {
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
				return output, ArErr{}
			}
		}
	}

	if x, ok := resp.(ArObject); ok {
		if y, ok := x.obj["__PostModulo__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{output},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if err.EXISTS {
				return nil, err
			}
			return val, ArErr{}
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

func calcIntDiv(o operationType, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(
		o.value1,
		stack,
		stacklevel+1,
	)
	if err.EXISTS {
		return nil, err
	}
	output := resp
	if isAnyNumber(resp) {
		output = newNumber().Set(resp.(number))
	}
	resp, err = runVal(
		o.value2,
		stack,
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
	if err.EXISTS {
		return nil, err
	}
	if typeof(resp) == "number" && typeof(output) == "number" {
		output = output.(number).Quo(output.(number), resp.(number))
		return output, ArErr{}
	} else if x, ok := output.(ArObject); ok {
		if y, ok := x.obj["__IntDivide__"]; ok {
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
				return output, ArErr{}
			}
		}
	}
	if x, ok := resp.(ArObject); ok {
		if y, ok := x.obj["__PostIntDivide__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{output},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if err.EXISTS {
				return nil, err
			}
			return val, ArErr{}
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

func calcPower(o operationType, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(
		o.value1,
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
	resp, err = runVal(
		o.value2,
		stack,
		stacklevel+1,
	)
	resp = ArValidToAny(resp)
	if err.EXISTS {
		return nil, err
	}
	if typeof(resp) == "number" {
		n := newNumber().Set(resp.(number))
		if n.Cmp(newNumber().SetInt64(10)) <= 0 {
			toOut := newNumber().SetInt64(1)
			clone := newNumber().Set(output)
			nAbs := (abs(newNumber().Set(n)))
			j := newNumber()
			for ; j.Cmp(nAbs) < 0; j.Add(j, one) {
				toOut.Mul(toOut, clone)
			}

			nAbs.Sub(nAbs, j)
			if nAbs.Cmp(newNumber()) < 0 {
				j.Sub(j, one)
				n1, _ := toOut.Float64()
				n2, _ := nAbs.Float64()
				calculated := newNumber().SetFloat64(math.Pow(n1, n2))
				if calculated == nil {
					calculated = infinity
				}
				toOut.Mul(toOut, calculated)
			}
			if n.Cmp(newNumber()) < 0 {
				toOut.Quo(newNumber().SetInt64(1), toOut)
			}
			output.Set(toOut)
		} else if n.Cmp(newNumber().SetInt64(1)) != 0 {
			n1, _ := output.Float64()
			n2, _ := n.Float64()
			calculated := newNumber().SetFloat64(math.Pow(n1, n2))
			if calculated == nil {
				calculated = infinity
			}
			output.Mul(output, calculated)
		}
		return output, ArErr{}
	} else if x, ok := resp.(ArObject); ok {
		if y, ok := x.obj["__Power__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{output},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if !err.EXISTS {
				return val, ArErr{}
			}
		}
	}

	if x, ok := resp.(ArObject); ok {
		if y, ok := x.obj["__PostPower__"]; ok {
			val, err := runCall(
				call{
					y,
					[]any{output},
					o.code,
					o.line,
					o.path,
				}, stack, stacklevel+1)
			if err.EXISTS {
				return nil, err
			}
			return val, ArErr{}
		}
	}
	return nil, ArErr{
		"Runtime Error",
		"Cannot calculate power of type '" + typeof(resp) + "'",
		o.line,
		o.path,
		o.code,
		true,
	}
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
