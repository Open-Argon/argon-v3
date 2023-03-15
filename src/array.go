package main

import "strings"

var arrayCompile = makeRegex(`( *)\[(.|\n)*\]( *)`)

type CreateArray struct {
	value []any
	line  int
	code  string
	path  string
}

func isArray(code UNPARSEcode) bool {
	return arrayCompile.MatchString(code.code)
}

func ArArray(arr []any) ArObject {
	val := ArObject{
		"array",
		anymap{
			"__value__": arr,
			"length":    len(arr),
		},
	}
	val.obj["remove"] = builtinFunc{
		"remove",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "number" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "argument must be a number",
					EXISTS:  true,
				}
			}
			if !args[0].(number).IsInt() {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "argument must be an integer",
					EXISTS:  true,
				}
			}
			num := int(args[0].(number).Num().Int64())
			if num < 0 || num >= len(arr) {
				return nil, ArErr{
					TYPE:    "IndexError",
					message: "index out of range",
					EXISTS:  true,
				}
			}
			arr = append(arr[:num], arr[num+1:]...)
			val.obj["length"] = len(arr)
			val.obj["__value__"] = arr
			return nil, ArErr{}
		}}
	val.obj["append"] = builtinFunc{
		"append",
		func(args ...any) (any, ArErr) {
			if len(args) == 0 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			arr = append(arr, args...)
			val.obj["length"] = len(arr)
			val.obj["__value__"] = arr
			return nil, ArErr{}
		},
	}
	val.obj["insert"] = builtinFunc{
		"insert",
		func(args ...any) (any, ArErr) {
			if len(args) < 2 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "number" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "argument must be a number",
					EXISTS:  true,
				}
			}
			if !args[0].(number).IsInt() {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "argument must be an integer",
					EXISTS:  true,
				}
			}
			num := int(args[0].(number).Num().Int64())
			if num < 0 || num > len(arr) {
				return nil, ArErr{
					TYPE:    "IndexError",
					message: "index out of range",
					EXISTS:  true,
				}
			}
			arr = append(arr[:num], append(args[1:], arr[num:]...)...)
			val.obj["length"] = len(arr)
			val.obj["__value__"] = arr
			return nil, ArErr{}
		},
	}
	val.obj["pop"] = builtinFunc{
		"pop",
		func(args ...any) (any, ArErr) {
			if len(args) > 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "too many arguments",
					EXISTS:  true,
				}
			}
			if len(args) == 1 {
				if typeof(args[0]) != "number" {
					return nil, ArErr{
						TYPE:    "TypeError",
						message: "argument must be a number",
						EXISTS:  true,
					}
				}
				if !args[0].(number).IsInt() {
					return nil, ArErr{
						TYPE:    "TypeError",
						message: "argument must be an integer",
						EXISTS:  true,
					}
				}
				num := int(args[0].(number).Num().Int64())
				if num < 0 || num >= len(arr) {
					return nil, ArErr{
						TYPE:    "IndexError",
						message: "index out of range",
						EXISTS:  true,
					}
				}
				v := arr[num]
				arr = append(arr[:num], arr[num+1:]...)
				val.obj["length"] = len(arr)
				val.obj["__value__"] = arr
				return v, ArErr{}
			}
			v := arr[len(arr)-1]
			arr = arr[:len(arr)-1]
			val.obj["length"] = len(arr)
			val.obj["__value__"] = arr
			return v, ArErr{}
		},
	}
	val.obj["clear"] = builtinFunc{
		"clear",
		func(args ...any) (any, ArErr) {
			if len(args) != 0 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "too many arguments",
					EXISTS:  true,
				}
			}
			arr = []any{}
			val.obj["length"] = len(arr)
			val.obj["__value__"] = arr
			return nil, ArErr{}
		},
	}
	val.obj["extend"] = builtinFunc{
		"extend",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "array" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "argument must be an array",
					EXISTS:  true,
				}
			}
			arr = append(arr, args[0].(ArObject).obj["__value__"].([]any)...)
			val.obj["length"] = len(arr)
			val.obj["__value__"] = arr
			return nil, ArErr{}
		},
	}
	val.obj["map"] = builtinFunc{
		"map",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "function" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "argument must be a function",
					EXISTS:  true,
				}
			}
			newarr := []any{}
			for _, v := range arr {
				vv, err := runCall(call{
					args[0],
					[]any{v}, "", 0, "",
				}, stack{vars, newscope()}, 0)
				if err.EXISTS {
					return nil, err
				}
				newarr = append(newarr, vv)
			}
			return ArArray(newarr), ArErr{}
		},
	}
	val.obj["filter"] = builtinFunc{
		"filter",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "function" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "argument must be a function",
					EXISTS:  true,
				}
			}
			newarr := []any{}
			for _, v := range arr {
				vv, err := runCall(call{
					args[0],
					[]any{v}, "", 0, "",
				}, stack{vars, newscope()}, 0)
				if err.EXISTS {
					return nil, err
				}
				if anyToBool(vv) {
					newarr = append(newarr, v)
				}
			}
			return ArArray(newarr), ArErr{}
		},
	}
	val.obj["reduce"] = builtinFunc{
		"reduce",
		func(args ...any) (any, ArErr) {
			if len(args) != 2 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "function" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "argument must be a function",
					EXISTS:  true,
				}
			}
			if len(arr) == 0 {
				return nil, ArErr{
					TYPE:    "ValueError",
					message: "array is empty",
					EXISTS:  true,
				}
			}
			v := args[1]
			for _, vv := range arr {
				var err ArErr
				v, err = runCall(call{
					args[0],
					[]any{v, vv}, "", 0, "",
				}, stack{vars, newscope()}, 0)
				if err.EXISTS {
					return nil, err
				}
			}
			return v, ArErr{}
		},
	}
	val.obj["join"] = builtinFunc{
		"join",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "string" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "argument must be a string",
					EXISTS:  true,
				}
			}
			output := []string{}
			for _, v := range arr {
				if typeof(v) != "string" {
					return nil, ArErr{
						TYPE:    "TypeError",
						message: "array must be an array of strings",
						EXISTS:  true,
					}
				}
				output = append(output, v.(string))
			}
			return strings.Join(output, args[0].(string)), ArErr{}
		},
	}
	val.obj["concat"] = builtinFunc{
		"concat",
		func(args ...any) (any, ArErr) {
			if len(args) < 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "missing argument(s)",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "array" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "argument must be an array",
					EXISTS:  true,
				}
			}
			newarr := append(arr, args[0].(ArObject).obj["__value__"].([]any)...)
			return ArArray(newarr), ArErr{}
		},
	}
	return val
}

func potentialAnyArrayToArArray(arr any) any {
	switch arr := arr.(type) {
	case []any:
		return ArArray(arr)
	default:
		return arr
	}
}

func parseArray(code UNPARSEcode, index int, codelines []UNPARSEcode) (any, bool, ArErr, int) {
	trimmed := strings.TrimSpace(code.code)
	trimmed = trimmed[1 : len(trimmed)-1]
	arguments, worked, err := getValuesFromLetter(trimmed, ",", index, codelines, true)
	return CreateArray{
		value: arguments,
		line:  code.line,
		code:  code.realcode,
		path:  code.path,
	}, worked, err, 1
}

func runArray(a CreateArray, stack stack, stacklevel int) (ArObject, ArErr) {
	var array []any
	for _, val := range a.value {
		val, err := runVal(val, stack, stacklevel+1)
		if err.EXISTS {
			return ArObject{}, err
		}
		array = append(array, val)
	}
	return ArArray(array), ArErr{}
}
