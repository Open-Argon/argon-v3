package main

import (
	"fmt"
	"math"
)

var PIFloatInaccuracy number = newNumber()

func init() {
	PIFloatInaccuracy.SetFloat64(math.Asin(1) * 2)
}

var ArSin = builtinFunc{"sin", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("sin expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("sin expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := newNumber().Set(args[0].(number))
	num.Quo(num, PI)
	num.Mul(num, PIFloatInaccuracy)
	n, _ := num.Float64()
	outputnum := newNumber().SetFloat64(math.Sin(n))
	return outputnum, ArErr{}
}}
var ArArcsin = builtinFunc{"arcsin", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arcsin expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arcsin expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := args[0].(number)
	n, _ := num.Float64()
	if n < -1 || n > 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arcsin expected number between -1 and 1, got %s", anyToArgon(n, true, true, 3, 0, false, 0)),
			EXISTS:  true,
		}
	}
	outputnum := newNumber().SetFloat64(math.Asin(n))
	outputnum.Quo(outputnum, PIFloatInaccuracy)
	outputnum.Mul(outputnum, PI)
	return outputnum, ArErr{}
}}

var ArCos = builtinFunc{"cos", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("cos expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("cos expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := newNumber().Set(args[0].(number))
	num.Quo(num, PI)
	num.Mul(num, PIFloatInaccuracy)
	n, _ := num.Float64()
	outputnum := newNumber().SetFloat64(math.Cos(n))
	return outputnum, ArErr{}
}}
var ArArccos = builtinFunc{"arccos", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arccos expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arccos expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := args[0].(number)
	n, _ := num.Float64()
	if n < -1 || n > 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arccos expected number between -1 and 1, got %s", anyToArgon(n, true, true, 3, 0, false, 0)),
			EXISTS:  true,
		}
	}
	outputnum := newNumber().SetFloat64(math.Acos(n))
	outputnum.Quo(outputnum, PIFloatInaccuracy)
	outputnum.Mul(outputnum, PI)
	return outputnum, ArErr{}
}}

var ArTan = builtinFunc{"tan", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("tan expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("tan expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := newNumber().Set(args[0].(number))
	num.Quo(num, PI)
	num.Mul(num, PIFloatInaccuracy)
	n, _ := num.Float64()
	outputnum := newNumber().SetFloat64(math.Tan(n))
	return outputnum, ArErr{}
}}
var ArArctan = builtinFunc{"arctan", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arctan expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arctan expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := args[0].(number)
	n, _ := num.Float64()
	outputnum := newNumber().SetFloat64(math.Atan(n))
	outputnum.Quo(outputnum, PIFloatInaccuracy)
	outputnum.Mul(outputnum, PI)
	return outputnum, ArErr{}
}}

var ArCosec = builtinFunc{"cosec", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("cosec expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("cosec expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := newNumber().Set(args[0].(number))
	num.Quo(num, PI)
	num.Mul(num, PIFloatInaccuracy)
	n, _ := num.Float64()
	outputnum := newNumber().SetFloat64(1 / math.Sin(n))
	return outputnum, ArErr{}
}}
var ArArccosec = builtinFunc{"arccosec", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arccosec expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arccosec expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := args[0].(number)
	n, _ := num.Float64()
	if n > -1 && n < 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arccosec expected number between -1 and 1, got %s", anyToArgon(n, true, true, 3, 0, false, 0)),
			EXISTS:  true,
		}
	}
	outputnum := newNumber().SetFloat64(math.Asin(1 / n))
	outputnum.Quo(outputnum, PIFloatInaccuracy)
	outputnum.Mul(outputnum, PI)
	return outputnum, ArErr{}
}}

var ArSec = builtinFunc{"sec", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("sec expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("sec expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := newNumber().Set(args[0].(number))
	num.Quo(num, PI)
	num.Mul(num, PIFloatInaccuracy)
	n, _ := num.Float64()
	outputnum := newNumber().SetFloat64(1 / math.Cos(n))
	return outputnum, ArErr{}
}}

var ArArcsec = builtinFunc{"arcsec", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arcsec expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arcsec expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := args[0].(number)
	n, _ := num.Float64()
	if n > -1 && n < 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arcsec expected number between -1 and 1, got %s", anyToArgon(n, true, true, 3, 0, false, 0)),
			EXISTS:  true,
		}
	}
	outputnum := newNumber().SetFloat64(math.Acos(1 / n))
	outputnum.Quo(outputnum, PIFloatInaccuracy)
	outputnum.Mul(outputnum, PI)
	return outputnum, ArErr{}
}}

var ArCot = builtinFunc{"cot", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("cot expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("cot expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := newNumber().Set(args[0].(number))
	num.Quo(num, PI)
	num.Mul(num, PIFloatInaccuracy)
	n, _ := num.Float64()
	outputnum := newNumber().SetFloat64(1 / math.Tan(n))
	return outputnum, ArErr{}
}}

var ArArccot = builtinFunc{"arccot", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arccot expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("arccot expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := args[0].(number)
	n, _ := num.Float64()
	outputnum := newNumber().SetFloat64(math.Atan(1 / n))
	outputnum.Quo(outputnum, PIFloatInaccuracy)
	outputnum.Mul(outputnum, PI)
	return outputnum, ArErr{}
}}

var ArToDeg = builtinFunc{"toDeg", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("toDeg expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("toDeg expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := newNumber().Set(args[0].(number))
	num.Quo(num, PI)
	num.Mul(num, newNumber().SetInt64(180))
	return num, ArErr{}
}}

var ArToRad = builtinFunc{"toRad", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("toRad expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("toRad expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	num := newNumber().Set(args[0].(number))
	num.Quo(num, newNumber().SetInt64(180))
	num.Mul(num, PI)
	return num, ArErr{}
}}
