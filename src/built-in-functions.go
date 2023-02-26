package main

import (
	"fmt"
	"time"
)

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

func ArgonSleep(args ...any) (any, ArErr) {
	if len(args) > 0 {
		float, _ := args[0].(number).Float64()
		time.Sleep(time.Duration(float*1000000000) * time.Nanosecond)
	}
	return nil, ArErr{}
}

func ArgonTimestamp(args ...any) (any, ArErr) {
	return newNumber().Quo(newNumber().SetInt64(time.Now().UnixNano()), newNumber().SetInt64(1000000000)), ArErr{}
}
