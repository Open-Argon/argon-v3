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

type Callable struct {
	name   string
	params []string
	run    any
	code   string
	stack  stack
	line   int
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
		args, success, argserr := getValuesFromLetter(argstr, ",", index, codelines, false)
		arguments = args
		if !success {
			if i == len(splitby)-1 {
				return nil, false, argserr, 1
			}
			continue
		}
		resp, worked, _, _ := translateVal(UNPARSEcode{code: name, realcode: code.realcode, line: index + 1, path: code.path}, index, codelines, 0)
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

func runCall(c call, stack stack, stacklevel int) (any, ArErr) {
	var callable any = c.callable
	switch x := c.callable.(type) {
	case builtinFunc:
		callable = x
	case Callable:
		callable = x
	default:
		callable_, err := runVal(c.callable, stack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		switch x := callable_.(type) {
		case ArObject:
			callable_, err := mapGet(ArMapGet{
				x,
				[]any{"__call__"},
				true,
				c.line,
				c.code,
				c.path,
			}, stack, stacklevel+1)
			if !err.EXISTS {
				callable = callable_
			}
		default:
			callable = callable_
		}
	}
	args := []any{}
	level := append(stack, newscope())
	for _, arg := range c.args {
		resp, err := runVal(arg, level, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		args = append(args, resp)
	}
	switch x := callable.(type) {
	case builtinFunc:
		debugPrintln(x.name, args)
		resp, err := x.FUNC(args...)
		resp = AnyToArValid(resp)
		if err.EXISTS {
			if err.line == 0 {
				err.line = c.line
			}
			if err.path == "" {
				err.path = c.path
			}
			if err.code == "" {
				err.code = c.code
			}
		}
		return resp, err
	case Callable:
		debugPrintln(x.name, args)
		if len(x.params) != len(args) {
			return nil, ArErr{"Runtime Error", "expected " + fmt.Sprint(len(x.params)) + " arguments, got " + fmt.Sprint(len(args)), c.line, c.path, c.code, true}
		}
		l := anymap{}
		for i, param := range x.params {
			l[param] = args[i]
		}
		resp, err := runVal(x.run, append(x.stack, Map(l)), stacklevel+1)
		return ThrowOnNonLoop(openReturn(resp), err)
	}
	return nil, ArErr{"Runtime Error", "type '" + typeof(callable) + "' is not callable", c.line, c.path, c.code, true}
}

func builtinCall(callable any, args []any) (any, ArErr) {
	debugPrintln(callable, args)

	switch x := callable.(type) {
	case builtinFunc:
		resp, err := x.FUNC(args...)
		resp = AnyToArValid(resp)
		return resp, err
	case Callable:
		if len(x.params) != len(args) {
			return nil, ArErr{TYPE: "Runtime Error", message: "expected " + fmt.Sprint(len(x.params)) + " arguments, got " + fmt.Sprint(len(args)), EXISTS: true}
		}
		level := newscope()
		for i, param := range x.params {
			level.obj[param] = args[i]
		}
		resp, err := runVal(x.run, append(x.stack, level), 0)
		return ThrowOnNonLoop(openReturn(resp), err)
	}
	return nil, ArErr{TYPE: "Runtime Error", message: "type '" + typeof(callable) + "' is not callable", EXISTS: true}
}
