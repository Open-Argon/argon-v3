package main

import "strings"

func anyToBool(x any) bool {
	switch x := x.(type) {
	case string:
		return x != ""
	case number:
		return x.Cmp(newNumber()) != 0
	case bool:
		return x
	case nil:
		return false
	case ArObject:
		if y, ok := x.obj["__Boolean__"]; ok {
			val, err := runCall(
				call{
					callable: y,
					args:     []any{},
				}, stack{}, 0)
			if err.EXISTS {
				return false
			}
			return anyToBool(val)
		}
		return false
	case builtinFunc:
		return true
	case Callable:
		return true
	default:
		return true
	}
}

var booleanCompile = makeRegex(`( )*(true|false|null)( )*`)

func isBoolean(code UNPARSEcode) bool {
	return booleanCompile.MatchString(code.code)
}

func parseBoolean(code UNPARSEcode) (any, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	if trim == "true" {
		return true, true, ArErr{}, 1
	} else if trim == "false" {
		return false, true, ArErr{}, 1
	}
	return nil, true, ArErr{}, 1
}
