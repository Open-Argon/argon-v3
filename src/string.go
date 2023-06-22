package main

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	return output, nil
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
		anymap{
			"__value__": str,
			"__name__":  "string",
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
			if len(a) > 3 {
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
				v = ArValidToAny(v)
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
			if typeof(a[0]) != "array" {
				return nil, ArErr{"TypeError", "expected array, got " + typeof(a[0]), 0, "", "", true}
			}
			output := []string{str}
			for _, v := range a[0].([]any) {
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
			splitby := ArValidToAny(a[0]).(string)
			output := []any{}
			splitted := (strings.Split(str, splitby))
			for _, v := range splitted {
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
			return cases.Title(language.English).String(str), ArErr{}
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
			a[0] = ArValidToAny(a[0])
			a[1] = ArValidToAny(a[1])
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
			return strings.Contains(str, a[0].(ArObject).obj["__value__"].(string)), ArErr{}
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
			return strings.HasPrefix(str, a[0].(ArObject).obj["__value__"].(string)), ArErr{}
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
			return strings.HasSuffix(str, a[0].(ArObject).obj["__value__"].(string)), ArErr{}
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
			return strings.Index(str, a[0].(ArObject).obj["__value__"].(string)), ArErr{}
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
			return strings.LastIndex(str, a[0].(ArObject).obj["__value__"].(string)), ArErr{}
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
			return strings.Count(str, a[0].(ArObject).obj["__value__"].(string)), ArErr{}
		}}

	obj.obj["sort"] = builtinFunc{
		"sort",
		func(args ...any) (any, ArErr) {
			if len(args) > 2 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "too many arguments",
					EXISTS:  true,
				}
			}
			reverse := false
			if len(args) >= 1 {
				if typeof(args[0]) != "boolean" {
					return nil, ArErr{
						TYPE:    "TypeError",
						message: "argument must be a boolean",
						EXISTS:  true,
					}
				}
				reverse = args[0].(bool)
			}
			bytes := []byte(str)
			anyarr := make([]any, len(bytes))
			for i, b := range bytes {
				anyarr[i] = b
			}
			if len(args) == 2 {
				if typeof(args[1]) != "function" {
					return nil, ArErr{
						TYPE:    "TypeError",
						message: "argument must be a function",
						EXISTS:  true,
					}
				}
				output, err := quickSort(anyarr, func(a any) (any, ArErr) {
					return runCall(call{
						args[1],
						[]any{a}, "", 0, "",
					}, stack{}, 0)
				})
				if err.EXISTS {
					return nil, err
				}
				bytes = make([]byte, len(output))
				for i, b := range output {
					bytes[i] = b.(byte)
				}
				str = string(bytes)
				obj.obj["length"] = len(str)
				obj.obj["__value__"] = str
				return nil, ArErr{}
			}
			output, err := quickSort(anyarr, func(a any) (any, ArErr) {
				return a, ArErr{}
			})
			if err.EXISTS {
				return nil, err
			}
			if reverse {
				for i, j := 0, len(output)-1; i < j; i, j = i+1, j-1 {
					output[i], output[j] = output[j], output[i]
				}
			}
			bytes = make([]byte, len(output))
			for i, b := range output {
				bytes[i] = b.(byte)
			}
			str = string(bytes)
			obj.obj["length"] = len(str)
			obj.obj["__value__"] = str
			return nil, ArErr{}
		},
	}
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
				cutset = a[0].(ArObject).obj["__value__"].(string)
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
				cutset = a[0].(ArObject).obj["__value__"].(string)
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
				cutset = a[0].(ArObject).obj["__value__"].(string)
			}
			return strings.TrimRight(str, cutset), ArErr{}
		}}
	obj.obj["__LessThanEqual__"] = builtinFunc{
		"__LessThanOrEqual__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "cannot get less than or equal to of type " + typeof(a[0]) + " from string", 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			return str <= a[0].(string), ArErr{}
		}}
	obj.obj["__LessThan__"] = builtinFunc{
		"__LessThan__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "cannot get less than of type " + typeof(a[0]) + " from string", 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			return str < a[0].(string), ArErr{}
		}}
	obj.obj["__GreaterThan__"] = builtinFunc{
		"__GreaterThan__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "cannot get greater than of type " + typeof(a[0]) + " from string", 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			return str > a[0].(string), ArErr{}
		}}

	obj.obj["__GreaterThanEqual__"] = builtinFunc{
		"__GreaterThanEqual__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "cannot get greater than or equal to of type " + typeof(a[0]) + " from string", 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			return str >= a[0].(string), ArErr{}
		}}
	obj.obj["__Equal__"] = builtinFunc{
		"__Equal__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			return str == a[0], ArErr{}
		}}
	obj.obj["__NotEqual__"] = builtinFunc{
		"__NotEqual__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			return str != a[0], ArErr{}
		}}
	obj.obj["__Add__"] = builtinFunc{
		"__Add__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			if typeof(a[0]) != "string" {
				a[0] = anyToArgon(a[0], false, false, 3, 0, false, 0)
			}
			return strings.Join([]string{str, a[0].(string)}, ""), ArErr{}
		}}
	obj.obj["__PostAdd__"] = builtinFunc{
		"__PostAdd__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			if typeof(a[0]) != "string" {
				a[0] = anyToArgon(a[0], false, false, 3, 0, false, 0)
			}
			return strings.Join([]string{a[0].(string), str}, ""), ArErr{}
		}}
	obj.obj["__Multiply__"] = builtinFunc{
		"__Multiply__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{"TypeError", "cannot multiply string by " + typeof(a[0]), 0, "", "", true}
			}
			n := a[0].(number)
			if !n.IsInt() {
				return nil, ArErr{"ValueError", "cannot multiply string by float", 0, "", "", true}
			}
			if n.Sign() < 0 {
				return nil, ArErr{"ValueError", "cannot multiply string by negative number", 0, "", "", true}
			}
			return strings.Repeat(str, int(n.Num().Int64())), ArErr{}
		}}
	obj.obj["__Contains__"] = builtinFunc{
		"__Contains__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "cannot check if string contains " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			return strings.Contains(str, a[0].(string)), ArErr{}
		}}
	obj.obj["__Subtract__"] = builtinFunc{
		"__Subtract__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "cannot subtract " + typeof(a[0]) + " from string", 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			return strings.Replace(str, a[0].(string), "", -1), ArErr{}
		}}
	obj.obj["__Divide__"] = builtinFunc{
		"__Divide__",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{"TypeError", "expected 1 argument, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{"TypeError", "cannot divide string by " + typeof(a[0]), 0, "", "", true}
			}
			a[0] = ArValidToAny(a[0])
			splitby := a[0].(string)
			output := []any{}
			splitted := (strings.Split(str, splitby))
			for _, v := range splitted {
				output = append(output, ArString(v))
			}
			return ArArray(output), ArErr{}
		}}
	obj.obj["__Boolean__"] = builtinFunc{
		"__Boolean__",
		func(a ...any) (any, ArErr) {
			return len(str) > 0, ArErr{}
		}}
	return obj
}
