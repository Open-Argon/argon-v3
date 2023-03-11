package main

import (
	"fmt"
	"math/big"
	"strings"
)

var numberCompile = makeRegex("( *)(-)?((([0-9]+(\\.[0-9]+)?)|(\\.[0-9]+))(e((\\-|\\+)?([0-9]+(\\.[0-9]+)?)))?)( *)")
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

func isAnyNumber(x any) bool {
	_, ok := x.(number)
	return ok
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
func parseNumber(code UNPARSEcode) (number, bool, ArErr, int) {
	output, _ := newNumber().SetString(strings.TrimSpace(code.code))
	return output, true, ArErr{}, 1
}
