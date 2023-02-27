package main

import (
	"fmt"
	"math/big"
	"strings"
)

type translateNumber struct {
	number number
	code   string
	line   int
}

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

// converts a number type to a string
func numberToString(num number, fraction int) string {
	if fraction != 0 {
		str := num.RatString()
		if fraction == 1 {
			return str
		}
		split := strings.SplitN(str, "/", 2)
		if len(str) == 1 {
			return split[0]
		}
		numerator := split[0]
		denominator := split[1]

		super := []string{}
		for i := 0; i < len(numerator); i++ {
			super = append(super, superscript[numerator[i]])
		}
		sub := []string{}
		for i := 0; i < len(denominator); i++ {
			sub = append(sub, subscript[denominator[i]])
		}
		return strings.Join(super, "") + "/" + strings.Join(sub, "")
	}
	x, _ := num.Float64()
	return fmt.Sprint(x)
}

var superscript = map[byte]string{
	'0': "⁰",
	'1': "¹",
	'2': "²",
	'3': "³",
	'4': "⁴",
	'5': "⁵",
	'6': "⁶",
	'7': "⁷",
	'8': "⁸",
	'9': "⁹",
}

var subscript = map[byte]string{
	'0': "₀",
	'1': "₁",
	'2': "₂",
	'3': "₃",
	'4': "₄",
	'5': "₅",
	'6': "₆",
	'7': "₇",
	'8': "₈",
	'9': "₉",
}

// returns translateNumber, success, error
func parseNumber(code UNPARSEcode) (number, bool, ArErr, int) {
	output, _ := newNumber().SetString(strings.TrimSpace(code.code))
	return output, true, ArErr{}, 1
}
