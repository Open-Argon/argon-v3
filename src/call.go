package main

import (
	"fmt"
	"strings"
)

var callCompile = makeRegex("( *)(.|\n)+\\((.|\n)*\\)( *)")

type call struct {
	callable any
	args     []any
	code     string
	line     int
	path     string
}

func isCall(code UNPARSEcode) bool {
	return callCompile.MatchString(code.code)
}

func parseCall(code UNPARSEcode, index int, codelines []UNPARSEcode) (any, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	trim = trim[:len(trim)-1]
	splitby := strings.Split(trim, "(")

	var works bool
	var callable any
	var arguments []any
	for i := 1; i < len(splitby); i++ {
		name := strings.Join(splitby[0:i], "(")
		argstr := strings.Join(splitby[i:], "(")
		args, success, argserr := getValuesFromCommas(argstr, index, codelines)
		arguments = args
		if !success {
			if i == len(splitby)-1 {
				return nil, false, argserr, 1
			}
			continue
		}
		resp, worked, _, _ := translateVal(UNPARSEcode{code: name, realcode: code.realcode, line: index + 1, path: code.path}, index, codelines, false)
		if !worked {
			if i == len(splitby)-1 {
				return nil, false, ArErr{"Syntax Error", "invalid callable", code.line, code.path, code.realcode, true}, 1
			}
			continue
		}
		works = true
		callable = resp
		break
	}
	if !works {
		return nil, false, ArErr{"Syntax Error", "invalid call", code.line, code.path, code.realcode, true}, 1
	}
	return call{callable: callable, args: arguments, line: code.line, code: code.realcode, path: code.path}, true, ArErr{}, 1
}

func runCall(c call, stack stack) (any, ArErr) {
	callable, err := runVal(c.callable, stack)
	if err.EXISTS {
		return nil, err
	}
	args := []any{}
	level := append(stack, map[string]any{})
	for _, arg := range c.args {
		resp, err := runVal(arg, level)
		if err.EXISTS {
			return nil, err
		}
		args = append(args, resp)
	}
	switch x := callable.(type) {
	case builtinFunc:
		resp, err := x.FUNC(args...)
		if err.EXISTS {
			err = ArErr{err.TYPE, err.message, c.line, c.path, c.code, true}
		}
		return resp, err
	case Callable:
		if len(x.params) != len(args) {
			return nil, ArErr{"Runtime Error", "expected " + fmt.Sprint(len(x.params)) + " arguments, got " + fmt.Sprint(len(args)), c.line, c.path, c.code, true}
		}
		level := map[string]any{}
		for i, param := range x.params {
			level[param] = args[i]
		}
		resp, err := runVal(x.run, append(stack, level))
		return resp, err
	}
	return nil, ArErr{"Runtime Error", "type '" + typeof(callable) + "' is not callable", c.line, c.path, c.code, true}
}
