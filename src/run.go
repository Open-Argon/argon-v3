package main

// returns (number|string|nil), error
func runVal(line any, stack []map[string]variableValue) (any, ArErr) {
	switch x := line.(type) {
	case translateNumber:
		return (x.number), ArErr{}
	case translateString:
		return (x.str), ArErr{}
	case call:
		return runCall(x, stack)
	case accessVariable:
		return readVariable(x, stack)
	}
	panic("unreachable")
}

// returns error
func run(translated []any, stack []map[string]variableValue) (any, ArErr) {
	for _, val := range translated {
		_, err := runVal(val, stack)
		if err.EXISTS {
			return nil, err
		}
	}
	return nil, ArErr{}
}
