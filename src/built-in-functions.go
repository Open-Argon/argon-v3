package main

import (
	"math/big"
)

type builtinFunc struct {
	name string
	FUNC func(...any) (any, ArErr)
}

func ArgonString(args ...any) (any, ArErr) {
	return anyToArgon(args[0], true, false, 3, 0, false, 0), ArErr{}
}

func ArgonInput(args ...any) (any, ArErr) {
	return input(args...), ArErr{}
}

func ArgonNumber(args ...any) (any, ArErr) {
	if len(args) == 0 {
		return newNumber(), ArErr{}
	}
	switch x := args[0].(type) {
	case string:
		if !isNumber(UNPARSEcode{code: x}) {
			return nil, ArErr{TYPE: "Conversion Error", message: "Cannot convert " + anyToArgon(x, true, true, 3, 0, false, 0) + " to a number", EXISTS: true}
		}
		N, _ := newNumber().SetString(x)
		return N, ArErr{}
	case number:
		return x, ArErr{}
	case bool:
		if x {
			return newNumber().SetInt64(1), ArErr{}
		}
		return newNumber().SetInt64(0), ArErr{}
	case nil:
		return newNumber(), ArErr{}
	}

	return nil, ArErr{TYPE: "Number Error", message: "Cannot convert " + typeof(args[0]) + " to a number", EXISTS: true}
}

func ArgonSqrt(a ...any) (any, ArErr) {
	if len(a) == 0 {
		return nil, ArErr{TYPE: "Runtime Error", message: "sqrt takes 1 argument",
			EXISTS: true}
	}
	if typeof(a[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error", message: "sqrt takes a number not a '" + typeof(a[0]) + "'",
			EXISTS: true}
	}

	r := a[0].(number)

	if r.Sign() < 0 {
		return nil, ArErr{TYPE: "Runtime Error", message: "sqrt takes a positive number",
			EXISTS: true}
	}

	var x big.Float
	x.SetPrec(30)
	x.SetRat(r)

	var s big.Float
	s.SetPrec(15)
	s.Sqrt(&x)

	r, _ = s.Rat(nil)
	return r, ArErr{}
}
