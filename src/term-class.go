package main

import "fmt"

func ArgonLog(args ...any) (any, ArErr) {
	output := []any{}
	for i := 0; i < len(args); i++ {
		output = append(output, anyToArgon(args[i], false, true, 3, 0))
	}
	fmt.Println(output...)
	return nil, ArErr{}
}

var ArTerm = ArMap{
	"log": builtinFunc{"log", ArgonLog},
}
