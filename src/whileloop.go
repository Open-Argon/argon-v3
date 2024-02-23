package main

import (
	"strings"
)

var whileLoopCompiled = makeRegex(`( *)while( )+\((.|\n)+\)( )+(.|\n)+`)
var foreverLoopCompiled = makeRegex(`( *)forever( )+(.|\n)+`)

type whileLoop struct {
	condition any
	body      any
	line      int
	code      string
	path      string
}

func isWhileLoop(code UNPARSEcode) bool {
	return whileLoopCompiled.MatchString(code.code)
}

func isForeverLoop(code UNPARSEcode) bool {
	return foreverLoopCompiled.MatchString(code.code)
}

func parseWhileLoop(code UNPARSEcode, index int, codeline []UNPARSEcode) (whileLoop, bool, ArErr, int) {
	trimmed := strings.TrimSpace(code.code)
	trimmed = strings.TrimSpace(trimmed[strings.Index(trimmed, "("):])
	trimmed = (trimmed[1:])
	split := strings.Split(trimmed, ")")
	for j := len(split) - 1; j > 0; j-- {
		conditionjoined := strings.Join(split[:j], ")")
		thenjoined := strings.Join(split[j:], ")")
		outindex := 0
		conditionval, worked, err, step := translateVal(
			UNPARSEcode{
				code:     conditionjoined,
				realcode: code.realcode,
				line:     code.line,
				path:     code.path,
			},
			index,
			codeline,
			0,
		)
		if err.EXISTS || !worked {
			if j == 1 {
				return whileLoop{}, worked, err, step
			} else {
				continue
			}
		}
		outindex += step
		thenval, worked, err, step := translateVal(
			UNPARSEcode{
				code:     thenjoined,
				realcode: code.realcode,
				line:     code.line,
				path:     code.path,
			},
			index+outindex-1,
			codeline,
			3,
		)
		if err.EXISTS || !worked {
			return whileLoop{}, worked, err, step
		}
		outindex += step - 1
		return whileLoop{
			condition: conditionval,
			body:      thenval,
			line:      code.line,
			code:      code.realcode,
			path:      code.path,
		}, true, ArErr{}, outindex
	}
	return whileLoop{}, false, ArErr{
		"Syntax Error",
		"Could not parse while loop",
		code.line,
		code.path,
		code.realcode,
		true,
	}, 0
}

func parseForeverLoop(code UNPARSEcode, index int, codeline []UNPARSEcode) (whileLoop, bool, ArErr, int) {
	trimmed := strings.TrimSpace(code.code)
	trimmed = strings.TrimSpace(trimmed[7:])
	thenval, worked, err, step := translateVal(
		UNPARSEcode{
			code:     trimmed,
			realcode: code.realcode,
			line:     code.line,
			path:     code.path,
		},
		index,
		codeline,
		3,
	)
	return whileLoop{
		condition: true,
		body:      thenval,
		line:      code.line,
		code:      code.realcode,
		path:      code.path,
	}, worked, err, step
}

func runWhileLoop(loop whileLoop, stack stack, stacklevel int) (any, ArErr) {

	newstack := append(stack, newscope())
	for {
		condition, err := runVal(loop.condition, newstack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		newbodystack := append(newstack, newscope())
		if !anyToBool(condition) {
			break
		}
		resp, err := runVal(loop.body, newbodystack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		switch x := resp.(type) {
		case Return:
			return x, ArErr{}
		case Break:
			return nil, ArErr{}
		case Continue:
			continue
		}
	}
	return nil, ArErr{}
}
