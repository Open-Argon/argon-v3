package main

import (
	"os"
)

// args without the program path
var Args = os.Args[1:]

type scope = ArMap
type stack = []scope

func main() {
	ex, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	if len(Args) == 0 {
		shell()
		os.Exit(0)
	}
	err := importMod(Args[0], ex, true)
	if err.EXISTS {
		panicErr(err)
		os.Exit(1)
	}
}
