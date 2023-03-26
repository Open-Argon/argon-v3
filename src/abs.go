package main

import (
	"fmt"
	"strings"
)

var AbsCompiled = makeRegex(`( *)\|(.|\n)+\|( *)`)

type ABS struct {
	body any
	code string
	line int
	path string
}

func isAbs(code UNPARSEcode) bool {
	return AbsCompiled.MatchString(code.code)
}

func parseAbs(code UNPARSEcode, index int, codelines []UNPARSEcode) (any, bool, ArErr, int) {
	trimmed := strings.TrimSpace(code.code)
	trimmed = trimmed[1 : len(trimmed)-1]

	val, worked, err, i := translateVal(UNPARSEcode{
		trimmed,
		code.realcode,
		code.line,
		code.path,
	}, index, codelines, 0)
	if !worked {
		return nil, false, err, 0
	}
	return ABS{
		val,
		code.realcode,
		code.line,
		code.path,
	}, true, ArErr{}, i
}

func runAbs(x ABS, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(x.body, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	if typeof(resp) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("abs expected number, got %s", typeof(resp)),
			EXISTS:  true,
		}
	}
	return abs(resp.(number)), ArErr{}
}

func abs(x number) number {
	if x.Sign() < 0 {
		return x.Neg(x)
	}
	return x
}

var ArAbs = builtinFunc{"abs", func(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("abs expected 1 argument, got %d", len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("abs expected number, got %s", typeof(args[0])),
			EXISTS:  true,
		}
	}
	return abs(args[0].(number)), ArErr{}
}}
