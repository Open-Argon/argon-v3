package main

import (
	"fmt"
)

func translate(code string) {
	output, _ := newNumber().SetString("3.1415")
	fmt.Println(numberToString(output, 0))
}
