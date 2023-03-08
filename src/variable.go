package main

import (
	"fmt"
	"strings"
)

var variableCompile = makeRegex(`( *)([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*( *)`)
var validname = makeRegex(`(.|\n)+(\(( *)((([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*)(( *)\,( *)([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*)*)?( *)\))`)
var setVariableCompile = makeRegex(`( *)(let( +))(.|\n)+( *)=(.|\n)+`)
var autoAsignVariableCompile = makeRegex(`(.|\n)+=(.|\n)+`)
var deleteVariableCompile = makeRegex(`( *)delete( +)( *)`)

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
	"true":     true,
	"false":    true,
	"null":     true,
	"delete":   true,
}

type accessVariable struct {
	name string
	line int
	code string
	path string
}

type setVariable struct {
	TYPE     string
	toset    any
	value    any
	function bool
	params   []string
	line     int
	code     string
	path     string
}

type setFunction struct {
	toset  any
	params []string
}

func isVariable(code UNPARSEcode) bool {
	return variableCompile.MatchString(code.code)
}

func parseVariable(code UNPARSEcode) (accessVariable, bool, ArErr, int) {
	name := strings.TrimSpace(code.code)
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

func isAutoAsignVariable(code UNPARSEcode) bool {
	return autoAsignVariableCompile.MatchString(code.code)
}

func isDeleteVariable(code UNPARSEcode) bool {
	return deleteVariableCompile.MatchString(code.code)
}

func nameToTranslated(code UNPARSEcode, index int, lines []UNPARSEcode) (any, bool, ArErr, int) {
	if validname.MatchString(code.code) {
		trimmed := strings.TrimSpace(code.code)
		start := strings.LastIndex(trimmed, "(")
		params := strings.Split(trimmed[start+1:len(trimmed)-1], ",")
		for i := range params {
			params[i] = strings.TrimSpace(params[i])
		}
		name := strings.TrimSpace(trimmed[:start])
		fmt.Println(name)
		if blockedVariableNames[name] {
			fmt.Println(name)
			return accessVariable{}, false, ArErr{"Naming Error", "\"" + name + "\" is a reserved keyword", code.line, code.path, code.realcode, true}, 1
		}
		value, success, err, i := translateVal(UNPARSEcode{
			code:     name,
			realcode: code.realcode,
			line:     code.line,
			path:     code.path,
		}, index, lines, false)
		return setFunction{toset: value, params: params}, success, err, i
	}
	return translateVal(code, index, lines, false)
}

func parseSetVariable(code UNPARSEcode, index int, lines []UNPARSEcode) (setVariable, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	equalsplit := strings.SplitN(trim, "=", 2)
	spacesplit := strings.SplitN(equalsplit[0], " ", 2)
	name := strings.TrimSpace(spacesplit[1])
	params := []string{}
	function := false
	if blockedVariableNames[name] {
		return setVariable{}, false, ArErr{"Naming Error", "\"" + name + "\" is a reserved keyword", code.line, code.path, code.realcode, true}, 1
	}
	toset, success, err, i := nameToTranslated(UNPARSEcode{code: name, realcode: code.realcode, line: code.line, path: code.path}, index, lines)
	if err.EXISTS {
		return setVariable{}, success, err, i
	}
	switch toset.(type) {
	case accessVariable:
		break
	case setFunction:
		function = true
		params = toset.(setFunction).params
		toset = toset.(setFunction).toset
	default:
		return setVariable{}, false, ArErr{"Type Error", "can't set for non variable, did you mean '=='?", code.line, code.path, code.realcode, true}, 1
	}
	value, success, err, i := translateVal(UNPARSEcode{code: equalsplit[1], realcode: code.realcode, line: code.line, path: code.path}, index, lines, false)
	if !success {
		return setVariable{}, false, err, i
	}
	return setVariable{TYPE: "let", toset: toset, value: value, function: function, params: params, line: code.line, code: code.code, path: code.path}, true, ArErr{}, i
}

func parseAutoAsignVariable(code UNPARSEcode, index int, lines []UNPARSEcode) (setVariable, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	equalsplit := strings.SplitN(trim, "=", 2)
	name := strings.TrimSpace(equalsplit[0])
	params := []string{}
	function := false
	if blockedVariableNames[name] {
		return setVariable{}, false, ArErr{"Naming Error", "\"" + name + "\" is a reserved keyword", code.line, code.path, code.realcode, true}, 1
	}
	toset, success, err, i := nameToTranslated(UNPARSEcode{code: name, realcode: code.realcode, line: code.line, path: code.path}, index, lines)
	if err.EXISTS {
		return setVariable{}, success, err, i
	}
	switch toset.(type) {
	case accessVariable:
		break
	case ArMapGet:
		break
	case setFunction:
		function = true
		params = toset.(setFunction).params
		toset = toset.(setFunction).toset
	default:
		return setVariable{}, false, ArErr{"Type Error", "can't set for non variable, did you mean '=='?", code.line, code.path, code.realcode, true}, 1
	}
	value, success, err, i := translateVal(UNPARSEcode{code: equalsplit[1], realcode: code.realcode, line: code.line, path: code.path}, index, lines, false)
	if !success {
		return setVariable{}, false, err, i
	}
	return setVariable{TYPE: "auto", toset: toset, value: value, function: function, params: params, line: code.line, code: code.code, path: code.path}, true, ArErr{}, i
}

func setVariableValue(v setVariable, stack stack) (any, ArErr) {
	var resp any
	if v.function {
		resp = Callable{v.params, v.value, v.code, stack, v.line}
	} else {
		respp, err := runVal(v.value, stack)
		if err.EXISTS {
			return nil, err
		}
		resp = respp
	}

	if v.TYPE == "let" {
		if _, ok := stack[len(stack)-1][v.toset.(accessVariable).name]; ok {
			return stack, ArErr{"Runtime Error", "variable \"" + v.toset.(accessVariable).name + "\" already exists", v.line, v.path, v.code, true}
		}
		stack[len(stack)-1][v.toset.(accessVariable).name] = resp
	} else {
		switch x := v.toset.(type) {
		case accessVariable:
			for i := len(stack) - 1; i >= 0; i-- {
				if _, ok := stack[i][x.name]; ok {
					stack[i][x.name] = resp
					return resp, ArErr{}
				}
			}
			stack[len(stack)-1][x.name] = resp
		case ArMapGet:
			respp, err := runVal(x.VAL, stack)
			if err.EXISTS {
				return nil, err
			}
			key, err := runVal(x.key, stack)
			if err.EXISTS {
				return nil, err
			}
			switch y := respp.(type) {
			case ArMap:
				y[key] = resp
			default:
				return nil, ArErr{"Runtime Error", "can't set for non map", v.line, v.path, v.code, true}
			}
		}
	}
	return resp, ArErr{}
}

/*

func parseAutosAsignVariable(code UNPARSEcode, index int, lines []UNPARSEcode) (setVariable, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	equalsplit := strings.SplitN(trim, "=", 2)
	name := strings.TrimSpace(equalsplit[0])
	params := []string{}
	function := false

	if blockedVariableNames[name] {
		return setVariable{}, false, ArErr{"Naming Error", "\"" + name + "\" is a reserved keyword", code.line, code.path, code.realcode, true}, 1
	}
	value, success, err, i := translateVal(UNPARSEcode{code: equalsplit[1], realcode: code.realcode, line: code.line, path: code.path}, index, lines, false)
	if !success {
		return setVariable{}, false, err, i
	}
	return setVariable{TYPE: "let", name: name, value: value, function: function, params: params, line: code.line, code: code.code, path: code.path}, true, ArErr{}, i
}


*/
