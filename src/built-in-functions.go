package main

import (
	"math/big"
)

type builtinFunc struct {
	name string
	FUNC func(...any) (any, ArErr)
}

func ArgonString(args ...any) (any, ArErr) {
	if len(args) == 0 {
		return ArString(""), ArErr{}
	}
	args[0] = ArValidToAny(args[0])
	return ArString(anyToArgon(args[0], false, false, 3, 0, false, 0)), ArErr{}
}

func ArgonNumber(args ...any) (any, ArErr) {
	if len(args) == 0 {
		return _zero_Number, ArErr{}
	}
	args[0] = ArValidToAny(args[0])
	switch x := args[0].(type) {
	case string:
		if !isNumber(UNPARSEcode{code: x}) {
			return nil, ArErr{TYPE: "Conversion Error", message: "Cannot convert " + anyToArgon(x, true, true, 3, 0, false, 0) + " to a number", EXISTS: true}
		}
		N := Number(x)
		return N, ArErr{}
	case int64, *big.Int, *big.Rat:
		return Number(x), ArErr{}
	case bool:
		if x {
			return _one_Number, ArErr{}
		}
		return _zero_Number, ArErr{}
	case nil:
		return _zero_Number, ArErr{}
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
