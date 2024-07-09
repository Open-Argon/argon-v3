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
func newNumber() number {
	return new(big.Rat)
}

func isNumber(code UNPARSEcode) bool {
	return numberCompile.MatchString(code.code) || binaryCompile.MatchString(code.code) || hexCompile.MatchString(code.code) || octalCompile.MatchString(code.code)
}

// converts a number type to a string
func numberToString(num number, simplify bool) string {
	if simplify {
		divPI := newNumber().Quo(num, PI)
		if divPI.Cmp(newNumber().SetInt64(1)) == 0 {
			return "π"
		} else if divPI.Cmp(newNumber().SetInt64(-1)) == 0 {
			return "-π"
		} else if divPI.Cmp(newNumber()) == 0 {
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

var _zero = big.NewInt(0)

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
		val.obj["__Boolean__"] = builtinFunc{
			"__Boolean__",
			func(a ...any) (any, ArErr) {
				return _zero, ArErr{}
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
				}
				return nil, ArErr{"Type Error", "expected number, got " + typeof(a[0]), 0, "", "", true}
			},
		}
	case *big.Rat:
		panic("not implemented")
	}

	return val
}
