package main

import (
	"fmt"
	"math/big"
)

type builtinFunc struct {
	name string
	FUNC func(...any) (any, ArErr)
}

func ArgonLog(args ...any) (any, ArErr) {
	output := []any{}
	for i := 0; i < len(args); i++ {
		output = append(output, anyToArgon(args[i], false))
	}
	fmt.Println(output...)
	return nil, ArErr{}
}

func ArgonAdd(args ...any) (any, ArErr) {
	return reduce(func(x any, y any) any {
		return newNumber().Add(x.(number), y.(number))
	}, args), ArErr{}
}
func ArgonDiv(args ...any) (any, ArErr) {
	if len(args) == 0 {
		return nil, ArErr{TYPE: "Division Error", message: "Cannot divide nothing", EXISTS: true}
	}
	output := args[0].(number)
	for i := 1; i < len(args); i++ {
		if args[i].(number).Cmp(newNumber()) == 0 {
			return nil, ArErr{TYPE: "Division Error", message: "Cannot divide by zero", EXISTS: true}
		}
		output = newNumber().Quo(output, args[i].(number))
	}
	return output, ArErr{}
}

func ArgonMult(args ...any) (any, ArErr) {
	return reduce(func(x any, y any) any {
		return newNumber().Mul(y.(number), x.(number))
	}, args), ArErr{}
}

func ArgonInput(args ...any) (any, ArErr) {
	// allow a message to be passed in as an argument
	if len(args) > 0 {
		fmt.Print(anyToArgon(args[0], false))
	}
	var input string
	fmt.Scanln(&input)
	return input, ArErr{}
}

func ArgonNumber(args ...any) (any, ArErr) {
	if len(args) == 0 {
		return newNumber(), ArErr{}
	}
	switch x := args[0].(type) {
	case string:
		if !numberCompile.MatchString(x) {
			return nil, ArErr{TYPE: "Number Error", message: "Cannot convert type '" + typeof(x) + "' to a number", EXISTS: true}
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
		return nil, ArErr{TYPE: "sqrt", message: "sqrt takes 1 argument",
			EXISTS: true}
	}
	r := a[0].(number)

	if r.Sign() < 0 {
		return nil, ArErr{TYPE: "sqrt", message: "sqrt takes a positive number",
			EXISTS: true}
	}

	var x big.Float
	x.SetPrec(30) // I didn't figure out the 'Prec' part correctly, read the docs more carefully than I did and experiement
	x.SetRat(r)

	var s big.Float
	s.SetPrec(15)
	s.Sqrt(&x)

	r, _ = s.Rat(nil)
	return r, ArErr{}
}
