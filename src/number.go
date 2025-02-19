package main

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

var numberCompile = makeRegex("( *)(-)?(((([0-9]+(\\.[0-9]+)?)|(\\.[0-9]+))(e((\\-|\\+)?([0-9]+(\\.[0-9]+)?)))?)|([0-9]+/[0-9]+))( *)")
var binaryCompile = makeRegex("( *)(-)?(0b[10]+(.\\[10]+)?(e((\\-|\\+)?([0-9]+(\\.[0-9]+)?)))?)( *)")
var hexCompile = makeRegex("( *)(-)?(0x[a-fA-F0-9]+(\\.[a-fA-F0-9]+)?)( *)")
var octalCompile = makeRegex("( *)(-)?(0o[0-7]+(\\.[0-7]+)?(e((\\-|\\+)?([0-9]+(\\.[0-9]+)?)))?)( *)")

// a number type
// type number = *big.Rat

// create a new number type
// func newNumber() *big.Rat {
// 	return new(big.Rat)
// }

func isNumber(code UNPARSEcode) bool {
	return numberCompile.MatchString(code.code) || binaryCompile.MatchString(code.code) || hexCompile.MatchString(code.code) || octalCompile.MatchString(code.code)
}

func exponentBySquaring(base *big.Rat, exp *big.Int) *big.Rat {
	if exp.Cmp(_zero) == 0 {
		return _one_Rat
	}
	if exp.Cmp(_one) == 0 {
		return base
	}
	if exp.Bit(0) == 0 {
		return exponentBySquaring(new(big.Rat).Mul(base, base), new(big.Int).Div(exp, _two))
	}
	return new(big.Rat).Mul(base, exponentBySquaring(new(big.Rat).Mul(base, base), new(big.Int).Div(new(big.Int).Sub(exp, _one), _two)))
}

// converts a number type to a string
func numberToString(num *big.Rat, simplify bool) string {
	if simplify {
		divPI := new(big.Rat).Quo(num, PI_RAT)
		if divPI.Cmp(_one_Rat) == 0 {
			return "π"
		} else if divPI.Cmp(new(big.Rat).SetInt64(-1)) == 0 {
			return "-π"
		} else if divPI.Cmp(_zero_Rat) == 0 {
			return "0"
		} else if divPI.Denom().Cmp(new(big.Int).SetInt64(1000)) <= 0 {
			num := divPI.RatString()

			return fmt.Sprint(num, "π")
		}
	}

	float := new(big.Float).SetRat(num)

	return float.String()
}

var int64_max = new(big.Int).SetInt64(9223372036854775807)
var int64_min = new(big.Int).SetInt64(-9223372036854775808)

// returns translateNumber, success, error
func parseNumber(code UNPARSEcode) (compiledNumber, bool, ArErr, int) {
	output, _ := new(big.Rat).SetString(strings.TrimSpace(code.code))
	if !output.IsInt() {
		return compiledNumber{output}, true, ArErr{}, 1
	}

	output_big := output.Num()

	if output_big.Cmp(int64_max) > 0 || output_big.Cmp(int64_min) < 0 {
		return compiledNumber{output_big}, true, ArErr{}, 1
	}

	return compiledNumber{output_big.Int64()}, true, ArErr{}, 1
}

func isNumberInt64(num ArObject) bool {
	if x, ok := num.obj["__value__"]; ok {
		if _, ok := x.(int64); ok {
			return true
		}
	}
	return false
}

type compiledNumber = struct {
	value any
}

func isNumberInt(num ArObject) bool {
	value := num.obj["__value__"]
	switch x := value.(type) {
	case *big.Rat:
		return x.IsInt()
	case *big.Int:
		return true
	}
	return false
}

func numberToInt64(num ArObject) (int64, error) {
	value := num.obj["__value__"]
	switch x := value.(type) {
	case *big.Rat:
		if x.IsInt() {
			return x.Num().Int64(), nil
		}
		return 0, fmt.Errorf("ration number cannot be converted to int64")
	case *big.Int:
		return x.Int64(), nil
	case int64:
		return x, nil
	}
	return 0, fmt.Errorf("object cannot be converted to int64")
}

func Int64ToNumber(num int64) ArObject {
	return Number(num)
}
func CompareObjects(A ArObject, B ArObject) (ArObject, ArErr) {
	if X, ok := A.obj["__Compare__"]; ok {
		resp, err := runCall(call{
			Callable: X,
			Args:     []any{B},
		}, stack{}, 0)
		if !err.EXISTS {
			if resp, ok := resp.(ArObject); ok {
				return resp, ArErr{}
			}
		}
	} else if X, ok := B.obj["__PostCompare__"]; ok {
		resp, err := runCall(call{
			Callable: X,
			Args:     []any{A},
		}, stack{}, 0)
		if !err.EXISTS {
			if resp, ok := resp.(ArObject); ok {
				return resp, ArErr{}
			}
		}
	}
	return ArObject{}, ArErr{"Type Error", "cannot add " + typeof(A) + " and " + typeof(B), 0, "", "", true}
}

func AddObjects(A ArObject, B ArObject) (ArObject, ArErr) {
	if X, ok := A.obj["__Add__"]; ok {
		resp, err := runCall(call{
			Callable: X,
			Args:     []any{B},
		}, stack{}, 0)
		if !err.EXISTS {
			if resp, ok := resp.(ArObject); ok {
				return resp, ArErr{}
			}
		}
	} else if X, ok := B.obj["__PostAdd__"]; ok {
		resp, err := runCall(call{
			Callable: X,
			Args:     []any{A},
		}, stack{}, 0)
		if !err.EXISTS {
			if resp, ok := resp.(ArObject); ok {
				return resp, ArErr{}
			}
		}
	}
	return ArObject{}, ArErr{"Type Error", "cannot add " + typeof(A) + " and " + typeof(B), 0, "", "", true}
}

var _two = big.NewInt(2)
var _one = big.NewInt(1)
var _one_Rat = big.NewRat(1, 1)
var _one_Number ArObject
var _zero = big.NewInt(0)
var _zero_Rat = big.NewRat(0, 1)
var _zero_Number ArObject

func init() {
	_zero_Number = Number(0)
	_one_Number = Number(1)
}

func Number(value any) ArObject {
	val := ArObject{
		anymap{
			"__name__": "number",
		},
	}
	switch x := value.(type) {
	case *big.Rat:
		if x.IsInt() {
			value = x.Num()
			if value.(*big.Int).Cmp(int64_max) <= 0 && value.(*big.Int).Cmp(int64_min) >= 0 {
				value = value.(*big.Int).Int64()
			}
		}
	case *big.Int:
		if x.Cmp(int64_max) <= 0 && x.Cmp(int64_min) >= 0 {
			value = x.Int64()
		}
	case int:
		value = int64(x)
	case int64:
	case float64:
		value = new(big.Rat).SetFloat64(x)
	case float32:
		value = new(big.Rat).SetFloat64(float64(x))
	case string:
		value, _ = new(big.Rat).SetString(x)
		if value.(*big.Rat).IsInt() {
			value = value.(*big.Rat).Num()
		}
	default:
		panic("invalid number type")
	}

	val.obj["__value__"] = value

	val.obj["__LessThan__"] = builtinFunc{
		"__LessThan__",
		func(a ...any) (any, ArErr) {
			resp, err := val.obj["__Compare__"].(builtinFunc).FUNC(a...)
			if err.EXISTS {
				return nil, err
			}
			resp, _ = numberToInt64(resp.(ArObject))
			return resp.(int64) == -1, ArErr{}
		},
	}
	val.obj["__LessThanEqual__"] = builtinFunc{
		"__LessThanEqual__",
		func(a ...any) (any, ArErr) {
			resp, err := val.obj["__Compare__"].(builtinFunc).FUNC(a...)
			if err.EXISTS {
				return nil, err
			}
			resp, _ = numberToInt64(resp.(ArObject))
			return resp.(int64) != 1, ArErr{}
		},
	}
	val.obj["__GreaterThan__"] = builtinFunc{
		"__GreaterThan__",
		func(a ...any) (any, ArErr) {
			resp, err := val.obj["__Compare__"].(builtinFunc).FUNC(a...)
			if err.EXISTS {
				return nil, err
			}
			resp, _ = numberToInt64(resp.(ArObject))
			return resp.(int64) == 1, ArErr{}
		},
	}
	val.obj["__GreaterThanEqual__"] = builtinFunc{
		"__GreaterThanEqual__",
		func(a ...any) (any, ArErr) {
			resp, err := val.obj["__Compare__"].(builtinFunc).FUNC(a...)
			if err.EXISTS {
				return nil, err
			}
			resp, _ = numberToInt64(resp.(ArObject))
			return resp.(int64) != -1, ArErr{}
		},
	}
	val.obj["__NotEqual__"] = builtinFunc{
		"__NotEqual__",
		func(a ...any) (any, ArErr) {
			resp, err := val.obj["__Equal__"].(builtinFunc).FUNC(a...)
			if err.EXISTS {
				return nil, err
			}
			return !resp.(bool), ArErr{}
		},
	}
	switch CurrentNumber := value.(type) {
	case *big.Int:
		_BigInt_logic(val, CurrentNumber)
	case *big.Rat:
		_BigRat_logic(val, CurrentNumber)
	case int64:
		debugPrintln("int64", CurrentNumber)
		_int64_logic(val, CurrentNumber)
	}

	return val
}

func _BigInt_logic(val ArObject, CurrentNumber *big.Int) {
	val.obj["__string__"] = builtinFunc{
		"__string__",
		func(a ...any) (any, ArErr) {
			return fmt.Sprint(CurrentNumber), ArErr{}
		},
	}
	val.obj["__repr__"] = builtinFunc{
		"__repr__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "boolean" {
				return nil, ArErr{"Type Error", "expected boolean, got " + typeof(a[0]), 0, "", "", true}
			}
			coloured := a[0].(bool)
			output := []string{}
			if coloured {
				output = append(output, "\x1b[34;5;25m")
			}
			output = append(output, fmt.Sprint(CurrentNumber))
			if coloured {
				output = append(output, "\x1b[0m")
			}
			return strings.Join(output, ""), ArErr{}
		},
	}
	val.obj["__Compare__"] = builtinFunc{
		"__Compare__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(CurrentNumber.Cmp(ReceivingNumber)), ArErr{}
			case int64:
				return Number(CurrentNumber.Cmp(big.NewInt(ReceivingNumber))), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).SetInt(CurrentNumber).Cmp(ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostCompare__"] = builtinFunc{
		"__PostCompare__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(ReceivingNumber.Cmp(CurrentNumber)), ArErr{}
			case int64:
				return Number(big.NewInt(ReceivingNumber).Cmp(CurrentNumber)), ArErr{}
			case *big.Rat:
				return Number(ReceivingNumber.Cmp(new(big.Rat).SetInt(CurrentNumber))), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__Equal__"] = builtinFunc{
		"__Equal__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return CurrentNumber.Cmp(ReceivingNumber) == 0, ArErr{}
			case int64:
				return CurrentNumber.Cmp(big.NewInt(ReceivingNumber)) == 0, ArErr{}
			case *big.Rat:
				return new(big.Rat).SetInt(CurrentNumber).Cmp(ReceivingNumber) == 0, ArErr{}
			}
			return false, ArErr{}
		},
	}
	val.obj["__json__"] = builtinFunc{
		"__json__",
		val.obj["__string__"].(builtinFunc).FUNC,
	}
	val.obj["__Boolean__"] = builtinFunc{
		"__Boolean__",
		func(a ...any) (any, ArErr) {
			return CurrentNumber.Cmp(_zero) != 0, ArErr{}
		},
	}

	val.obj["__factorial__"] = builtinFunc{
		"__factorial__",
		func(a ...any) (any, ArErr) {
			if CurrentNumber.Cmp(_zero) < 0 {
				return nil, ArErr{"Runtime Error", "factorial of a negative number", 0, "", "", true}
			}
			output := new(big.Int).SetInt64(1)
			for i := new(big.Int).SetInt64(2); i.Cmp(CurrentNumber) <= 0; i.Add(i, _one) {
				output.Mul(output, i)
			}
			return Number(output), ArErr{}
		},
	}

	val.obj["__Add__"] = builtinFunc{
		"__Add__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Int).Add(CurrentNumber, ReceivingNumber)), ArErr{}
			case int64:
				return Number(new(big.Int).Add(CurrentNumber, big.NewInt(ReceivingNumber))), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Add(new(big.Rat).SetInt(CurrentNumber), ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostAdd__"] = builtinFunc{
		"__PostAdd__",
		val.obj["__Add__"].(builtinFunc).FUNC,
	}
	val.obj["__Subtract__"] = builtinFunc{
		"__Subtract__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Int).Sub(CurrentNumber, ReceivingNumber)), ArErr{}
			case int64:
				return Number(new(big.Int).Sub(CurrentNumber, big.NewInt(ReceivingNumber))), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Sub(new(big.Rat).SetInt(CurrentNumber), ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostSubtract__"] = builtinFunc{
		"__PostSubtract__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Int).Sub(ReceivingNumber, CurrentNumber)), ArErr{}
			case int64:
				return Number(new(big.Int).Sub(big.NewInt(ReceivingNumber), CurrentNumber)), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Sub(ReceivingNumber, new(big.Rat).SetInt(CurrentNumber))), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__Multiply__"] = builtinFunc{
		"__Multiply__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Int).Mul(CurrentNumber, ReceivingNumber)), ArErr{}
			case int64:
				return Number(new(big.Int).Mul(CurrentNumber, big.NewInt(ReceivingNumber))), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Mul(new(big.Rat).SetInt(CurrentNumber), ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostMultiply__"] = builtinFunc{
		"__PostMultiply__",
		val.obj["__Multiply__"].(builtinFunc).FUNC,
	}
	val.obj["__Divide__"] = builtinFunc{
		"__Divide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Quo(CurrentNumber, ReceivingNumber)), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Quo(CurrentNumber, big.NewInt(ReceivingNumber))), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(new(big.Rat).SetInt(CurrentNumber), ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostDivide__"] = builtinFunc{
		"__PostDivide_",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if CurrentNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Quo(ReceivingNumber, CurrentNumber)), ArErr{}
			case int64:
				if CurrentNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Quo(big.NewInt(ReceivingNumber), CurrentNumber)), ArErr{}
			case *big.Rat:
				if CurrentNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(ReceivingNumber, new(big.Rat).SetInt(CurrentNumber))), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__Modulo__"] = builtinFunc{
		"__Modulo__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Mod(CurrentNumber, ReceivingNumber)), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Mod(CurrentNumber, big.NewInt(ReceivingNumber))), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				x := new(big.Rat).Set(ReceivingNumber)
				output := new(big.Rat).SetInt(CurrentNumber)
				x.Quo(output, x)
				x = floor(x)
				x.Mul(x, ReceivingNumber)
				output.Sub(output, x)
				return Number(output), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostModulo__"] = builtinFunc{
		"__PostModulo__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if CurrentNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Mod(ReceivingNumber, CurrentNumber)), ArErr{}
			case int64:
				if CurrentNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Mod(big.NewInt(ReceivingNumber), CurrentNumber)), ArErr{}
			case *big.Rat:
				if CurrentNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				currentNumber_RAT := new(big.Rat).SetInt(CurrentNumber)
				x := new(big.Rat).Set(currentNumber_RAT)
				output := new(big.Rat).Set(ReceivingNumber)
				x.Quo(output, x)
				x = floor(x)
				x.Mul(x, currentNumber_RAT)
				output.Sub(output, x)
				return Number(output), ArErr{}
			}
			return false, ArErr{}
		},
	}
	val.obj["__IntDivide__"] = builtinFunc{
		"__IntDivide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Div(CurrentNumber, ReceivingNumber)), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Div(CurrentNumber, big.NewInt(ReceivingNumber))), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				x := new(big.Rat).Set(ReceivingNumber)
				output := new(big.Rat).SetInt(CurrentNumber)
				x.Quo(output, x)
				x = floor(x)
				return Number(x), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__PostIntDivide__"] = builtinFunc{
		"__PostIntDivide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if CurrentNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Div(ReceivingNumber, CurrentNumber)), ArErr{}
			case int64:
				if CurrentNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Div(big.NewInt(ReceivingNumber), CurrentNumber)), ArErr{}
			case *big.Rat:
				if CurrentNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				currentNumber_RAT := new(big.Rat).SetInt(CurrentNumber)
				x := new(big.Rat).Set(currentNumber_RAT)
				output := new(big.Rat).Set(ReceivingNumber)
				x.Quo(output, x)
				x = floor(x)
				return Number(x), ArErr{}
			}
			return false, ArErr{}
		},
	}
	val.obj["__Power__"] = builtinFunc{
		"__Power__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(_zero) == 0 {
					return Number(1), ArErr{}
				}
				output := new(big.Int).Set(CurrentNumber)
				output.Exp(output, new(big.Int).Abs(ReceivingNumber), nil)
				if ReceivingNumber.Cmp(_zero) < 0 {
					output = new(big.Int).Quo(_one, output)
				}
				return Number(output), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return Number(1), ArErr{}
				}
				output := new(big.Int).Set(CurrentNumber)
				output.Exp(output, new(big.Int).Abs(big.NewInt(ReceivingNumber)), nil)
				if ReceivingNumber < 0 {
					output = new(big.Int).Quo(_one, output)
				}
				return Number(output), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return Number(1), ArErr{}
				}
				exponent_numerator := new(big.Int).Abs(ReceivingNumber.Num())
				exponent_denominator := ReceivingNumber.Denom()
				output := new(big.Rat).SetInt(CurrentNumber)
				output_float, _ := output.Float64()
				// error if output_float is infinity
				if math.IsInf(output_float, 0) {
					return nil, ArErr{"Runtime Error", "number too large to perform rational exponential calculations on it.", 0, "", "", true}
				}
				exponent_denominator_float, _ := exponent_denominator.Float64()
				if math.IsInf(exponent_denominator_float, 0) {
					return nil, ArErr{"Runtime Error", "demominator too large to perform rational exponential calculations on it.", 0, "", "", true}
				}
				if exponent_denominator_float != 1 {
					output_float = math.Pow(output_float, 1/exponent_denominator_float)
				}
				output = new(big.Rat).SetFloat64(output_float)
				output = exponentBySquaring(output, exponent_numerator)
				if ReceivingNumber.Cmp(_zero_Rat) < 0 {
					output = new(big.Rat).Quo(_one_Rat, output)
				}
				return Number(output), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
}

func _BigRat_logic(val ArObject, CurrentNumber *big.Rat) {
	val.obj["__string__"] = builtinFunc{
		"__string__",
		func(a ...any) (any, ArErr) {
			return ArString(numberToString(CurrentNumber, false)), ArErr{}
		},
	}
	val.obj["__repr__"] = builtinFunc{
		"__repr__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "boolean" {
				return nil, ArErr{"Type Error", "expected boolean, got " + typeof(a[0]), 0, "", "", true}
			}
			coloured := a[0].(bool)
			output := []string{}
			if coloured {
				output = append(output, "\x1b[34;5;25m")
			}
			output = append(output, numberToString(CurrentNumber, true))
			if coloured {
				output = append(output, "\x1b[0m")
			}
			return ArString(strings.Join(output, "")), ArErr{}
		},
	}

	val.obj["__json__"] = builtinFunc{
		"__json__",
		val.obj["__string__"].(builtinFunc).FUNC,
	}

	val.obj["__Boolean__"] = builtinFunc{
		"__Boolean__",
		func(a ...any) (any, ArErr) {
			return CurrentNumber.Cmp(_zero_Rat) != 0, ArErr{}
		},
	}

	val.obj["__Compare__"] = builtinFunc{
		"__Compare__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(CurrentNumber.Cmp(new(big.Rat).SetInt(ReceivingNumber))), ArErr{}
			case int64:
				return Number(CurrentNumber.Cmp(new(big.Rat).SetInt64(ReceivingNumber))), ArErr{}
			case *big.Rat:
				return Number(CurrentNumber.Cmp(ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostCompare__"] = builtinFunc{
		"__PostCompare__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(CurrentNumber.Cmp(new(big.Rat).SetInt(ReceivingNumber))), ArErr{}
			case int64:
				return Number(CurrentNumber.Cmp(new(big.Rat).SetInt64(ReceivingNumber))), ArErr{}
			case *big.Rat:
				return Number(CurrentNumber.Cmp(ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__Equal__"] = builtinFunc{
		"__Equal__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return CurrentNumber.Cmp(new(big.Rat).SetInt(ReceivingNumber)) == 0, ArErr{}
			case int64:
				return CurrentNumber.Cmp(new(big.Rat).SetInt64(ReceivingNumber)) == 0, ArErr{}
			case *big.Rat:
				return CurrentNumber.Cmp(ReceivingNumber) == 0, ArErr{}
			}
			return false, ArErr{}
		},
	}

	val.obj["__Add__"] = builtinFunc{
		"__Add__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Rat).Add(CurrentNumber, new(big.Rat).SetInt(ReceivingNumber))), ArErr{}
			case int64:
				return Number(new(big.Rat).Add(CurrentNumber, new(big.Rat).SetInt64(ReceivingNumber))), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Add(CurrentNumber, ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostAdd__"] = builtinFunc{
		"__PostAdd__",
		val.obj["__Add__"].(builtinFunc).FUNC,
	}
	val.obj["__Subtract__"] = builtinFunc{
		"__Subtract__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Rat).Sub(CurrentNumber, new(big.Rat).SetInt(ReceivingNumber))), ArErr{}
			case int64:
				return Number(new(big.Rat).Sub(CurrentNumber, new(big.Rat).SetInt64(ReceivingNumber))), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Sub(CurrentNumber, ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostSubtract__"] = builtinFunc{
		"__PostSubtract__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Rat).Sub(new(big.Rat).SetInt(ReceivingNumber), CurrentNumber)), ArErr{}
			case int64:
				return Number(new(big.Rat).Sub(new(big.Rat).SetInt64(ReceivingNumber), CurrentNumber)), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Sub(ReceivingNumber, CurrentNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__Multiply__"] = builtinFunc{
		"__Multiply__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Rat).Mul(CurrentNumber, new(big.Rat).SetInt(ReceivingNumber))), ArErr{}
			case int64:
				return Number(new(big.Rat).Mul(CurrentNumber, new(big.Rat).SetInt64(ReceivingNumber))), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Mul(CurrentNumber, ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostMultiply__"] = builtinFunc{
		"__PostMultiply__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Rat).Mul(new(big.Rat).SetInt(ReceivingNumber), CurrentNumber)), ArErr{}
			case int64:
				return Number(new(big.Rat).Mul(new(big.Rat).SetInt64(ReceivingNumber), CurrentNumber)), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Mul(ReceivingNumber, CurrentNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__Divide__"] = builtinFunc{
		"__Divide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(CurrentNumber, new(big.Rat).SetInt(ReceivingNumber))), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(CurrentNumber, new(big.Rat).SetInt64(ReceivingNumber))), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(CurrentNumber, ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostDivide__"] = builtinFunc{
		"__PostDivide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if CurrentNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(new(big.Rat).SetInt(ReceivingNumber), CurrentNumber)), ArErr{}
			case int64:
				if CurrentNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(new(big.Rat).SetInt64(ReceivingNumber), CurrentNumber)), ArErr{}
			case *big.Rat:
				if CurrentNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(ReceivingNumber, CurrentNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__Modulo__"] = builtinFunc{
		"__Modulo__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				ReceivingNumber_RAT := new(big.Rat).SetInt(ReceivingNumber)
				x := new(big.Rat).Set(ReceivingNumber_RAT)
				output := new(big.Rat).Set(CurrentNumber)
				x.Quo(output, x)
				x = floor(x)
				x.Mul(x, ReceivingNumber_RAT)
				output.Sub(output, x)
				return Number(output), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				ReceivingNumber_RAT := new(big.Rat).SetInt64(ReceivingNumber)
				x := new(big.Rat).Set(ReceivingNumber_RAT)
				output := new(big.Rat).Set(CurrentNumber)
				x.Quo(output, x)
				x = floor(x)
				x.Mul(x, ReceivingNumber_RAT)
				output.Sub(output, x)
				return Number(output), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				x := new(big.Rat).Set(ReceivingNumber)
				output := new(big.Rat).Set(CurrentNumber)
				x.Quo(output, x)
				x = floor(x)
				x.Mul(x, ReceivingNumber)
				output.Sub(output, x)
				return Number(output), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostModulo__"] = builtinFunc{
		"__PostModulo__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if CurrentNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				ReceivingNumber_RAT := new(big.Rat).SetInt(ReceivingNumber)
				x := new(big.Rat).Set(ReceivingNumber_RAT)
				output := new(big.Rat).Set(CurrentNumber)
				x.Quo(output, x)
				x = floor(x)
				x.Mul(x, ReceivingNumber_RAT)
				output.Sub(output, x)
				return Number(output), ArErr{}
			case int64:
				if CurrentNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				ReceivingNumber_RAT := new(big.Rat).SetInt64(ReceivingNumber)
				x := new(big.Rat).Set(ReceivingNumber_RAT)
				output := new(big.Rat).Set(CurrentNumber)
				x.Quo(output, x)
				x = floor(x)
				x.Mul(x, ReceivingNumber_RAT)
				output.Sub(output, x)
				return Number(output), ArErr{}
			case *big.Rat:
				if CurrentNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				x := new(big.Rat).Set(ReceivingNumber)
				output := new(big.Rat).Set(CurrentNumber)
				x.Quo(output, x)
				x = floor(x)
				x.Mul(x, ReceivingNumber)
				output.Sub(output, x)
				return Number(output), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__IntDivide__"] = builtinFunc{
		"__IntDivide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(floor(new(big.Rat).Quo(CurrentNumber, new(big.Rat).SetInt(ReceivingNumber)))), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(floor(new(big.Rat).Quo(CurrentNumber, new(big.Rat).SetInt64(ReceivingNumber)))), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(floor(new(big.Rat).Quo(CurrentNumber, ReceivingNumber))), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__PostIntDivide__"] = builtinFunc{
		"__PostIntDivide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if CurrentNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(floor(new(big.Rat).Quo(new(big.Rat).SetInt(ReceivingNumber), CurrentNumber))), ArErr{}
			case int64:
				if CurrentNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(floor(new(big.Rat).Quo(new(big.Rat).SetInt64(ReceivingNumber), CurrentNumber))), ArErr{}
			case *big.Rat:
				if CurrentNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(floor(new(big.Rat).Quo(ReceivingNumber, CurrentNumber))), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__Power__"] = builtinFunc{
		"__Power__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(_zero) == 0 {
					return Number(1), ArErr{}
				}
				output := new(big.Rat).Set(CurrentNumber)
				output = exponentBySquaring(output, new(big.Int).Abs(ReceivingNumber))
				if ReceivingNumber.Cmp(_zero) < 0 {
					output = new(big.Rat).Quo(_one_Rat, output)
				}
				return Number(output), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return Number(1), ArErr{}
				}
				output := new(big.Rat).Set(CurrentNumber)
				output = exponentBySquaring(output, new(big.Int).Abs(big.NewInt(ReceivingNumber)))
				if ReceivingNumber < 0 {
					output = new(big.Rat).Quo(_one_Rat, output)
				}
				return Number(output), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return Number(1), ArErr{}
				}
				exponent_numerator := new(big.Int).Abs(ReceivingNumber.Num())
				exponent_denominator := ReceivingNumber.Denom()
				output := new(big.Rat).Set(CurrentNumber)
				output_float, _ := output.Float64()
				// error if output_float is infinity
				if math.IsInf(output_float, 0) {
					return nil, ArErr{"Runtime Error", "number too large to perform rational exponential calculations on it.", 0, "", "", true}
				}
				exponent_denominator_float, _ := exponent_denominator.Float64()
				if math.IsInf(exponent_denominator_float, 0) {
					return nil, ArErr{"Runtime Error", "demominator too large to perform rational exponential calculations on it.", 0, "", "", true}
				}
				if exponent_denominator_float != 1 {
					output_float = math.Pow(output_float, 1/exponent_denominator_float)
				}
				output = new(big.Rat).SetFloat64(output_float)
				output = exponentBySquaring(output, exponent_numerator)
				if ReceivingNumber.Cmp(_zero_Rat) < 0 {
					output = new(big.Rat).Quo(_one_Rat, output)
				}
				return Number(output), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__factorial__"] = builtinFunc{
		"__factorial__",
		func(a ...any) (any, ArErr) {
			return nil, ArErr{"Runtime Error", "factorial of a non-integer number", 0, "", "", true}
		},
	}
}

func _int64_logic(val ArObject, CurrentNumber int64) {
	val.obj["__string__"] = builtinFunc{
		"__string__",
		func(a ...any) (any, ArErr) {
			return ArString(fmt.Sprint(CurrentNumber)), ArErr{}
		},
	}
	val.obj["__repr__"] = builtinFunc{
		"__repr__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "boolean" {
				return nil, ArErr{"Type Error", "expected boolean, got " + typeof(a[0]), 0, "", "", true}
			}
			coloured := a[0].(bool)
			output := []string{}
			if coloured {
				output = append(output, "\x1b[34;5;25m")
			}
			output = append(output, fmt.Sprint(CurrentNumber))
			if coloured {
				output = append(output, "\x1b[0m")
			}
			return ArString(strings.Join(output, "")), ArErr{}
		},
	}

	val.obj["__json__"] = builtinFunc{
		"__json__",
		val.obj["__string__"].(builtinFunc).FUNC,
	}

	val.obj["__Boolean__"] = builtinFunc{
		"__Boolean__",
		func(a ...any) (any, ArErr) {
			return CurrentNumber != 0, ArErr{}
		},
	}

	val.obj["__Compare__"] = builtinFunc{
		"__Compare__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(big.NewInt(CurrentNumber).Cmp(ReceivingNumber)), ArErr{}
			case int64:
				if CurrentNumber < ReceivingNumber {
					return Number(-1), ArErr{}
				}
				if CurrentNumber > ReceivingNumber {
					return Number(1), ArErr{}
				}
				return Number(0), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).SetInt64(CurrentNumber).Cmp(ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostCompare__"] = builtinFunc{
		"__PostCompare__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(ReceivingNumber.Cmp(big.NewInt(CurrentNumber))), ArErr{}
			case int64:
				if ReceivingNumber < CurrentNumber {
					return Number(-1), ArErr{}
				}
				if ReceivingNumber > CurrentNumber {
					return Number(1), ArErr{}
				}
				return Number(0), ArErr{}
			case *big.Rat:
				return Number(ReceivingNumber.Cmp(new(big.Rat).SetInt64(CurrentNumber))), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__Equal__"] = builtinFunc{
		"__Equal__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return big.NewInt(CurrentNumber).Cmp(ReceivingNumber) == 0, ArErr{}
			case int64:
				return CurrentNumber == ReceivingNumber, ArErr{}
			case *big.Rat:
				return new(big.Rat).SetInt64(CurrentNumber).Cmp(ReceivingNumber) == 0, ArErr{}
			}
			return false, ArErr{}
		},
	}

	val.obj["__Add__"] = builtinFunc{
		"__Add__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Int).Add(big.NewInt(CurrentNumber), ReceivingNumber)), ArErr{}
			case int64:
				calc := CurrentNumber + ReceivingNumber
				// check for overflow
				if (ReceivingNumber > 0 && CurrentNumber > 0 && calc < 0) || (ReceivingNumber < 0 && CurrentNumber < 0 && calc > 0) {
					return Number(new(big.Int).Add(big.NewInt(CurrentNumber), big.NewInt(ReceivingNumber))), ArErr{}
				}
				return Number(calc), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Add(new(big.Rat).SetInt64(CurrentNumber), ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostAdd__"] = builtinFunc{
		"__PostAdd__",
		val.obj["__Add__"].(builtinFunc).FUNC,
	}
	val.obj["__Subtract__"] = builtinFunc{
		"__Subtract__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Int).Sub(big.NewInt(CurrentNumber), ReceivingNumber)), ArErr{}
			case int64:
				calc := CurrentNumber - ReceivingNumber
				// check for overflow
				if (ReceivingNumber > 0 && CurrentNumber < 0 && calc > 0) || (ReceivingNumber < 0 && CurrentNumber > 0 && calc < 0) {
					return Number(new(big.Int).Sub(big.NewInt(CurrentNumber), big.NewInt(ReceivingNumber))), ArErr{}
				}
				return Number(calc), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Sub(new(big.Rat).SetInt64(CurrentNumber), ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostSubtract__"] = builtinFunc{
		"__PostSubtract__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Int).Sub(ReceivingNumber, big.NewInt(CurrentNumber))), ArErr{}
			case int64:
				calc := ReceivingNumber - CurrentNumber
				// check for overflow
				if (ReceivingNumber < 0 && CurrentNumber > 0 && calc > 0) || (ReceivingNumber > 0 && CurrentNumber < 0 && calc < 0) {
					return Number(new(big.Int).Sub(big.NewInt(ReceivingNumber), big.NewInt(CurrentNumber))), ArErr{}
				}
				return Number(calc), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Sub(ReceivingNumber, new(big.Rat).SetInt64(CurrentNumber))), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__Multiply__"] = builtinFunc{
		"__Multiply__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Int).Mul(big.NewInt(CurrentNumber), ReceivingNumber)), ArErr{}
			case int64:
				calc := CurrentNumber * ReceivingNumber
				// check for overflow
				if ReceivingNumber != 0 && calc/ReceivingNumber != CurrentNumber {
					return Number(new(big.Int).Mul(big.NewInt(CurrentNumber), big.NewInt(ReceivingNumber))), ArErr{}
				}
				return Number(calc), ArErr{}
			case *big.Rat:
				return Number(new(big.Rat).Mul(new(big.Rat).SetInt64(CurrentNumber), ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostMultiply__"] = builtinFunc{
		"__PostMultiply__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				return Number(new(big.Int).Mul(ReceivingNumber, big.NewInt(CurrentNumber))), ArErr{}
			case int64:
				calc := ReceivingNumber * CurrentNumber
				// check for overflow
				if CurrentNumber != 0 && calc/CurrentNumber != ReceivingNumber {
					return Number(new(big.Int).Mul(big.NewInt(ReceivingNumber), big.NewInt(CurrentNumber))), ArErr{}
				}
			case *big.Rat:
				return Number(new(big.Rat).Mul(ReceivingNumber, new(big.Rat).SetInt64(CurrentNumber))), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__Divide__"] = builtinFunc{
		"__Divide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(big.NewInt(0)) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(new(big.Rat).SetInt64(CurrentNumber), new(big.Rat).SetInt(ReceivingNumber))), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(new(big.Rat).SetInt64(CurrentNumber), new(big.Rat).SetInt64(ReceivingNumber))), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(new(big.Rat).SetInt64(CurrentNumber), ReceivingNumber)), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostDivide__"] = builtinFunc{
		"__PostDivide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if CurrentNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Rat).Quo(new(big.Rat).SetInt(ReceivingNumber), new(big.Rat).SetInt64(CurrentNumber))), ArErr{}
			case int64:
				if CurrentNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
			case *big.Rat:
				if CurrentNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__Modulo__"] = builtinFunc{
		"__Modulo__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				CurrentNumber_BigInt := big.NewInt(CurrentNumber)
				return Number(new(big.Int).Mod(CurrentNumber_BigInt, ReceivingNumber)), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(CurrentNumber % ReceivingNumber), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				CurrentNumber_Rat := new(big.Rat).SetInt64(CurrentNumber)
				x := new(big.Rat).Set(ReceivingNumber)
				output := new(big.Rat).Set(CurrentNumber_Rat)
				x.Quo(output, x)
				x = floor(x)
				x.Mul(x, ReceivingNumber)
				output.Sub(output, x)
				return Number(output), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}
	val.obj["__PostModulo__"] = builtinFunc{
		"__PostModulo__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if CurrentNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Mod(ReceivingNumber, big.NewInt(CurrentNumber))), ArErr{}
			case int64:
				if CurrentNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(ReceivingNumber % CurrentNumber), ArErr{}
			case *big.Rat:
				if CurrentNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				x := new(big.Rat).Set(ReceivingNumber)
				output := new(big.Rat).SetInt64(CurrentNumber)
				x.Quo(output, x)
				x = floor(x)
				x.Mul(x, ReceivingNumber)
				output.Sub(output, x)
				return Number(output), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__IntDivide__"] = builtinFunc{
		"__IntDivide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(_zero) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Div(new(big.Int).SetInt64(CurrentNumber), ReceivingNumber)), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(CurrentNumber / ReceivingNumber), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(floor(new(big.Rat).Quo(new(big.Rat).SetInt64(CurrentNumber), ReceivingNumber))), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__PostIntDivide__"] = builtinFunc{
		"__PostIntDivide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if CurrentNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(new(big.Int).Div(ReceivingNumber, new(big.Int).SetInt64(CurrentNumber))), ArErr{}
			case int64:
				if CurrentNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(CurrentNumber / ReceivingNumber), ArErr{}
			case *big.Rat:
				if CurrentNumber == 0 {
					return nil, ArErr{"Runtime Error", "division by zero", 0, "", "", true}
				}
				return Number(floor(new(big.Rat).Quo(ReceivingNumber, new(big.Rat).SetInt64(CurrentNumber)))), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__Power__"] = builtinFunc{
		"__Power__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"Type Error", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			switch ReceivingNumber := a[0].(type) {
			case *big.Int:
				if ReceivingNumber.Cmp(_zero) == 0 {
					return Number(1), ArErr{}
				}
				output := new(big.Int).SetInt64(CurrentNumber)
				output.Exp(output, new(big.Int).Abs(ReceivingNumber), nil)
				if ReceivingNumber.Cmp(_zero) < 0 {
					output = new(big.Int).Quo(_one, output)
				}
				return Number(output), ArErr{}
			case int64:
				if ReceivingNumber == 0 {
					return Number(1), ArErr{}
				}
				output := new(big.Int).SetInt64(CurrentNumber)
				output.Exp(output, new(big.Int).Abs(big.NewInt(ReceivingNumber)), nil)
				if ReceivingNumber < 0 {
					output = new(big.Int).Quo(_one, output)
				}
				return Number(output), ArErr{}
			case *big.Rat:
				if ReceivingNumber.Cmp(_zero_Rat) == 0 {
					return Number(1), ArErr{}
				}
				exponent_numerator := new(big.Int).Abs(ReceivingNumber.Num())
				exponent_denominator := ReceivingNumber.Denom()
				output := new(big.Rat).SetInt64(CurrentNumber)
				output_float, _ := output.Float64()
				// error if output_float is infinity
				if math.IsInf(output_float, 0) {
					return nil, ArErr{"Runtime Error", "number too large to perform rational exponential calculations on it.", 0, "", "", true}
				}
				exponent_denominator_float, _ := exponent_denominator.Float64()
				if math.IsInf(exponent_denominator_float, 0) {
					return nil, ArErr{"Runtime Error", "demominator too large to perform rational exponential calculations on it.", 0, "", "", true}
				}
				if exponent_denominator_float != 1 {
					output_float = math.Pow(output_float, 1/exponent_denominator_float)
				}
				output = new(big.Rat).SetFloat64(output_float)
				output = exponentBySquaring(output, exponent_numerator)
				if ReceivingNumber.Cmp(_zero_Rat) < 0 {
					output = new(big.Rat).Quo(_one_Rat, output)
				}
				return Number(output), ArErr{}
			}
			return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
		},
	}

	val.obj["__factorial__"] = builtinFunc{
		"__factorial__",
		func(a ...any) (any, ArErr) {
			if CurrentNumber < 0 {
				return nil, ArErr{"Runtime Error", "factorial of a negative number", 0, "", "", true}
			}
			if CurrentNumber == 0 {
				return Number(1), ArErr{}
			}
			result := big.NewInt(1)
			for i := int64(1); i <= CurrentNumber; i++ {
				result.Mul(result, big.NewInt(i))
			}
			return Number(result), ArErr{}
		},
	}
}
