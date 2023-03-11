package main

import (
	"fmt"
	"reflect"
)

// returns (number|string|nil), error
func runVal(line any, stack stack) (any, ArErr) {
	if len(stack) > 500 {
		return nil, ArErr{
			TYPE:    "Stack overflow",
			message: "the stack has exceeded 500 levels",
			EXISTS:  true,
		}
	}
	switch x := line.(type) {
	case number:
		return x, ArErr{}
	case string:
		return x, ArErr{}
	case call:
		return runCall(x, stack)
	case factorial:
		return runFactorial(x, stack)
	case accessVariable:
		return readVariable(x, stack)
	case ArMapGet:
		return mapGet(x, stack)
	case ArClass:
		return x.MAP, ArErr{}
	case setVariable:
		return setVariableValue(x, stack)
	case negative:
		resp, err := runVal(x.VAL, stack)
		resp = classVal(resp)
		if err.EXISTS {
			return nil, err
		}
		switch y := resp.(type) {
		case number:
			return newNumber().Neg(y), ArErr{}
		}
		return nil, ArErr{
			TYPE:    "TypeError",
			message: "cannot negate a non-number",
			EXISTS:  true,
		}
	case brackets:
		return runVal(x.VAL, stack)
	case operationType:
		return runOperation(x, stack)
	case dowrap:
		return runDoWrap(x, stack)
	case CallJumpStatment:
		return runJumpStatment(x, stack)
	case ArDelete:
		return runDelete(x, stack)
	}
	fmt.Println("unreachable", reflect.TypeOf(line))
	panic("unreachable")
}

// returns error
func run(translated []any, stack stack) (any, ArErr, any) {
	var output any = nil
	for _, val := range translated {
		val, err := runVal(val, stack)
		output = val
		if err.EXISTS {
			return nil, err, output
		}
	}
	return nil, ArErr{}, output
}
