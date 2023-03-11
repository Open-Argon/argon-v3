package main

import "strings"

var returnCompile = makeRegex(`( *)return( +)(.|\n)+`)

type CallJumpStatment struct {
	TYPE  string
	value any
	line  int
	code  string
	path  string
}

type PassBackJumpStatment struct {
	TYPE  string
	value any
	line  int
	code  string
	path  string
}

func isReturn(code UNPARSEcode) bool {
	return returnCompile.MatchString(code.code)
}

func parseReturn(code UNPARSEcode, index int, codeline []UNPARSEcode) (CallJumpStatment, bool, ArErr, int) {
	resp, worked, err, i := translateVal(UNPARSEcode{
		code:     strings.TrimSpace(code.code)[6:],
		realcode: code.realcode,
		line:     code.line,
		path:     code.path,
	}, index, codeline, 1)
	return CallJumpStatment{
		TYPE:  "return",
		value: resp,
		line:  code.line,
		code:  code.realcode,
		path:  code.path,
	}, worked, err, i
}

func runJumpStatment(code CallJumpStatment, stack stack, stacklevel int) (any, ArErr) {
	var val any
	if code.value != nil {
		v, err := runVal(code.value, stack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		val = v
	}
	return PassBackJumpStatment{
		TYPE:  code.TYPE,
		value: val,
		line:  code.line,
		code:  code.code,
		path:  code.path,
	}, ArErr{}
}

func openJump(resp any) any {
	switch x := resp.(type) {
	case PassBackJumpStatment:
		return x.value
	default:
		return resp
	}
}
