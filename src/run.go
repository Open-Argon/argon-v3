package main

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
	}
	panic("unreachable")
}

// returns error
func run(translated []any, stack stack) (any, ArErr) {
	for _, val := range translated {
		_, err := runVal(val, stack)
		if err.EXISTS {
			return nil, err
		}
	}
	return nil, ArErr{}
}
