package main

import (
	"fmt"
	"math/big"
	"strings"
)

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
		anymap{
			"__name__":  "array",
			"__value__": arr,
		},
	}
	val.obj["__json__"] = builtinFunc{
		"__json__",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return "", ArErr{
					TYPE:    "Type Error",
					message: "expected 1 argument, got " + fmt.Sprint(len(args)),
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "number" {
				return "", ArErr{
					TYPE:    "Type Error",
					message: "expected number, got " + typeof(args[0]),
					EXISTS:  true,
				}
			}

			output := []string{}

			level, err := numberToInt64(args[0].(ArObject))
			if err != nil {
				return "", ArErr{
					TYPE:    "Runtime Error",
					message: err.Error(),
					EXISTS:  true,
				}
			}
			for _, value := range arr {
				str, err := jsonstringify(value, level+1)
				if err != nil {
					return "", ArErr{
						TYPE:    "Runtime Error",
						message: err.Error(),
						EXISTS:  true,
					}
				}
				output = append(output, str)
			}
			return "[" + strings.Join(output, ", ") + "]", ArErr{}
		},
	}
	val.obj["__setindex__"] = builtinFunc{
		"__setindex__",
		func(a ...any) (any, ArErr) {
			if len(a) != 2 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "expected 2 arguments, got " + fmt.Sprint(len(a)),
					EXISTS:  true,
				}
			}
			if typeof(a[0]) != "number" {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "dex must be a number",
					EXISTS:  true,
				}
			}
			if !isNumberInt(a[0].(ArObject)) {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "index must be an integer",
					EXISTS:  true,
				}
			}
			num64, err := numberToInt64(a[0].(ArObject))
			if err != nil {
				return nil, ArErr{
					TYPE:    "Runtime Error",
					message: err.Error(),
					EXISTS:  true,
				}
			}
			num := int(num64)
			if num < 0 || num >= len(arr) {
				return nil, ArErr{
					TYPE:    "Index Error",
					message: "index out of range",
					EXISTS:  true,
				}
			}
			arr[num] = a[1]
			return nil, ArErr{}
		},
	}
	val.obj["__getindex__"] = builtinFunc{
		"__getindex__",
		func(a ...any) (any, ArErr) {
			// a[0] is start
			// a[1] is end
			// a[2] is step
			if len(a) > 3 || len(a) == 0 {
				return nil, ArErr{"Type Error", "expected 1 to 3 arguments, got " + fmt.Sprint(len(a)), 0, "", "", true}
			}
			{
				if len(a) == 1 {
					if typeof(a[0]) == "string" {
						var name = ArValidToAny(a[0]).(string)
						if name == "length" {
							return Number(compiledNumber{big.NewInt(int64(len(arr)))}), ArErr{}
						}
					}
				}
			}
			var (
				start int = 0
				end   any = nil
				step  int = 1
			)
			{
				if a[0] == nil {
					start = 0
				} else if typeof(a[0]) != "number" || !isNumberInt(a[0].(ArObject)) {
					return "", ArErr{
						TYPE:    "Type Error",
						message: "slice index must be an integer",
						EXISTS:  true,
					}
				} else {
					start64, err := numberToInt64(a[0].(ArObject))
					if err != nil {
						return "", ArErr{
							TYPE:    "Runtime Error",
							message: err.Error(),
							EXISTS:  true,
						}
					}
					start = int(start64)
				}
			}
			if len(a) > 1 {
				if a[1] == nil {
					end = len(arr)
				} else if typeof(a[1]) != "number" || !isNumberInt(a[1].(ArObject)) {
					return "", ArErr{
						TYPE:    "Type Error",
						message: "slice index must be an integer",
						EXISTS:  true,
					}
				} else {
					end64, err := numberToInt64(a[1].(ArObject))
					if err != nil {
						return "", ArErr{
							TYPE:    "Runtime Error",
							message: err.Error(),
							EXISTS:  true,
						}
					}
					end = int(end64)
				}
			}
			if len(a) > 2 {
				if a[2] == nil {
					step = 1
				} else if typeof(a[2]) != "number" || !isNumberInt(a[2].(ArObject)) {
					return "", ArErr{
						TYPE:    "Type Error",
						message: "slice index must be an integer",
						EXISTS:  true,
					}
				} else {
					step64, err := numberToInt64(a[2].(ArObject))
					if err != nil {
						return "", ArErr{
							TYPE:    "Runtime Error",
							message: err.Error(),
							EXISTS:  true,
						}
					}
					step = int(step64)
				}
			}
			var ogStart = start
			if start < 0 {
				start = len(arr) + start
			}
			if _, ok := end.(int); ok && end.(int) < 0 {
				end = len(arr) + end.(int)
			}
			if end != nil && end.(int) > len(arr) {
				end = len(arr)
			}
			if start >= len(arr) || start < 0 {
				return "", ArErr{
					TYPE:    "Index Error",
					message: "index out of range, trying to access index " + fmt.Sprint(ogStart) + " in array of length " + fmt.Sprint(len(arr)),
					EXISTS:  true,
				}
			}
			if end == nil {
				return arr[start], ArErr{}
			} else if step == 1 {
				return arr[start:end.(int)], ArErr{}
			} else {
				output := []any{}
				if step > 0 {
					for i := start; i < end.(int); i += step {
						output = append(output, arr[i])
					}
				} else {
					for i := end.(int) - 1; i >= start; i += step {
						output = append(output, arr[i])
					}
				}
				return output, ArErr{}
			}
		}}
	val.obj["remove"] = builtinFunc{
		"remove",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "number" {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "argument must be a number",
					EXISTS:  true,
				}
			}
			if !isNumberInt(args[0].(ArObject)) {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "argument must be an integer",
					EXISTS:  true,
				}
			}
			num64, err := (numberToInt64(args[0].(ArObject)))
			if err != nil {
				return nil, ArErr{
					TYPE:    "Runtime Error",
					message: err.Error(),
					EXISTS:  true,
				}
			}
			num := int(num64)
			if num < 0 || num >= len(arr) {
				return nil, ArErr{
					TYPE:    "Index Error",
					message: "index out of range",
					EXISTS:  true,
				}
			}
			arr = append(arr[:num], arr[num+1:]...)
			val.obj["__value__"] = arr
			return nil, ArErr{}
		}}
	val.obj["append"] = builtinFunc{
		"append",
		func(args ...any) (any, ArErr) {
			if len(args) == 0 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			arr = append(arr, args...)
			val.obj["__value__"] = arr
			return nil, ArErr{}
		},
	}
	val.obj["insert"] = builtinFunc{
		"insert",
		func(args ...any) (any, ArErr) {
			if len(args) < 2 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "number" {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "argument must be a number",
					EXISTS:  true,
				}
			}
			if !isNumberInt(args[0].(ArObject)) {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "argument must be an integer",
					EXISTS:  true,
				}
			}
			num64, err := numberToInt64(args[0].(ArObject))
			if err != nil {
				return nil, ArErr{
					TYPE:    "Runtime Error",
					message: err.Error(),
					EXISTS:  true,
				}
			}
			num := int(num64)
			if num < 0 || num > len(arr) {
				return nil, ArErr{
					TYPE:    "Index Error",
					message: "index out of range",
					EXISTS:  true,
				}
			}
			arr = append(arr[:num], append(args[1:], arr[num:]...)...)
			val.obj["__value__"] = arr
			return nil, ArErr{}
		},
	}
	val.obj["pop"] = builtinFunc{
		"pop",
		func(args ...any) (any, ArErr) {
			if len(args) > 1 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "too many arguments",
					EXISTS:  true,
				}
			}
			if len(args) == 1 {
				if typeof(args[0]) != "number" {
					return nil, ArErr{
						TYPE:    "Type Error",
						message: "argument must be a number",
						EXISTS:  true,
					}
				}
				if !isNumberInt(args[0].(ArObject)) {
					return nil, ArErr{
						TYPE:    "Type Error",
						message: "argument must be an integer",
						EXISTS:  true,
					}
				}
				num64, err := (numberToInt64(args[0].(ArObject)))
				if err != nil {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: err.Error(),
						EXISTS:  true,
					}
				}
				num := int(num64)
				if num < 0 || num >= len(arr) {
					return nil, ArErr{
						TYPE:    "Index Error",
						message: "index out of range",
						EXISTS:  true,
					}
				}
				v := arr[num]
				arr = append(arr[:num], arr[num+1:]...)
				val.obj["__value__"] = arr
				return v, ArErr{}
			}
			v := arr[len(arr)-1]
			arr = arr[:len(arr)-1]
			val.obj["__value__"] = arr
			return v, ArErr{}
		},
	}
	val.obj["clear"] = builtinFunc{
		"clear",
		func(args ...any) (any, ArErr) {
			if len(args) != 0 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "too many arguments",
					EXISTS:  true,
				}
			}
			arr = []any{}
			val.obj["__value__"] = arr
			return nil, ArErr{}
		},
	}
	val.obj["extend"] = builtinFunc{
		"extend",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "array" {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "argument must be an array",
					EXISTS:  true,
				}
			}
			arr = append(arr, args[0].(ArObject).obj["__value__"].([]any)...)
			val.obj["__value__"] = arr
			return nil, ArErr{}
		},
	}
	val.obj["sort"] = builtinFunc{
		"sort",
		func(args ...any) (any, ArErr) {
			if len(args) > 2 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "too many arguments",
					EXISTS:  true,
				}
			}
			reverse := false
			if len(args) >= 1 {
				if typeof(args[0]) != "boolean" {
					return nil, ArErr{
						TYPE:    "Type Error",
						message: "argument must be a boolean",
						EXISTS:  true,
					}
				}
				reverse = args[0].(bool)
			}
			if len(args) == 2 {
				if typeof(args[1]) != "function" {
					return nil, ArErr{
						TYPE:    "Type Error",
						message: "argument must be a function",
						EXISTS:  true,
					}
				}
				output, err := quickSort(arr, func(a any) (any, ArErr) {
					return runCall(call{
						args[1],
						[]any{a}, "", 0, "",
					}, stack{}, 0)
				})
				if err.EXISTS {
					return nil, err
				}
				if reverse {
					for i, j := 0, len(output)-1; i < j; i, j = i+1, j-1 {
						output[i], output[j] = output[j], output[i]
					}
				}
				arr = output
				val.obj["__value__"] = arr
				return nil, ArErr{}
			}
			output, err := quickSort(arr, func(a any) (any, ArErr) {
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
			arr = output
			val.obj["__value__"] = arr
			return nil, ArErr{}
		},
	}
	val.obj["map"] = builtinFunc{
		"map",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "function" {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "argument must be a function",
					EXISTS:  true,
				}
			}
			newarr := []any{}
			for _, v := range arr {
				vv, err := runCall(call{
					args[0],
					[]any{v}, "", 0, "",
				}, stack{}, 0)
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
					TYPE:    "Type Error",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "function" {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "argument must be a function",
					EXISTS:  true,
				}
			}
			newarr := []any{}
			for _, v := range arr {
				vv, err := runCall(call{
					args[0],
					[]any{v}, "", 0, "",
				}, stack{}, 0)
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
					TYPE:    "Type Error",
					message: "missing argument, expected 2 got " + fmt.Sprint(len(args)),
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "function" {
				return nil, ArErr{
					TYPE:    "Type Error",
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
				}, stack{}, 0)
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
					TYPE:    "Type Error",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "string" {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "argument must be a string",
					EXISTS:  true,
				}
			}
			output := []string{}
			for _, v := range arr {
				if typeof(v) != "string" {
					return nil, ArErr{
						TYPE:    "Type Error",
						message: "array must be an array of strings",
						EXISTS:  true,
					}
				}
				output = append(output, v.(ArObject).obj["__value__"].(string))
			}
			return ArString(strings.Join(output, args[0].(ArObject).obj["__value__"].(string))), ArErr{}
		},
	}
	val.obj["concat"] = builtinFunc{
		"concat",
		func(args ...any) (any, ArErr) {
			if len(args) < 1 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "missing argument(s)",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "array" {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "argument must be an array",
					EXISTS:  true,
				}
			}
			newarr := append(arr, args[0].(ArObject).obj["__value__"].([]any)...)
			return ArArray(newarr), ArErr{}
		},
	}
	val.obj["__Equal__"] = builtinFunc{
		"__Equal__",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "array" {
				return false, ArErr{}
			}
			if len(arr) != len(args[0].(ArObject).obj["__value__"].([]any)) {
				return false, ArErr{}
			}
			for i, v := range arr {
				res, err := runOperation(operationType{
					operation: 8,
					values:    []any{v, args[0].(ArObject).obj["__value__"].([]any)[i]},
				}, stack{}, 0)
				if err.EXISTS {
					return nil, err
				}
				if anyToBool(res) {
					return false, ArErr{}
				}
			}
			return true, ArErr{}
		}}
	val.obj["__Contains__"] = builtinFunc{
		"__Contains__",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			for _, v := range arr {
				res, err := runOperation(operationType{
					operation: 9,
					values:    []any{v, args[0]},
				}, stack{}, 0)
				if err.EXISTS {
					return nil, err
				}
				if anyToBool(res) {
					return true, ArErr{}
				}
			}
			return false, ArErr{}
		},
	}
	val.obj["__NotContains__"] = builtinFunc{
		"__NotContains__",
		func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "missing argument",
					EXISTS:  true,
				}
			}
			for _, v := range arr {
				res, err := runOperation(operationType{
					operation: 9,
					values:    []any{v, args[0]},
				}, stack{}, 0)
				if err.EXISTS {
					return nil, err
				}
				if anyToBool(res) {
					return false, ArErr{}
				}
			}
			return true, ArErr{}
		},
	}
	val.obj["copy"] = builtinFunc{
		"copy",
		func(args ...any) (any, ArErr) {
			arrCopy := make([]any, len(arr))
			copy(arrCopy, arr)
			return ArArray(arrCopy), ArErr{}
		},
	}
	val.obj["__Boolean__"] = builtinFunc{
		"__Boolean__",
		func(args ...any) (any, ArErr) {
			return len(
				arr,
			) > 0, ArErr{}
		},
	}
	return val
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
