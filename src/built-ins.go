package main

var vars = Map(anymap{})

func init() {
	vars.obj["global"] = vars
	vars.obj["term"] = ArTerm
	vars.obj["input"] = builtinFunc{"input", ArgonInput}
	vars.obj["passwordInput"] = builtinFunc{"passwordInput", ArgonPassworInput}
	vars.obj["number"] = builtinFunc{"number", ArgonNumber}
	vars.obj["string"] = builtinFunc{"string", ArgonString}
	vars.obj["infinity"] = infinity
	vars.obj["map"] = builtinFunc{"map", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return Map(anymap{}), ArErr{}
		}
		switch x := a[0].(type) {
		case ArObject:
			if x.TYPE == "array" {
				newmap := anymap{}
				for i, v := range x.obj["__value__"].([]any) {
					switch y := v.(type) {
					case []any:
						if len(y) == 2 {
							if isUnhashable(y[0]) {
								return nil, ArErr{TYPE: "TypeError", message: "Cannot use unhashable value as key: " + typeof(y[0]), EXISTS: true}
							}
							newmap[y[0]] = y[1]
							continue
						}
					}
					newmap[i] = v
				}
				return Map(newmap), ArErr{}
			}
			return x, ArErr{}
		case string:
			newmap := anymap{}
			for i, v := range x {
				newmap[i] = ArString(string(v))
			}
			return Map(newmap), ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot create map from '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars.obj["array"] = builtinFunc{"array", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return ArArray([]any{}), ArErr{}
		}
		switch x := a[0].(type) {
		case string:
			newarray := []any{}
			for _, v := range x {
				newarray = append(newarray, ArString(string(v)))
			}
			return ArArray(newarray), ArErr{}
		case ArObject:
			if x.TYPE == "array" {
				return x, ArErr{}
			}
			newarray := []any{}
			for key, val := range x.obj {
				newarray = append(newarray, []any{key, val})
			}
			return ArArray(newarray), ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot create array from '" + typeof(a[0]) + "'", EXISTS: true}
	}}
	vars.obj["boolean"] = builtinFunc{"boolean", func(a ...any) (any, ArErr) {
		if len(a) == 0 {
			return false, ArErr{}
		}
		return anyToBool(a[0]), ArErr{}
	}}
	vars.obj["time"] = ArTime
	vars.obj["PI"] = PI
	vars.obj["π"] = PI
	vars.obj["e"] = e
	vars.obj["ln"] = builtinFunc{"ln", ArgonLn}
	vars.obj["log"] = builtinFunc{"log", ArgonLog}
	vars.obj["logN"] = builtinFunc{"logN", ArgonLogN}
	vars.obj["thread"] = builtinFunc{"thread", ArThread}
	vars.obj["round"] = builtinFunc{"round", func(a ...any) (any, ArErr) {
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
	vars.obj["floor"] = builtinFunc{"floor", func(a ...any) (any, ArErr) {
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
	vars.obj["ceil"] = builtinFunc{"ceil", func(a ...any) (any, ArErr) {
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
	vars.obj["sqrt"] = builtinFunc{"sqrt", ArgonSqrt}
	vars.obj["open"] = builtinFunc{"open", ArOpen}
	vars.obj["random"] = ArRandom
	vars.obj["json"] = ArJSON
	vars.obj["sin"] = ArSin
	vars.obj["arcsin"] = ArArcsin
}
