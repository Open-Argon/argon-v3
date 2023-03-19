package main

import (
	"fmt"
	"strconv"
	"strings"
)

var stringCompile = makeRegex("(( *)\"((\\\\([a-z\\\"'`]))|[^\\\"])*\"( *))|(( *)'((\\\\([a-z\\'\"`]))|[^\\'])*'( *))")

func isString(code UNPARSEcode) bool {
	return stringCompile.MatchString(code.code)
}

func unquoted(
	str string,
) (string, error) {
	str = strings.Trim(str, " ")
	if str[0] == '\'' {
		str = strings.Replace(str, "\\\"", "\"", -1)
		str = strings.Replace(str, "\"", "\\\"", -1)
	}
	str = str[1 : len(str)-1]
	str = strings.Replace(str, "\\'", "'", -1)
	str = "\"" + str + "\""
	output, err := strconv.Unquote(str)
	if err != nil {
		return "", err
	}
	classoutput := (output)
	return classoutput, nil
}

// returns translateString, success, error
func parseString(code UNPARSEcode) (string, bool, ArErr, int) {
	trim := strings.Trim(code.code, " ")

	unquoted, err := unquoted(trim)
	if err != nil {
		return "", false, ArErr{"Syntax Error", "invalid string", code.line, code.path, code.realcode, true}, 1
	}

	return unquoted, true, ArErr{}, 1
}

func ArString(str string) ArObject {
	obj := ArObject{
		"string",
		anymap{
			"__value__": str,
			"length":    newNumber().SetUint64(uint64(len(str))),
		},
	}

	obj.obj["__setindex__"] = builtinFunc{
		"__setindex__",
		func(a ...any) (any, ArErr) {
			if len(a) != 2 {
				return nil, ArErr{"TypeError", "expected 2 arguments, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"TypeError", "expected number, got " + typeof(a[0]), 0, "", "", true}
			}
			if typeof(a[1]) != "string" {
				return nil, ArErr{"TypeError", "expected string, got " + typeof(a[1]), 0, "", "", true}
			}
			if len(a[1].(string)) != 1 {
				return nil, ArErr{"TypeError", "expected string of length 1, got " + fmt.Sprint(len(a[1].(string))), 0, "", "", true}
			}
			if !a[0].(number).IsInt() {
				return nil, ArErr{"TypeError", "expected integer, got float", 0, "", "", true}
			}
			index := a[0].(number).Num().Int64()
			if index < 0 {
				index = int64(len(str)) + index
			}
			if index < 0 || index >= int64(len(str)) {
				return nil, ArErr{"IndexError", "index out of range", 0, "", "", true}
			}
			str = strings.Join([]string{str[:index], a[1].(string), str[index+1:]}, "")
			obj.obj["__value__"] = str
			obj.obj["length"] = newNumber().SetUint64(uint64(len(str)))
			return nil, ArErr{}
		}}
	obj.obj["__getindex__"] = builtinFunc{
		"__getindex__",
		func(a ...any) (any, ArErr) {
			// a[0] is start
			// a[1] is end
			// a[2] is step
			if len(a) < 0 || len(a) > 3 {
				return nil, ArErr{"TypeError", "expected 1 to 3 arguments, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			var (
				start int = 0
				end   any = nil
				step  int = 1
			)
			{
				if a[0] == nil {
					start = 0
				} else if typeof(a[0]) != "number" || !a[0].(number).IsInt() {
					return "", ArErr{
						TYPE:    "TypeError",
						message: "slice index must be an integer",
						EXISTS:  true,
					}
				} else {
					start = int(a[0].(number).Num().Int64())
				}
			}
			if len(a) > 1 {
				if a[1] == nil {
					end = len(str)
				} else if typeof(a[1]) != "number" || !a[1].(number).IsInt() {
					return "", ArErr{
						TYPE:    "TypeError",
						message: "slice index must be an integer",
						EXISTS:  true,
					}
				} else {
					end = int(a[1].(number).Num().Int64())
				}
			}
			if len(a) > 2 {
				if a[2] == nil {
					step = 1
				} else if typeof(a[2]) != "number" || !a[2].(number).IsInt() {
					return "", ArErr{
						TYPE:    "TypeError",
						message: "slice index must be an integer",
						EXISTS:  true,
					}
				} else {
					step = int(a[2].(number).Num().Int64())
				}
			}
			if start < 0 {
				start = len(str) + start
			}
			if _, ok := end.(int); ok && end.(int) < 0 {
				end = len(str) + end.(int)
			}

			if end != nil && end.(int) > len(str) {
				end = len(str)
			}
			if end == nil {
				return string(str[start]), ArErr{}
			} else if step == 1 {
				return str[start:end.(int)], ArErr{}
			} else {
				output := []byte{}
				if step > 0 {
					for i := start; i < end.(int); i += step {
						output = append(output, str[i])
					}
				} else {
					for i := end.(int) - 1; i >= start; i += step {
						output = append(output, str[i])
					}
				}
				return string(output), ArErr{}
			}
		}}
	obj.obj["append"] = builtinFunc{
		"append",
		func(a ...any) (any, ArErr) {
			if len(a) == 0 {
				return nil, ArErr{"TypeError", "expected 1 or more argument, got 0", 0, "", "", true}
			}
			output := []string{str}
			for _, v := range a {
				if typeof(v) != "string" {
					return nil, ArErr{"TypeError", "expected string, got " + typeof(v), 0, "", "", true}
				}
				output = append(output, v.(string))
			}
			str = strings.Join(output, "")
			obj.obj["__value__"] = str
			obj.obj["length"] = newNumber().SetUint64(uint64(len(str)))
			return nil, ArErr{}
		}}
	obj.obj["extend"] = builtinFunc{
		"extend",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			output := []string{str}
			for _, v := range a {
				if typeof(v) != "string" {
					return nil, ArErr{"TypeError", "expected string, got " + typeof(v), 0, "", "", true}
				}
				output = append(output, v.(string))
			}
			str = strings.Join(output, "")
			obj.obj["__value__"] = str
			obj.obj["length"] = newNumber().SetUint64(uint64(len(str)))
			return nil, ArErr{}
		},
	}
	obj.obj["insert"] = builtinFunc{
		"insert",
		func(a ...any) (any, ArErr) {
			if len(a) != 2 {
				return nil, ArErr{"TypeError", "expected 2 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" || !a[0].(number).IsInt() {
				return nil, ArErr{"TypeError", "expected integer, got " + typeof(a[0]), 0, "", "", true}
			}
			if typeof(a[1]) != "string" {
				return nil, ArErr{"TypeError", "expected string, got " + typeof(a[1]), 0, "", "", true}
			}
			index := int(a[0].(number).Num().Int64())
			if index < 0 {
				index = len(str) + index
			}
			if index > len(str) {
				index = len(str)
			}
			str = str[:index] + a[1].(string) + str[index:]
			obj.obj["__value__"] = str
			obj.obj["length"] = newNumber().SetUint64(uint64(len(str)))
			return nil, ArErr{}
		}}
	obj.obj["concat"] = builtinFunc{
		"concat",
		func(a ...any) (any, ArErr) {
			if len(a) == 0 {
				return nil, ArErr{"TypeError", "expected 1 or more argument, got 0", 0, "", "", true}
			}
			output := []string{str}
			for _, v := range a {
				if typeof(v) != "string" {
					return nil, ArErr{"TypeError", "expected string, got " + typeof(v), 0, "", "", true}
				}
				output = append(output, v.(string))
			}
			return strings.Join(output, ""), ArErr{}
		}}
	obj.obj["split"] = builtinFunc{
		"split",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 or more argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "expected string, got " + typeof(a[0]), 0, "", "", true}
			}
			splitby := a[0].(string)
			output := []any{}
			splitted := any(strings.Split(str, splitby))
			for _, v := range splitted.([]string) {
				output = append(output, ArString(v))
			}
			return output, ArErr{}
		}}
	obj.obj["capitalise"] = builtinFunc{
		"capitalise",
		func(a ...any) (any, ArErr) {
			if len(a) != 0 {
				return nil, ArErr{"TypeError", "expected 0 arguments, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			return strings.Title(str), ArErr{}
		}}
	obj.obj["lower"] = builtinFunc{
		"lower",
		func(a ...any) (any, ArErr) {
			if len(a) != 0 {
				return nil, ArErr{"TypeError", "expected 0 arguments, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			return strings.ToLower(str), ArErr{}
		}}
	obj.obj["upper"] = builtinFunc{
		"upper",
		func(a ...any) (any, ArErr) {
			if len(a) != 0 {
				return nil, ArErr{"TypeError", "expected 0 arguments, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			return strings.ToUpper(str), ArErr{}
		}}
	obj.obj["replace"] = builtinFunc{
		"replace",
		func(a ...any) (any, ArErr) {
			if len(a) != 2 {
				return nil, ArErr{"TypeError", "expected 2 arguments, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "expected string, got " + typeof(a[0]), 0, "", "", true}
			}
			if typeof(a[1]) != "string" {
				return nil, ArErr{"TypeError", "expected string, got " + typeof(a[1]), 0, "", "", true}
			}
			return strings.Replace(str, a[0].(string), a[1].(string), -1), ArErr{}
		}}
	obj.obj["contains"] = builtinFunc{
		"contains",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "expected string, got " + typeof(a[0]), 0, "", "", true}
			}
			return strings.Contains(str, a[0].(string)), ArErr{}
		}}
	obj.obj["startswith"] = builtinFunc{
		"startswith",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "expected string, got " + typeof(a[0]), 0, "", "", true}
			}
			return strings.HasPrefix(str, a[0].(string)), ArErr{}
		}}
	obj.obj["endswith"] = builtinFunc{
		"endswith",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "expected string, got " + typeof(a[0]), 0, "", "", true}
			}
			return strings.HasSuffix(str, a[0].(string)), ArErr{}
		}}
	obj.obj["index"] = builtinFunc{
		"index",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "expected string, got " + typeof(a[0]), 0, "", "", true}
			}
			return strings.Index(str, a[0].(string)), ArErr{}
		}}
	obj.obj["rindex"] = builtinFunc{
		"rindex",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "expected string, got " + typeof(a[0]), 0, "", "", true}
			}
			return strings.LastIndex(str, a[0].(string)), ArErr{}
		}}
	obj.obj["count"] = builtinFunc{
		"count",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {

				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "expected string, got " + typeof(a[0]), 0, "", "", true}
			}
			return strings.Count(str, a[0].(string)), ArErr{}
		}}
	obj.obj["strip"] = builtinFunc{
		"strip",
		func(a ...any) (any, ArErr) {
			if len(a) > 1 {
				return nil, ArErr{"TypeError", "expected 0 or 1 arguments, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			cutset := " "
			if len(a) == 1 {
				if typeof(a[0]) != "string" {
					return nil, ArErr{"TypeError", "expected string, got " + typeof(a[0]), 0, "", "", true}
				}
				cutset = a[0].(string)
			}
			return strings.Trim(str, cutset), ArErr{}
		}}
	obj.obj["leftstrip"] = builtinFunc{
		"leftstrip",
		func(a ...any) (any, ArErr) {
			if len(a) > 1 {
				return nil, ArErr{"TypeError", "expected 0 or 1 arguments, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			cutset := " "
			if len(a) == 1 {
				if typeof(a[0]) != "string" {
					return nil, ArErr{"TypeError", "expected string, got " + typeof(a[0]), 0, "", "", true}
				}
				cutset = a[0].(string)
			}
			return strings.TrimLeft(str, cutset), ArErr{}
		}}
	obj.obj["rightstrip"] = builtinFunc{
		"rightstrip",
		func(a ...any) (any, ArErr) {
			if len(a) > 1 {
				return nil, ArErr{"TypeError", "expected 0 or 1 arguments, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			cutset := " "
			if len(a) == 1 {
				if typeof(a[0]) != "string" {
					return nil, ArErr{"TypeError", "expected string, got " + typeof(a[0]), 0, "", "", true}
				}
				cutset = a[0].(string)
			}
			return strings.TrimRight(str, cutset), ArErr{}
		}}
	return obj
}
