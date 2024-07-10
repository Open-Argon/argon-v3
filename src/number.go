package main

import (
	"fmt"
	"math/big"
	"strings"
)

var numberCompile = makeRegex("( *)(-)?(((([0-9]+(\\.[0-9]+)?)|(\\.[0-9]+))(e((\\-|\\+)?([0-9]+(\\.[0-9]+)?)))?)|([0-9]+/[0-9]+))( *)")
var binaryCompile = makeRegex("( *)(-)?(0b[10]+(.\\[10]+)?(e((\\-|\\+)?([0-9]+(\\.[0-9]+)?)))?)( *)")
var hexCompile = makeRegex("( *)(-)?(0x[a-fA-F0-9]+(\\.[a-fA-F0-9]+)?)( *)")
var octalCompile = makeRegex("( *)(-)?(0o[0-7]+(\\.[0-7]+)?(e((\\-|\\+)?([0-9]+(\\.[0-9]+)?)))?)( *)")

// a number type
type number = *big.Rat

// create a new number type
func newNumber() *big.Rat {
	return new(big.Rat)
}

func isNumber(code UNPARSEcode) bool {
	return numberCompile.MatchString(code.code) || binaryCompile.MatchString(code.code) || hexCompile.MatchString(code.code) || octalCompile.MatchString(code.code)
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

	x, _ := num.Float64()

	return fmt.Sprint(x)
}

// returns translateNumber, success, error
func parseNumber(code UNPARSEcode) (compiledNumber, bool, ArErr, int) {
	output, _ := new(big.Rat).SetString(strings.TrimSpace(code.code))
	if !output.IsInt() {
		return compiledNumber{output}, true, ArErr{}, 1
	}

	return compiledNumber{output.Num()}, true, ArErr{}, 1
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
		return floor(x).Num().Int64(), nil
	case *big.Int:
		return x.Int64(), nil
	}
	return 0, fmt.Errorf("object cannot be converted to int64")
}

func Int64ToNumber(num int64) ArObject {
	return Number(compiledNumber{big.NewInt(num)})
}
func CompareObjects(A ArObject, B ArObject) (ArObject, ArErr) {
	if X, ok := A.obj["__compare__"]; ok {
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

var _one = big.NewInt(1)
var _one_Rat = big.NewRat(1, 1)
var _one_Number ArObject
var _zero = big.NewInt(0)
var _zero_Rat = big.NewRat(0, 1)
var _zero_Number ArObject

func init() {
	_zero_Number = Number(compiledNumber{_zero})
	_one_Number = Number(compiledNumber{_one})
}

func Number(number compiledNumber) ArObject {
	// copy value to new number
	var value any = number.value
	val := ArObject{
		anymap{
			"__name__": "number",
		},
	}
	switch x := value.(type) {
	case *big.Rat:
		if x.IsInt() {
			value = x.Num()
		}
	case *big.Int:
	default:
		panic("invalid number type")
	}

	val.obj["__value__"] = value

	switch CurrentNumber := value.(type) {
	case *big.Int:
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
					output = append(output, "\x1b[34;5;240m")
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
					return Number(compiledNumber{CurrentNumber.Cmp(ReceivingNumber)}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).SetInt(CurrentNumber).Cmp(ReceivingNumber)}), ArErr{}
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
					return Number(compiledNumber{ReceivingNumber.Cmp(CurrentNumber)}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{ReceivingNumber.Cmp(new(big.Rat).SetInt(CurrentNumber))}), ArErr{}
				}
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
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
				return Number(compiledNumber{output}), ArErr{}
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
					return Number(compiledNumber{new(big.Int).Add(CurrentNumber, ReceivingNumber)}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Add(new(big.Rat).SetInt(CurrentNumber), ReceivingNumber)}), ArErr{}
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
					return Number(compiledNumber{new(big.Int).Sub(CurrentNumber, ReceivingNumber)}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Sub(new(big.Rat).SetInt(CurrentNumber), ReceivingNumber)}), ArErr{}
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
					return Number(compiledNumber{new(big.Int).Sub(ReceivingNumber, CurrentNumber)}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Sub(ReceivingNumber, new(big.Rat).SetInt(CurrentNumber))}), ArErr{}
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
					return Number(compiledNumber{new(big.Int).Mul(CurrentNumber, ReceivingNumber)}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Mul(new(big.Rat).SetInt(CurrentNumber), ReceivingNumber)}), ArErr{}
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
					return Number(compiledNumber{new(big.Int).Quo(CurrentNumber, ReceivingNumber)}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Quo(new(big.Rat).SetInt(CurrentNumber), ReceivingNumber)}), ArErr{}
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
					return Number(compiledNumber{new(big.Int).Quo(ReceivingNumber, CurrentNumber)}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Quo(ReceivingNumber, new(big.Rat).SetInt(CurrentNumber))}), ArErr{}
				}
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			},
		}

	case *big.Rat:
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
					output = append(output, "\x1b[34;5;240m")
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
					return Number(compiledNumber{CurrentNumber.Cmp(new(big.Rat).SetInt(ReceivingNumber))}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{CurrentNumber.Cmp(ReceivingNumber)}), ArErr{}
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
					return Number(compiledNumber{CurrentNumber.Cmp(new(big.Rat).SetInt(ReceivingNumber))}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{CurrentNumber.Cmp(ReceivingNumber)}), ArErr{}
				}
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
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
					return Number(compiledNumber{new(big.Rat).Add(CurrentNumber, new(big.Rat).SetInt(ReceivingNumber))}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Add(CurrentNumber, ReceivingNumber)}), ArErr{}
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
					return Number(compiledNumber{new(big.Rat).Sub(CurrentNumber, new(big.Rat).SetInt(ReceivingNumber))}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Sub(CurrentNumber, ReceivingNumber)}), ArErr{}
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
					return Number(compiledNumber{new(big.Rat).Sub(new(big.Rat).SetInt(ReceivingNumber), CurrentNumber)}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Sub(ReceivingNumber, CurrentNumber)}), ArErr{}
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
					return Number(compiledNumber{new(big.Rat).Mul(CurrentNumber, new(big.Rat).SetInt(ReceivingNumber))}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Mul(CurrentNumber, ReceivingNumber)}), ArErr{}
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
					return Number(compiledNumber{new(big.Rat).Mul(new(big.Rat).SetInt(ReceivingNumber), CurrentNumber)}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Mul(ReceivingNumber, CurrentNumber)}), ArErr{}
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
					return Number(compiledNumber{new(big.Rat).Quo(CurrentNumber, new(big.Rat).SetInt(ReceivingNumber))}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Quo(CurrentNumber, ReceivingNumber)}), ArErr{}
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
					return Number(compiledNumber{new(big.Rat).Quo(new(big.Rat).SetInt(ReceivingNumber), CurrentNumber)}), ArErr{}
				case *big.Rat:
					return Number(compiledNumber{new(big.Rat).Quo(ReceivingNumber, CurrentNumber)}), ArErr{}
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

	return val
}
