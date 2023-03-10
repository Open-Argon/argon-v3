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
	val, success, err, i := translateVal(UNPARSEcode{code: trim, realcode: code.realcode, line: 1, path: ""}, 0, []UNPARSEcode{}, false)
	if !success {
		return factorial{}, false, err, i
	}

	return factorial{val, code.code, code.line, code.path}, success, ArErr{}, i
}

func isFactorial(code UNPARSEcode) bool {
	return factorialCompiled.MatchString(code.code)
}

func fact(n number) number {
	if n.Cmp(newNumber().SetInt64(0)) == 0 {
		return newNumber().SetInt64(1)
	}
	result := newNumber().SetInt64(1)
	for i := newNumber().SetInt64(2); i.Cmp(n) <= 0; i.Add(i, newNumber().SetInt64(1)) {
		result.Mul(result, i)
	}
	return result
}

func runFactorial(f factorial, stack stack) (any, ArErr) {
	val, err := runVal(f.value, stack)
	if err.EXISTS {
		return nil, err
	}
	switch x := val.(type) {
	case number:
		if !x.IsInt() {
			return nil, ArErr{"Runtime Error", "cannot use factorial on non-integer", f.line, f.path, f.code, true}
		}
		if x.Cmp(newNumber().SetInt64(0)) == -1 {
			return nil, ArErr{"Runtime Error", "cannot use factorial on negative number", f.line, f.path, f.code, true}
		}
		return fact(x), ArErr{}
	default:
		return nil, ArErr{"Runtime Error", "cannot use factorial on non-number of type '" + typeof(val) + "'", f.line, f.path, f.code, true}
	}
}
