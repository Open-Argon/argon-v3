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
	value, err := runVal(x, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	switch value := value.(type) {
	case ArObject:
		if Callable, ok := value.obj["__abs__"]; ok {
			return runCall(call{
				Callable: Callable,
				Args:     []any{},
				Code:     x.code,
				Line:     x.line,
				Path:     x.path,
			}, stack, stacklevel)
		}
	}
	return nil, ArErr{
		"TypeError",
		fmt.Sprint("abs() not supported on ", typeof(value)),
		x.line,
		x.path,
		x.code,
		true,
	}
}
