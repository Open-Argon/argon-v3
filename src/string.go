package main

import (
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
	return obj
}
