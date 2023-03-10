package main

import "github.com/wadey/go-rounding"

var vars = scope{}

func init() {
	vars["global"] = vars
	vars["term"] = ArTerm
	vars["true"] = true
	vars["false"] = false
	vars["null"] = nil
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
		case []any:
			newmap := ArMap{}
			for i, v := range x {
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
			for _, v := range x {
				newarray = append(newarray, v)
			}
			return newarray, ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot create array from '" + typeof(a[0]) + "'", EXISTS: true}
	}}
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
			return rounding.Round(newNumber().Set(x), int(precision.Num().Int64()), rounding.HalfUp), ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot round '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars["time"] = ArTime
	vars["PI"] = PI
	vars["π"] = PI
	vars["e"] = e
	sqrt := builtinFunc{"sqrt", ArgonSqrt}
	vars["sqrt"] = sqrt
	vars["thread"] = builtinFunc{"thread", ArThread}
}
