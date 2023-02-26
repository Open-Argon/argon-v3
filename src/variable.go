package main

import (
	"strings"
)

var variableCompile = makeRegex(`([a-zA-Z_])([a-zA-Z0-9_])*`)

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

type variableValue struct {
	VAL    any
	EXISTS any
	origin string
}

type accessVariable struct {
	name string
	line int
	code string
	path string
}

func isVariable(code UNPARSEcode) bool {
	return variableCompile.MatchString(code.code)
}

func parseVariable(code UNPARSEcode) (accessVariable, bool, ArErr, int) {
	name := strings.TrimSpace(code.code)
	if blockedVariableNames[name] {
		return accessVariable{}, false, ArErr{"Naming Error", "Naming Error: \"" + name + "\" is a reserved keyword", code.line, code.path, code.realcode, true}, 1
	}
	return accessVariable{name: name, code: code.code, line: code.line}, true, ArErr{}, 1
}

func readVariable(v accessVariable, stack stack) (any, ArErr) {
	for i := len(stack) - 1; i >= 0; i-- {
		if val, ok := stack[i][v.name]; ok {
			return val.VAL, ArErr{}
		}
	}
	return nil, ArErr{"Runtime Error", "variable \"" + v.name + "\" does not exist", v.line, v.path, v.code, true}
}
