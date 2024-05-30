package main

import "fmt"

func AReval(a ...any) (any, ArErr) {
	if len(a) < 1 || len(a) > 2 {
		return nil, ArErr{TYPE: "Type Error", message: "expected 1 or 2 argument, got " + fmt.Sprint(len(a)), EXISTS: true}
	}
	var expression string
	if typeof(a[0]) != "string" {
		return nil, ArErr{TYPE: "Type Error", message: "expected string as first argument, got " + typeof(a[0]), EXISTS: true}
	}
	expression = ArValidToAny(a[0]).(string)

	// translate the expression
	var code UNPARSEcode = UNPARSEcode{
		code:     expression,
		realcode: expression,
		line:     0,
		path:     "eval",
	}
	translated, err := translate([]UNPARSEcode{code})
	if err.EXISTS {
		return nil, err
	}

	var scope ArObject
	if len(a) == 2 {
		if typeof(a[1]) != "map" {
			return nil, ArErr{TYPE: "Type Error", message: "expected map as second argument, got " + typeof(a[1]), EXISTS: true}
		}
		scope = a[1].(ArObject)
	} else {
		scope = newscope()
	}

	var stack stack = []ArObject{scope}

	// run the translated expression
	return run(translated, stack)
}
