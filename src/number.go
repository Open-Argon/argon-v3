package main

import (
	"fmt"
	"math/big"
)

// a number type
type number = *big.Rat

// create a number from two integers
var createNumber = big.NewRat

// create a new number type
func newNumber() *big.Rat {
	return new(big.Rat)
}

// converts a string into a number
func stringToNumber(str string) (*big.Rat, bool) {
	return newNumber().SetString(str)
}

// converts a number type to a string
func numberToString(num number, fraction bool) string {
	if fraction {
		return num.RatString()
	}
	x, _ := num.Float64()
	return fmt.Sprint(x)
}
