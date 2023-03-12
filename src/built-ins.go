package main

import (
	"fmt"
	"strings"
)

var vars = scope{}

func init() {
	vars["global"] = vars
	vars["term"] = ArTerm
	vars["input"] = builtinFunc{"input", ArgonInput}
	vars["number"] = builtinFunc{"number", ArgonNumber}
	vars["string"] = builtinFunc{"string", ArgonString}
	vars["infinity"] = infinity
	vars["length"] = builtinFunc{"length", func(a ...any) (any, ArErr) {
		switch x := a[0].(type) {
		case string:
			return len(x), ArErr{}
		case ArMap:
			return len(x), ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot get length of " + typeof(a[0]), EXISTS: true}
	}}
	vars["map"] = builtinFunc{"map", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return ArMap{}, ArErr{}
		}
		switch x := a[0].(type) {
		case ArMap:
			return x, ArErr{}
		case string:
			newmap := ArMap{}
			for i, v := range x {
				newmap[i] = string(v)
			}
			return newmap, ArErr{}
		case ArArray:
			newmap := ArMap{}
			for i, v := range x {
				switch y := v.(type) {
				case ArArray:
					if len(y) == 2 {
						newmap[y[0]] = y[1]
						continue
					}
				}
				newmap[i] = v
			}
			return newmap, ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot create map from '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["array"] = builtinFunc{"array", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return ArArray{}, ArErr{}
		}
		switch x := a[0].(type) {
		case ArArray:
			return x, ArErr{}
		case string:
			newarray := ArArray{}
			for _, v := range x {
				newarray = append(newarray, string(v))
			}
			return newarray, ArErr{}
		case ArMap:
			newarray := ArArray{}
			for key, val := range x {
				newarray = append(newarray, ArArray{key, val})
			}
			return newarray, ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot create array from '" + typeof(a[0]) + "'", EXISTS: true}
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
	vars["round"] = builtinFunc{"round", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return nil, ArErr{TYPE: "round", message: "round takes 1 argument",
				EXISTS: true}
		}
		precision := newNumber()
		if len(a) > 1 {
			switch x := a[1].(type) {
			case number:
				if !x.IsInt() {
					return nil, ArErr{TYPE: "TypeError", message: "Cannot round to '" + typeof(a[1]) + "'", EXISTS: true}
				}
				precision = x
			default:
				return nil, ArErr{TYPE: "TypeError", message: "Cannot round to '" + typeof(a[1]) + "'", EXISTS: true}
			}
		}

		switch x := a[0].(type) {
		case number:
			return round(newNumber().Set(x), int(precision.Num().Int64())), ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot round '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["floor"] = builtinFunc{"floor", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return nil, ArErr{TYPE: "floor", message: "floor takes 1 argument",
				EXISTS: true}
		}
		switch x := a[0].(type) {
		case number:
			return floor(x), ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot floor '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["ceil"] = builtinFunc{"ceil", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return nil, ArErr{TYPE: "ceil", message: "ceil takes 1 argument",
				EXISTS: true}
		}

		switch x := a[0].(type) {
		case number:
			return ceil(x), ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot ceil '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["append"] = builtinFunc{"append", func(a ...any) (any, ArErr) {
		if len(a) != 2 {
			return nil, ArErr{TYPE: "append", message: "append takes 2 arguments, got " + fmt.Sprint(len(a)),
				EXISTS: true}
		}
		switch x := a[0].(type) {
		case ArArray:
			return append(x, a[1]), ArErr{}
		case string:
			if typeof(a[1]) != "string" {
				return nil, ArErr{TYPE: "TypeError", message: "Cannot append '" + typeof(a[1]) + "' to string", EXISTS: true}
			}
			return strings.Join([]string{x, a[1].(string)}, ""), ArErr{}
		case ArMap:
			if typeof(a[1]) != "array" {
				return nil, ArErr{TYPE: "TypeError", message: "Cannot append '" + typeof(a[1]) + "' to map", EXISTS: true}
			}
			y := a[1].(ArArray)
			if len(y) != 2 {
				return nil, ArErr{TYPE: "TypeError", message: "Cannot append '" + typeof(a[1]) + "' to map", EXISTS: true}
			}
			x[y[0]] = y[1]
			return x, ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot append to '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["sqrt"] = builtinFunc{"sqrt", ArgonSqrt}
	vars["file"] = ArFile
	vars["random"] = ArRandom
	vars["json"] = ArJSON
}
