package main

var vars = scope{}

func init() {
	vars["window"] = vars
	vars["term"] = ArTerm
	vars["true"] = true
	vars["false"] = false
	vars["null"] = nil
	vars["input"] = builtinFunc{"input", ArgonInput}
	vars["number"] = builtinFunc{"number", ArgonNumber}
	vars["mult"] = builtinFunc{"mult", ArgonMult}
	vars["length"] = builtinFunc{"length", func(a ...any) (any, ArErr) {
		switch x := a[0].(type) {
		case string:
			return len(x), ArErr{}
		case ArMap:
			return len(x), ArErr{}
		}
		return nil, ArErr{TYPE: "TypeError", message: "Cannot get length of " + typeof(a[0]), EXISTS: true}
	}}
	vars["time"] = ArTime
	vars["PI"] = PI
	vars["π"] = PI
	vars["e"] = e
	sqrt := builtinFunc{"sqrt", ArgonSqrt}
	vars["sqrt"] = sqrt
}
