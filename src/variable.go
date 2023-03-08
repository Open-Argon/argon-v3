package main

import (
	"strings"
)

var variableCompile = makeRegex(`( *)([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*( *)`)
var setVariableCompile = makeRegex(`( *)(let( +))?([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*(\(( *)((([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*)(( *)\,( *)([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*)*)?( *)\))?( *)=(.|\n)+`)

var blockedVariableNames = map[string]bool{
	"if":       true,
	"else":     true,
	"elif":     true,
	"while":    true,
	"for":      true,
	"break":    true,
	"continue": true,
	"return":   true,
	"let":      true,
	"import":   true,
	"from":     true,
	"do":       true,
}

type accessVariable struct {
	name string
	line int
	code string
	path string
}

type setVariable struct {
	TYPE     string
	name     string
	value    any
	function bool
	params   []string
	line     int
	code     string
	path     string
}

func isVariable(code UNPARSEcode) bool {
	return variableCompile.MatchString(code.code)
}

func parseVariable(code UNPARSEcode) (accessVariable, bool, ArErr, int) {
	name := strings.TrimSpace(code.code)
	if blockedVariableNames[name] {
		return accessVariable{}, false, ArErr{"Naming Error", "Naming Error: \"" + name + "\" is a reserved keyword", code.line, code.path, code.realcode, true}, 1
	}
	return accessVariable{name: name, code: code.code, line: code.line, path: code.path}, true, ArErr{}, 1
}

func readVariable(v accessVariable, stack stack) (any, ArErr) {
	for i := len(stack) - 1; i >= 0; i-- {
		if val, ok := stack[i][v.name]; ok {
			return val, ArErr{}
		}
	}
	return nil, ArErr{"Runtime Error", "variable \"" + v.name + "\" does not exist", v.line, v.path, v.code, true}
}

func isSetVariable(code UNPARSEcode) bool {
	return setVariableCompile.MatchString(code.code)
}

func parseSetVariable(code UNPARSEcode, index int, lines []UNPARSEcode) (setVariable, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	equalsplit := strings.SplitN(trim, "=", 2)
	spacesplit := strings.SplitN(equalsplit[0], " ", 2)
	TYPE := "auto"
	name := strings.TrimSpace(equalsplit[0])
	params := []string{}
	function := false
	if spacesplit[0] == "let" {
		TYPE = "let"
		name = strings.TrimSpace(spacesplit[1])
	}
	if name[len(name)-1] == ')' {
		function = true
		bracketsplit := strings.SplitN(name, "(", 2)
		name = bracketsplit[0]
		params = strings.Split(bracketsplit[1][:len(bracketsplit[1])-1], ",")
		for i := 0; i < len(params); i++ {
			params[i] = strings.TrimSpace(params[i])
		}
	}
	if blockedVariableNames[name] {
		return setVariable{}, false, ArErr{"Naming Error", "Naming Error: \"" + name + "\" is a reserved keyword", code.line, code.path, code.realcode, true}, 1
	}
	value, success, err, i := translateVal(UNPARSEcode{code: equalsplit[1], realcode: code.realcode, line: code.line, path: code.path}, index, lines, true)
	if !success {
		return setVariable{}, false, err, i
	}
	return setVariable{TYPE: TYPE, name: name, value: value, function: function, params: params, line: code.line, code: code.code, path: code.path}, true, ArErr{}, i
}

func setVariableValue(v setVariable, stack stack) (any, ArErr) {
	var resp any
	if v.function {
		resp = Callable{v.name, v.params, v.value, v.code, stack, v.line}
	} else {
		respp, err := runVal(v.value, stack)
		if err.EXISTS {
			return nil, err
		}
		resp = respp
	}
	if v.TYPE == "let" {
		if _, ok := stack[len(stack)-1][v.name]; ok {
			return stack, ArErr{"Runtime Error", "variable \"" + v.name + "\" already exists", v.line, v.path, v.code, true}
		}
		stack[len(stack)-1][v.name] = resp
	} else {
		for i := len(stack) - 1; i >= 0; i-- {
			if _, ok := stack[i][v.name]; ok {
				stack[i][v.name] = resp
				return resp, ArErr{}
			}
		}
		stack[len(stack)-1][v.name] = resp
	}
	return resp, ArErr{}
}
