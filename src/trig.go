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
