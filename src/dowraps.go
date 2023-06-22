package main

import (
	"strings"
)

var dowrapCompile = makeRegex("( )*do( )*")

type dowrap struct {
	run  []any
	line int
	path string
	code string
}

func isDoWrap(code UNPARSEcode) bool {
	return dowrapCompile.MatchString(code.code)
}

func parseDoWrap(code UNPARSEcode, index int, codelines []UNPARSEcode) (any, bool, ArErr, int) {
	currentindent := len(code.realcode) - len(strings.TrimLeft(code.realcode, " "))
	var setindent int = -1
	var i = index + 1
	translated := []any{}
	for i < len(codelines) {
		if isBlank(codelines[i]) {
			i++
			continue
		}
		indent := len(codelines[i].code) - len(strings.TrimLeft(codelines[i].code, " "))
		if setindent == -1 {
			setindent = indent
		}
		if indent <= currentindent {
			break
		}
		if indent != setindent {
			return nil, false, ArErr{"Syntax Error", "invalid indent", i, code.path, codelines[i].code, true}, 1
		}

		val, _, err, step := translateVal(codelines[i], i, codelines, 2)
		i += step
		if err.EXISTS {
			return nil, false, err, i - index
		}
		translated = append(translated, val)
	}
	return dowrap{run: translated, line: code.line, path: code.path, code: code.realcode}, true, ArErr{}, i - index
}

func runDoWrap(d dowrap, stack stack, stacklevel int) (any, ArErr) {
	newstack := append(stack, newscope())
	for _, v := range d.run {
		val, err := runVal(v, newstack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		switch x := val.(type) {
		case Return, Break, Continue:
			return x, ArErr{}
		}
	}
	return nil, ArErr{}
}
