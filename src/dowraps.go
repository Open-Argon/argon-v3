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
	var allCodelines = []UNPARSEcode{}
	i := index + 1
	for ; i < len(codelines); i++ {

		if isBlank(codelines[i]) {
			continue
		}
		indent := len(codelines[i].code) - len(strings.TrimLeft(codelines[i].code, " "))
		if setindent == -1 {
			setindent = indent
		}
		if indent <= currentindent {
			break
		}
		allCodelines = append(allCodelines, codelines[i])
	}
	finalLines := i
	codelines = allCodelines
	translated := []any{}
	for i := 0; i < len(codelines); {
		indent := len(codelines[i].code) - len(strings.TrimLeft(codelines[i].code, " "))
		if indent != setindent {
			return nil, false, ArErr{"Syntax Error", "invalid indent", code.line, code.path, codelines[i].code, true}, 1
		}

		val, _, err, step := translateVal(codelines[i], i, codelines, 3)
		i += step
		if err.EXISTS {
			return nil, false, err, i - index
		}
		translated = append(translated, val)
	}
	return dowrap{run: translated, line: code.line, path: code.path, code: code.realcode}, true, ArErr{}, finalLines - index
}

func runDoWrap(d dowrap, Stack stack, stacklevel int) (any, ArErr) {
	newstack := append(Stack, newscope())
	newstackCopy := make(stack, len(newstack))
	copy(newstackCopy, newstack)
	for _, v := range d.run {
		val, err := runVal(v, newstackCopy, stacklevel+1)
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
