package main

import (
	"os"
)

// args without the program path
var Args = os.Args[1:]

func main() {
	ex, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	if len(Args) == 0 {
		panic("No file specified")
	}
	err := importMod(Args[0], ex, true)
	if err.EXISTS {
		panicErr(err)
	}
}
