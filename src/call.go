package main

import (
	"fmt"
	"strings"
)

var callCompile = makeRegex("( *)(.|\n)+\\((.|\n)*\\)( *)")

type call struct {
	Callable any
	Args     []any
	Code     string
	Line     int
	Path     string
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
	for i := 1; i < len(splitby); i++ {
		name := strings.Join(splitby[0:i], "(")
		argstr := strings.Join(splitby[i:], "(")
		args, success, argserr := getValuesFromLetter(argstr, ",", index, codelines, false)
		if !success {
			if i == len(splitby)-1 {
				return nil, false, argserr, 1
			}
			continue
		}
		resp, worked, err, linecount := translateVal(UNPARSEcode{code: name, realcode: code.realcode, line: index + 1, path: code.path}, index, codelines, 0)
		if !worked {
			if i == len(splitby)-1 {
				return nil, false, err, 1
			}
			continue
		}
		return call{Callable: resp, Args: args, Line: code.line, Code: code.realcode, Path: code.path}, true, ArErr{}, linecount
	}
	return nil, false, ArErr{"Syntax Error", "invalid call", code.line, code.path, code.realcode, true}, 1
}

func runCall(c call, stack stack, stacklevel int) (any, ArErr) {
	var callable any = c.Callable
	switch x := c.Callable.(type) {
	case builtinFunc:
		callable = x
	case Callable:
		callable = x
	default:
		callable_, err := runVal(c.Callable, stack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		switch x := callable_.(type) {
		case ArObject:
			callable_, err := mapGet(ArMapGet{
				x,
				[]any{"__call__"},
				true,
				c.Line,
				c.Code,
				c.Path,
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
	for _, arg := range c.Args {
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
				err.line = c.Line
			}
			if err.path == "" {
				err.path = c.Path
			}
			if err.code == "" {
				err.code = c.Code
			}
		}
		return resp, err
	case Callable:
		debugPrintln(x.name, args)
		if len(x.params) != len(args) {
			return nil, ArErr{"Runtime Error", "expected " + fmt.Sprint(len(x.params)) + " arguments, got " + fmt.Sprint(len(args)), c.Line, c.Path, c.Code, true}
		}
		l := anymap{}
		for i, param := range x.params {
			l[param] = args[i]
		}
		resp, err := runVal(x.run, append(x.stack, Map(l)), stacklevel+1)
		return openReturn(resp), err
	}
	return nil, ArErr{"Runtime Error", "type '" + typeof(callable) + "' is not callable", c.Line, c.Path, c.Code, true}
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
		return openReturn(resp), err
	}
	return nil, ArErr{TYPE: "Runtime Error", message: "type '" + typeof(callable) + "' is not callable", EXISTS: true}
}
