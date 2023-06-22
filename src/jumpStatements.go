package main

import (
	"strings"
)

var returnCompile = makeRegex(`( *)return(( +)(.|\n)+)?`)
var breakCompile = makeRegex(`( *)break( *)`)
var continueCompile = makeRegex(`( *)continue( *)`)

type CallReturn struct {
	value any
	line  int
	code  string
	path  string
}

type Return struct {
	value any
	line  int
	code  string
	path  string
}
type Break struct {
	line int
	code string
	path string
}
type Continue struct {
	line int
	code string
	path string
}

func isReturn(code UNPARSEcode) bool {
	return returnCompile.MatchString(code.code)
}

func isBreak(code UNPARSEcode) bool {
	return breakCompile.MatchString(code.code)
}

func isContinue(code UNPARSEcode) bool {
	return continueCompile.MatchString(code.code)
}

func parseReturn(code UNPARSEcode, index int, codeline []UNPARSEcode) (CallReturn, bool, ArErr, int) {
	val := strings.TrimSpace(code.code)[6:]
	var resp any
	var worked, err, i = true, ArErr{}, 1
	if val != "" {
		resp, worked, err, i = translateVal(UNPARSEcode{
			code:     val,
			realcode: code.realcode,
			line:     code.line,
			path:     code.path,
		}, index, codeline, 1)
	}
	return CallReturn{
		value: resp,
		line:  code.line,
		code:  code.realcode,
		path:  code.path,
	}, worked, err, i
}

func runReturn(code CallReturn, stack stack, stacklevel int) (any, ArErr) {
	var val any
	if code.value != nil {
		v, err := runVal(code.value, stack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		val = v
	}
	return Return{
		value: val,
		line:  code.line,
		code:  code.code,
		path:  code.path,
	}, ArErr{}
}

func openReturn(resp any) any {
	switch x := resp.(type) {
	case Return:
		return x.value
	default:
		return resp
	}
}

func parseBreak(code UNPARSEcode) (Break, bool, ArErr, int) {
	return Break{
		line: code.line,
		code: code.realcode,
		path: code.path,
	}, true, ArErr{}, 1
}

func parseContinue(code UNPARSEcode) (Continue, bool, ArErr, int) {
	return Continue{
		line: code.line,
		code: code.realcode,
		path: code.path,
	}, true, ArErr{}, 1
}

func ThrowOnNonLoop(val any, err ArErr) (any, ArErr) {
	switch x := val.(type) {
	case Break:
		return nil, ArErr{"Break Error", "break can only be used in loops", x.line, x.path, x.code, true}
	case Continue:
		return nil, ArErr{"Continue Error", "continue can only be used in loops", x.line, x.path, x.code, true}
	default:
		return x, err
	}
}
