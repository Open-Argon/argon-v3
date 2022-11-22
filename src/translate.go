package main

import (
	"fmt"
)

func translate(code string) {
	a := createNumber(7, 1)
	b, worked := stringToNumber("1e-1")
	if worked {
		output := newNumber().Mul(a, b)
		output.Add(output, a)
		output.Sub(output, a)
		fmt.Println(numberToString(output, true))
	}
}
