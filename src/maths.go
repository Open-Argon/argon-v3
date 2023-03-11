package main

import (
	"github.com/wadey/go-rounding"
)

var maths = ArMap{
	"round": builtinFunc{"round", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return nil, ArErr{TYPE: "round", message: "round takes 1 argument",
				EXISTS: true}
		}
		precision := newNumber()
		if len(a) > 1 {
			switch x := a[1].(type) {
			case number:
				if !x.IsInt() {
					return nil, ArErr{TYPE: "TypeError", message: "Cannot round to '" + typeof(a[1]) + "'", EXISTS: true}
				}
				precision = x
			default:
				return nil, ArErr{TYPE: "TypeError", message: "Cannot round to '" + typeof(a[1]) + "'", EXISTS: true}
			}
		}

		switch x := a[0].(type) {
		case number:
			return rounding.Round(newNumber().Set(x), int(precision.Num().Int64()), rounding.HalfUp), ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot round '" + typeof(a[0]) + "'", EXISTS: true}
	}},
	"floor": builtinFunc{"floor", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return nil, ArErr{TYPE: "floor", message: "floor takes 1 argument",
				EXISTS: true}
		}
		switch x := a[0].(type) {
		case number:
			n := newNumber().Set(x)
			if n.Sign() < 0 {
				return rounding.Round(n, 0, rounding.Up), ArErr{}
			}
			return rounding.Round(n, 0, rounding.Down), ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot floor '" + typeof(a[0]) + "'", EXISTS: true}
	}},
	"ceil": builtinFunc{"ceil", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return nil, ArErr{TYPE: "ceil", message: "ceil takes 1 argument",
				EXISTS: true}
		}

		switch x := a[0].(type) {
		case number:
			n := newNumber().Set(x)
			if n.Sign() < 0 {
				return rounding.Round(n, 0, rounding.Down), ArErr{}
			}
			return rounding.Round(n, 0, rounding.Up), ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot ceil '" + typeof(a[0]) + "'", EXISTS: true}
	}},
	"sqrt": builtinFunc{"sqrt", ArgonSqrt},
	"ln":   builtinFunc{"ln", ArgonLn},
	"log":  builtinFunc{"log", ArgonLog},
	"logN": builtinFunc{"logN", ArgonLogN},
}
