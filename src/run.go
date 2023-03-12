package main

import (
	"fmt"
	"reflect"
)

// returns (number|string|nil), error
func runVal(line any, stack stack, stacklevel int) (any, ArErr) {
	if stacklevel > 10000 {
		return nil, ArErr{
			TYPE:    "RuntimeError",
			message: "stack overflow",
			line:    0,
			path:    "",
			code:    "",
			EXISTS:  true,
		}
	}
	switch x := line.(type) {
	case number:
		return x, ArErr{}
	case string:
		return x, ArErr{}
	case call:
		return runCall(x, stack, stacklevel+1)
	case factorial:
		return runFactorial(x, stack, stacklevel+1)
	case accessVariable:
		return readVariable(x, stack)
	case ArMapGet:
		return mapGet(x, stack, stacklevel+1)
	case ArClass:
		return x.MAP, ArErr{}
	case setVariable:
		return setVariableValue(x, stack, stacklevel+1)
	case negative:
		resp, err := runVal(x.VAL, stack, stacklevel+1)
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
		return runVal(x.VAL, stack, stacklevel+1)
	case operationType:
		return runOperation(x, stack, stacklevel+1)
	case dowrap:
		return runDoWrap(x, stack, stacklevel+1)
	case CallReturn:
		return runReturn(x, stack, stacklevel+1)
	case CallBreak:
		return runBreak(x, stack, stacklevel+1)
	case ArDelete:
		return runDelete(x, stack, stacklevel+1)
	case not:
		return runNot(x, stack, stacklevel+1)
	case ifstatement:
		return runIfStatement(x, stack, stacklevel+1)
	case whileLoop:
		return runWhileLoop(x, stack, stacklevel+1)
	case bool:
		return x, ArErr{}
	case nil:
		return nil, ArErr{}
	}
	fmt.Println("unreachable", reflect.TypeOf(line))
	panic("unreachable")
}

// returns error
func run(translated []any, stack stack) (any, ArErr, any) {
	var output any = nil
	for _, val := range translated {
		val, err := runVal(val, stack, 0)
		output = val
		if err.EXISTS {
			return nil, err, output
		}
	}
	return nil, ArErr{}, output
}
