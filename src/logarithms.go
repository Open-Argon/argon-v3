package main

import (
	"fmt"
	"math"
)

var N = Number(1e6)

func Ln(x ArObject) (any, ArErr) {
	var output any = Number(1)
	var err ArErr
	output, err = runOperation(
		operationType{
			operation: 15,
			values:    []any{x},
		},
		stack{},
		0,
	)
	if err.EXISTS {
		return nil, err
	}

	n1, _ := x_rational.Float64()
	n2, _ := output.Float64()
	output = newNumber().SetFloat64(math.Pow(n1, n2))
	output.Sub(output, newNumber().SetInt64(1))
	output.Mul(output, N)
	return output
}

func ArgonLn(a ...any) (any, ArErr) {
	if len(a) != 1 {
		return nil, ArErr{TYPE: "Runtime Error", message: "ln takes 1 argument, got " + fmt.Sprint(len(a)),
			EXISTS: true}
	}
	if typeof(a[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error", message: "ln takes a number not a '" + typeof(a[0]) + "'",
			EXISTS: true}
	}
	x := a[0].(number)
	if x.Sign() <= 0 {
		return nil, ArErr{TYPE: "Runtime Error", message: "ln takes a positive number",
			EXISTS: true}
	}
	return Ln(x)
	
}

var __ln10, _ = Ln(Number(10))

func ArgonLog(a ...any) (any, ArErr) {
	if len(a) != 1 {
		return nil, ArErr{TYPE: "Runtime Error", message: "log takes 1 argument, got " + fmt.Sprint(len(a)),
			EXISTS: true}
	}
	if typeof(a[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error", message: "log takes a number not a '" + typeof(a[0]) + "'",
			EXISTS: true}
	}
	x := a[0].(number)
	if x.Sign() <= 0 {
		return nil, ArErr{TYPE: "Runtime Error", message: "log takes a positive number",
			EXISTS: true}
	}
	return Ln(x).Quo(Ln(x), __ln10)
}

func ArgonLogN(a ...any) (any, ArErr) {
	if len(a) != 2 {
		return nil, ArErr{TYPE: "Runtime Error", message: "logN takes 2 argument, got " + fmt.Sprint(len(a)),
			EXISTS: true}
	}
	if typeof(a[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error", message: "logN takes a number not a '" + typeof(a[0]) + "'",
			EXISTS: true}
	}
	if typeof(a[1]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error", message: "logN takes a number not a '" + typeof(a[0]) + "'",
			EXISTS: true}
	}
	N := a[0].(number)
	if N.Sign() <= 0 {
		return nil, ArErr{TYPE: "Runtime Error", message: "logN takes a positive number",
			EXISTS: true}
	}
	x := a[1].(number)
	if x.Sign() <= 0 {
		return nil, ArErr{TYPE: "Runtime Error", message: "logN takes a positive number",
			EXISTS: true}
	}
	return Ln(x).Quo(Ln(x), Ln(N)), ArErr{}
}
