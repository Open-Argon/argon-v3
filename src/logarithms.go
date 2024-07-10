package main

import (
	"fmt"
	"math"
)

var N = newNumber().SetInt64(1e6)

func Ln(x number) number {
	output := newNumber()
	output.SetInt64(1)
	output.Quo(output, N)

	n1, _ := x.Float64()
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
	return Ln(x), ArErr{}
}

var __ln10 = Ln(newNumber().SetInt64(10))

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
	return Ln(x).Quo(Ln(x), __ln10), ArErr{}
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
