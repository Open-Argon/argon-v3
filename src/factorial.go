package main

import (
	"strings"
)

var factorialCompiled = makeRegex(`( *)(.|\n)+\!( *)`)

type factorial struct {
	value any
	code  string
	line  int
	path  string
}

func parseFactorial(code UNPARSEcode, index int, codeline []UNPARSEcode) (factorial, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	trim = trim[:len(trim)-1]
	val, success, err, i := translateVal(UNPARSEcode{code: trim, realcode: code.realcode, line: code.line, path: code.path}, index, codeline, 0)
	if !success {
		return factorial{}, false, err, i
	}

	return factorial{val, code.code, code.line, code.path}, success, ArErr{}, i
}

func isFactorial(code UNPARSEcode) bool {
	return factorialCompiled.MatchString(code.code)
}

func runFactorial(f factorial, stack stack, stacklevel int) (any, ArErr) {
	val, err := runVal(f.value, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	switch val := val.(type) {
	case ArObject:
		if callable, ok := val.obj["__factorial__"]; ok {
			return runCall(call{
				Callable: callable,
				Args:     []any{},
				Code:     f.code,
				Line:     f.line,
				Path:     f.path,
			}, stack, stacklevel)
		}
	}
	return nil, ArErr{
		TYPE:    "TypeError",
		message: "factorial not defined for type",
		EXISTS:  true,
	}
}
