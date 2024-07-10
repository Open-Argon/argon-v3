package main

import (
	"strings"
)

var spacelessVariable = `([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*`
var SpacelessVariableCompiled = makeRegex(spacelessVariable)
var variableCompile = makeRegex(`( *)` + spacelessVariable + `( *)`)
var validname = makeRegex(`(.|\n)*(\(( *)((([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*)(( *)\,( *)([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*)*)?( *)\))`)
var setVariableCompile = makeRegex(`( *)(let( +))(.|\n)+( *)=(.|\n)+`)
var autoAsignVariableCompile = makeRegex(`(.|\n)+=(.|\n)+`)
var deleteVariableCompile = makeRegex(`( *)delete( +)(.|\n)+( *)`)

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
	"not":      true,
	"and":      true,
	"or":       true,
	"try":      true,
	"catch":    true,
}

type accessVariable struct {
	Name string
	Line int
	Code string
	Path string
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

type ArDelete struct {
	value any
	line  int
	code  string
	path  string
}

func isVariable(code UNPARSEcode) bool {
	return variableCompile.MatchString(code.code)
}

func parseVariable(code UNPARSEcode) (accessVariable, bool, ArErr, int) {
	name := strings.TrimSpace(code.code)
	return accessVariable{Name: name, Code: code.realcode, Line: code.line, Path: code.path}, true, ArErr{}, 1
}

func readVariable(v accessVariable, stack stack) (any, ArErr) {
	name := ArString(v.Name)
	for i := len(stack) - 1; i >= 0; i-- {
		callable, ok := stack[i].obj["__Contains__"]
		if !ok {
			continue
		}
		contains, err := builtinCall(callable, []any{name})
		if err.EXISTS {
			return nil, err
		}
		if anyToBool(contains) {
			callable, ok := stack[i].obj["__getindex__"]
			if !ok {
				continue
			}
			return builtinCall(callable, []any{name})
		}
	}
	return nil, ArErr{"Name Error", "variable \"" + v.Name + "\" does not exist", v.Line, v.Path, v.Code, true}
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
			if params[i] == "" {
				params = append(params[:i], params[i+1:]...)
			}
		}
		name := strings.TrimSpace(trimmed[:start])
		if name == "" {
			return setFunction{toset: nil, params: params}, true, ArErr{}, 1
		}
		if blockedVariableNames[name] {
			return accessVariable{}, false, ArErr{"Naming Error", "\"" + name + "\" is a reserved keyword", code.line, code.path, code.realcode, true}, 1
		}
		value, success, err, i := translateVal(UNPARSEcode{
			code:     name,
			realcode: code.realcode,
			line:     code.line,
			path:     code.path,
		}, index, lines, 0)
		return setFunction{toset: value, params: params}, success, err, i
	}
	return translateVal(code, index, lines, 0)
}

func parseSetVariable(code UNPARSEcode, index int, lines []UNPARSEcode, isLine int) (setVariable, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	equalsplit := strings.SplitN(trim, "=", 2)
	spacesplit := strings.SplitN(equalsplit[0], " ", 2)
	name := strings.TrimSpace(spacesplit[1])
	params := []string{}
	function := false
	if blockedVariableNames[name] {
		return setVariable{}, false, ArErr{"Naming Error", "\"" + name + "\" is a reserved keyword", code.line, code.path, code.realcode, true}, 1
	}
	toset, success, err, namei := nameToTranslated(UNPARSEcode{code: name, realcode: code.realcode, line: code.line, path: code.path}, index, lines)
	if err.EXISTS {
		return setVariable{}, success, err, namei
	}
	switch x := toset.(type) {
	case accessVariable:
		break
	case setFunction:
		function = true
		params = x.params
		toset = x.toset
		if toset == nil {
			return setVariable{}, false, ArErr{"Type Error", "can't set for non variable, did you mean '=='?", code.line, code.path, code.realcode, true}, 1
		}
	default:
		return setVariable{}, false, ArErr{"Type Error", "can't set for non variable, did you mean '=='?", code.line, code.path, code.realcode, true}, 1
	}
	value, success, err, i := translateVal(UNPARSEcode{code: equalsplit[1], realcode: code.realcode, line: code.line, path: code.path}, index, lines, isLine)
	if !success {
		return setVariable{}, false, err, i
	}
	return setVariable{TYPE: "let", toset: toset, value: value, function: function, params: params, line: code.line, code: code.code, path: code.path}, true, ArErr{}, i + namei - 1
}

var operationsToInt = map[byte]int{
	'+': 10,
	'-': 11,
	'*': 12,
	'/': 15,
	'^': 16,
	'&': 0,
	'|': 1,
}

func parseAutoAsignVariable(code UNPARSEcode, index int, lines []UNPARSEcode, isLine int) (setVariable, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	equalsplit := strings.Split(trim, "=")
	for i := 1; i < len(equalsplit); i++ {
		name := strings.TrimSpace(strings.Join(equalsplit[:i], "="))
		if name == "" {
			continue
		}
		operation := name[len(name)-1]
		var operationtype int = -1
		if operation == '+' || operation == '-' || operation == '*' || operation == '/' || operation == '^' || operation == '&' || operation == '|' {
			name = strings.TrimSpace(name[:len(name)-1])
			if n, ok := operationsToInt[operation]; ok {
				operationtype = n
			}
		}
		params := []string{}
		function := false
		if blockedVariableNames[name] {
			if i == len(equalsplit)-1 {
				return setVariable{}, false, ArErr{"Naming Error", "\"" + name + "\" is a reserved keyword", code.line, code.path, code.realcode, true}, 1
			}
			continue
		}
		toset, success, err, namei := nameToTranslated(UNPARSEcode{code: name, realcode: code.realcode, line: code.line, path: code.path}, index, lines)
		if err.EXISTS {
			if i == len(equalsplit)-1 {
				return setVariable{}, success, err, namei
			}
			continue
		}
		switch x := toset.(type) {
		case accessVariable:
			break
		case ArMapGet:
			break
		case setFunction:
			function = true
			params = x.params
			toset = x.toset
		default:
			if i == len(equalsplit)-1 {
				return setVariable{}, false, ArErr{"Type Error", "can't set for non variable, did you mean '=='?", code.line, code.path, code.realcode, true}, 1
			}
			continue
		}
		value, success, err, i := translateVal(UNPARSEcode{code: strings.Join(equalsplit[i:], "="), realcode: code.realcode, line: code.line, path: code.path}, index, lines, isLine)
		if !success {
			return setVariable{}, false, err, i
		}
		if operationtype != -1 {
			value = operationType{
				operation: operationtype,
				values:    []any{toset, value},
				line:      code.line,
				code:      code.code,
				path:      code.path,
			}
		}
		return setVariable{TYPE: "auto", toset: toset, value: value, function: function, params: params, line: code.line, code: code.code, path: code.path}, true, ArErr{}, i + namei - 1
	}
	return setVariable{}, false, ArErr{"Syntax Error", "invalid syntax", code.line, code.path, code.realcode, true}, 1
}

func setVariableValue(v setVariable, stack stack, stacklevel int) (any, ArErr) {
	var resp any
	if v.function {
		resp = Callable{"anonymous", v.params, v.value, v.code, stack, v.line}
	} else {
		respp, err := runVal(v.value, stack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		resp = openReturn(respp)
	}

	if v.TYPE == "let" {
		name := v.toset.(accessVariable).Name
		if v.function {
			resp = Callable{name, v.params, v.value, v.code, stack, v.line}
		}
		stackcallable, ok := stack[len(stack)-1].obj["__setindex__"]
		if !ok {
			return nil, ArErr{"Type Error", "stack doesn't have __setindex__", v.line, v.path, v.code, true}
		}
		_, err := builtinCall(stackcallable, []any{name, resp})
		if err.EXISTS {
			return nil, err
		}
	} else {
		switch x := v.toset.(type) {
		case accessVariable:
			name := x.Name
			hasSet := false
			if v.function {
				resp = Callable{name, v.params, v.value, v.code, stack, v.line}
			}
			for i := len(stack) - 1; i >= 0; i-- {
				callable, ok := stack[i].obj["__Contains__"]
				if !ok {
					continue
				}
				contains, err := builtinCall(callable, []any{name})
				if err.EXISTS {
					return nil, err
				}
				if anyToBool(contains) {
					callable, ok := stack[i].obj["__setindex__"]
					if !ok {
						continue
					}
					builtinCall(callable, []any{name, resp})
					hasSet = true
					break
				}
			}
			if !hasSet {
				callable, ok := stack[len(stack)-1].obj["__setindex__"]
				if !ok {
					return nil, ArErr{"Type Error", "stack doesn't have __setindex__", v.line, v.path, v.code, true}
				}
				builtinCall(callable, []any{name, resp})
			}
		case ArMapGet:
			respp, err := runVal(x.VAL, stack, stacklevel+1)
			if err.EXISTS {
				return nil, err
			}
			if len(x.Args) != 1 {
				return nil, ArErr{"Runtime Error", "cannot set by slice", v.line, v.path, v.code, true}
			}
			key, err := runVal(x.Args[0], stack, stacklevel+1)
			key = ArValidToAny(key)
			if err.EXISTS {
				return nil, err
			}
			switch y := respp.(type) {
			case ArObject:
				if _, ok := y.obj["__setindex__"]; ok {
					callable := y.obj["__setindex__"]
					_, err := runCall(call{
						Callable: callable,
						Args:     []any{key, resp},
						Line:     v.line,
						Path:     v.path,
						Code:     v.code,
					}, stack, stacklevel+1)
					if err.EXISTS {
						return nil, err
					}
				}
			default:
				return nil, ArErr{"Runtime Error", "can't set for non object", v.line, v.path, v.code, true}
			}
		}
	}
	return ThrowOnNonLoop(resp, ArErr{})
}

func parseDelete(code UNPARSEcode, index int, lines []UNPARSEcode) (ArDelete, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	spacesplit := strings.SplitN(trim, " ", 2)
	name := strings.TrimSpace(spacesplit[1])
	if blockedVariableNames[name] {
		return ArDelete{}, false, ArErr{"Naming Error", "\"" + name + "\" is a reserved keyword", code.line, code.path, code.realcode, true}, 1
	}
	toset, success, err, i := translateVal(UNPARSEcode{code: name, realcode: code.realcode, line: code.line, path: code.path}, index, lines, 0)

	if !success {
		return ArDelete{}, false, err, i
	}
	return ArDelete{
		toset,
		code.line,
		code.code,
		code.path,
	}, true, ArErr{}, i
}

func runDelete(d ArDelete, stack stack, stacklevel int) (any, ArErr) {
	switch x := d.value.(type) {
	case accessVariable:
		for i := len(stack) - 1; i >= 0; i-- {
			callable, ok := stack[i].obj["__Contains__"]
			if !ok {
				continue
			}
			contains, err := builtinCall(callable, []any{x.Name})
			if err.EXISTS {
				return nil, err
			}
			if anyToBool(contains) {
				callable, ok := stack[i].obj["__deleteindex__"]
				if !ok {
					continue
				}
				return builtinCall(callable, []any{x.Name})
			}
		}
		return nil, ArErr{"Name Error", "variable \"" + x.Name + "\" does not exist", d.line, d.path, d.code, true}
	case ArMapGet:
		respp, err := runVal(x.VAL, stack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		if len(x.Args) != 1 {
			return nil, ArErr{"Runtime Error", "can't delete by slice", d.line, d.path, d.code, true}
		}
		key, err := runVal(x.Args[0], stack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		switch y := respp.(type) {
		case ArObject:
			if typeof(y) == "array" {
				return nil, ArErr{"Runtime Error", "can't delete from array", d.line, d.path, d.code, true}
			}
			if isUnhashable(key) {
				return nil, ArErr{"Runtime Error", "can't use unhashable type as map key: " + typeof(key), d.line, d.path, d.code, true}
			}
			delete(y.obj, key)
		default:
			return nil, ArErr{"Runtime Error", "can't delete for non map", d.line, d.path, d.code, true}
		}
	default:
		return nil, ArErr{"Runtime Error", "can't delete for non variable", d.line, d.path, d.code, true}
	}
	return nil, ArErr{}
}
