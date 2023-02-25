package main

import "fmt"

type builtinFunc struct {
	name string
	FUNC func(...any) (any, ArErr)
}

func ArgonLog(args ...any) (any, ArErr) {
	output := []any{}
	for i := 0; i < len(args); i++ {
		output = append(output, anyToArgon(args[i], false))
	}
	fmt.Println(output...)
	return nil, ArErr{}
}

func ArgonAdd(args ...any) (any, ArErr) {
	return reduce(func(x any, y any) any {
		return newNumber().Add(x.(number), y.(number))
	}, args), ArErr{}
}
func ArgonDiv(args ...any) (any, ArErr) {
	return reduce(func(x any, y any) any {
		return newNumber().Quo(y.(number), x.(number))
	}, args), ArErr{}
}

func ArgonMult(args ...any) (any, ArErr) {
	return reduce(func(x any, y any) any {
		return newNumber().Mul(y.(number), x.(number))
	}, args), ArErr{}
}
