package main

import (
	"fmt"
	"os"
)

func makeGlobal() ArObject {
	var vars = anymap{}
	vars["global"] = vars
	vars["env"] = env
	vars["term"] = ArTerm
	vars["ArgonVersion"] = ArString(VERSION)
	vars["ArgonVersionNumber"] = Number(VERSION_NUM)
	vars["number"] = builtinFunc{"number", ArgonNumber}
	vars["string"] = builtinFunc{"string", ArgonString}
	vars["socket"] = Map(anymap{
		"server": builtinFunc{"server", ArSocketServer},
		"client": builtinFunc{"client", ArSocketClient},
	})
	vars["eval"] = builtinFunc{"eval", AReval}
	vars["map"] = builtinFunc{"map", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return Map(anymap{}), ArErr{}
		}
		switch x := a[0].(type) {
		case ArObject:
			if typeof(x) == "array" {
				newmap := anymap{}
				for i, v := range x.obj["__value__"].([]any) {
					v := ArValidToAny(v)
					switch y := v.(type) {
					case []any:
						if len(y) == 2 {
							if isUnhashable(y[0]) {
								return nil, ArErr{TYPE: "Type Error", message: "Cannot use unhashable value as key: " + typeof(y[0]), EXISTS: true}
							}
							key := ArValidToAny(y[0])
							newmap[key] = y[1]
							continue
						}
					}
					newmap[i] = v
				}
				return Map(newmap), ArErr{}
			} else if typeof(x) == "string" {
				newmap := anymap{}
				for i, v := range x.obj["__value__"].(string) {
					newmap[i] = ArString(string(v))
				}
				return Map(newmap), ArErr{}
			}

			newmap := anymap{}
			for key, val := range x.obj {
				newmap[key] = val
			}
			return Map(newmap), ArErr{}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot create map from '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["hex"] = builtinFunc{"hex", func(a ...any) (any, ArErr) {
		if len(a) != 1 {
			return nil, ArErr{TYPE: "Type Error", message: "expected 1 argument, got " + fmt.Sprint(len(a)), EXISTS: true}
		}
		a[0] = ArValidToAny(a[0])
		switch x := a[0].(type) {
		case int64:
			return ArString(fmt.Sprintf("%x", x)), ArErr{}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot convert '" + typeof(a[0]) + "' to hex", EXISTS: true}
	}}
	vars["buffer"] = builtinFunc{"buffer", func(a ...any) (any, ArErr) {
		if len(a) != 0 {
			return nil, ArErr{TYPE: "Type Error", message: "expected 0 arguments, got " + fmt.Sprint(len(a)), EXISTS: true}
		}
		return ArBuffer([]byte{}), ArErr{}
	}}
	vars["byte"] = builtinFunc{"byte", func(a ...any) (any, ArErr) {
		if len(a) != 0 {
			return nil, ArErr{TYPE: "Type Error", message: "expected 0 arguments, got " + fmt.Sprint(len(a)), EXISTS: true}
		}
		return ArByte(0), ArErr{}
	}}
	vars["throwError"] = builtinFunc{"throwError", ArThrowError}
	vars["array"] = builtinFunc{"array", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return ArArray([]any{}), ArErr{}
		}
		switch x := a[0].(type) {
		case ArObject:
			if typeof(x) == "array" {
				return x, ArErr{}
			} else if typeof(x) == "string" {

				newarray := []any{}
				for _, v := range x.obj["__value__"].(string) {
					newarray = append(newarray, ArString(string(v)))
				}
				return ArArray(newarray), ArErr{}
			}
			newarray := []any{}
			for key, val := range x.obj {
				newarray = append(newarray, []any{key, val})
			}
			return ArArray(newarray), ArErr{}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot create array from '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["boolean"] = builtinFunc{"boolean", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return false, ArErr{}
		}
		return anyToBool(a[0]), ArErr{}
	}}
	vars["time"] = ArTime
	vars["PI"] = PI
	vars["π"] = PI
	vars["e"] = e
	vars["ln"] = builtinFunc{"ln", ArgonLn}
	vars["log"] = builtinFunc{"log", ArgonLog}
	vars["logN"] = builtinFunc{"logN", ArgonLogN}
	vars["thread"] = builtinFunc{"thread", ArThread}
	vars["input"] = ArInput
	vars["round"] = builtinFunc{"round", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return nil, ArErr{TYPE: "round", message: "round takes 1 argument",
				EXISTS: true}
		}
		var precision int64 = 0
		if len(a) > 1 {
			a[1] = ArValidToAny(a[1])
			switch x := a[1].(type) {
			case int64:
				precision = x
			default:
				return nil, ArErr{TYPE: "Type Error", message: "Cannot round to '" + typeof(a[1]) + "'", EXISTS: true}
			}
		}
		switch x := a[0].(type) {
		case ArObject:
			if round_method, ok := x.obj["__round__"]; ok {
				return builtinCall(round_method, []any{Number(precision)})
			}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot round '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["floor"] = builtinFunc{"floor", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return nil, ArErr{TYPE: "floor", message: "floor takes 1 argument",
				EXISTS: true}
		}
		switch x := a[0].(type) {
		case ArObject:
			if floor_method, ok := x.obj["__floor__"]; ok {
				return builtinCall(floor_method, []any{})
			}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot floor '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["ceil"] = builtinFunc{"ceil", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return nil, ArErr{TYPE: "ceil", message: "ceil takes 1 argument",
				EXISTS: true}
		}
		switch x := a[0].(type) {
		case ArObject:
			if ceil_method, ok := x.obj["__ceil__"]; ok {
				return builtinCall(ceil_method, []any{})
			}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot ceil '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["sqrt"] = builtinFunc{"sqrt", ArgonSqrt}
	vars["file"] = ArFile
	vars["random"] = ArRandom
	vars["json"] = ArJSON
	vars["sin"] = ArSin
	vars["arcsin"] = ArArcsin
	vars["cos"] = ArCos
	vars["arccos"] = ArArccos
	vars["tan"] = ArTan
	vars["arctan"] = ArArctan
	vars["cosec"] = ArCosec
	vars["arccosec"] = ArArccosec
	vars["sec"] = ArSec
	vars["arcsec"] = ArArcsec
	vars["cot"] = ArCot
	vars["arccot"] = ArArccot
	vars["todeg"] = ArToDeg
	vars["colour"] = ArColour
	vars["torad"] = ArToRad
	vars["fraction"] = builtinFunc{"fraction", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return nil, ArErr{TYPE: "fraction", message: "fraction takes 1 argument",
				EXISTS: true}
		}
		switch x := a[0].(type) {
		case ArObject:
			if callable, ok := x.obj["__fraction__"]; ok {
				resp, err := runCall(
					call{
						Callable: callable,
						Args:     []any{},
					},
					stack{},
					0,
				)
				if err.EXISTS {
					return nil, err
				}
				return resp, ArErr{}
			}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot fraction '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["dir"] = builtinFunc{"dir", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return ArArray([]any{}), ArErr{}
		}
		t := AnyToArValid(a[0])
		switch x := t.(type) {
		case ArObject:
			newarray := []any{}
			for key := range x.obj {
				newarray = append(newarray, AnyToArValid(key))
			}
			return ArArray(newarray), ArErr{}
		}
		return ArArray([]any{}), ArErr{}
	}}
	vars["subprocess"] = builtinFunc{"subprocess", ArSubprocess}
	vars["sequence"] = builtinFunc{"sequence", ArSequence}
	vars["exit"] = builtinFunc{"exit", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			os.Exit(0)
		}
		a[0] = ArValidToAny(a[0])
		switch x := a[0].(type) {
		case int64:
			os.Exit(int(x))
			return nil, ArErr{}
		}
		os.Exit(0)
		return nil, ArErr{}
	}}
	vars["chr"] = builtinFunc{"chr", func(a ...any) (any, ArErr) {
		if len(a) != 1 {
			return nil, ArErr{TYPE: "chr", message: "chr takes 1 argument, got " + fmt.Sprint(len(a)), EXISTS: true}
		}
		a[0] = ArValidToAny(a[0])
		switch x := a[0].(type) {
		case int64:
			return string([]rune{rune(x)}), ArErr{}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot convert '" + typeof(a[0]) + "' to string", EXISTS: true}
	}}
	vars["ord"] = builtinFunc{"ord", func(a ...any) (any, ArErr) {
		if len(a) != 1 {
			return nil, ArErr{TYPE: "ord", message: "ord takes 1 argument, got " + fmt.Sprint(len(a)), EXISTS: true}
		}
		a[0] = ArValidToAny(a[0])
		switch x := a[0].(type) {
		case string:
			if len(x) != 1 {
				return nil, ArErr{TYPE: "ord", message: "ord takes a string with only one character, got " + fmt.Sprint(len(a)), EXISTS: true}
			}
			return Number(int64([]rune(x)[0])), ArErr{}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot convert '" + typeof(a[0]) + "' to string", EXISTS: true}
	}}
	vars["max"] = builtinFunc{"max", func(a ...any) (any, ArErr) {
		if len(a) != 1 {
			return nil, ArErr{TYPE: "runtime Error", message: "max takes 1 argument, got " + fmt.Sprint(len(a)), EXISTS: true}
		}
		a[0] = ArValidToAny(a[0])
		switch x := a[0].(type) {
		case []any:
			if len(x) == 0 {
				return nil, ArErr{TYPE: "runtime Error", message: "max takes a non-empty array", EXISTS: true}
			}
			var max ArObject
			for i, v := range x {
				switch x := v.(type) {
				case ArObject:
					if i == 0 {
						max = x
					} else {
						compared, err := CompareObjects(max, x)

						if err.EXISTS {
							return nil, err
						}

						compared_int, Err := numberToInt64(compared)

						if Err != nil {
							return nil, ArErr{TYPE: "Type Error", message: Err.Error(), EXISTS: true}
						}

						if compared_int == 1 {
							max = x
						}
					}
				}
			}
			return max, ArErr{}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot get max of type '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["min"] = builtinFunc{"min", func(a ...any) (any, ArErr) {
		if len(a) != 1 {
			return nil, ArErr{TYPE: "runtime Error", message: "max takes 1 argument, got " + fmt.Sprint(len(a)), EXISTS: true}
		}
		a[0] = ArValidToAny(a[0])
		switch x := a[0].(type) {
		case []any:
			if len(x) == 0 {
				return nil, ArErr{TYPE: "runtime Error", message: "max takes a non-empty array", EXISTS: true}
			}
			var max ArObject
			for i, v := range x {
				switch x := v.(type) {
				case ArObject:
					if i == 0 {
						max = x
					} else {
						compared, err := CompareObjects(max, x)

						if err.EXISTS {
							return nil, err
						}

						compared_int, Err := numberToInt64(compared)

						if Err != nil {
							return nil, ArErr{TYPE: "Type Error", message: Err.Error(), EXISTS: true}
						}

						if compared_int == -1 {
							max = x
						}
					}
				}
			}
			return max, ArErr{}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot get max of type '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["path"] = ArPath
	vars["typeof"] = builtinFunc{"typeof", func(a ...any) (any, ArErr) {
		if len(a) != 1 {
			return nil, ArErr{TYPE: "typeof", message: "typeof takes 1 argument, got " + fmt.Sprint(len(a)), EXISTS: true}
		}
		return ArString(typeof(a[0])), ArErr{}
	}}
	vars["sha256"] = builtinFunc{"sha256", func(a ...any) (any, ArErr) {
		if len(a) != 1 {
			return nil, ArErr{TYPE: "sha256", message: "sha256 takes 1 argument, got " + fmt.Sprint(len(a)), EXISTS: true}
		}
		a[0] = ArValidToAny(a[0])
		switch x := a[0].(type) {
		case string:
			return ArString(sha256Hash(x)), ArErr{}
		}
		return nil, ArErr{TYPE: "Type Error", message: "Cannot hash type '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	return Map(vars)
}
