package main

import (
	"os"
)

// args without the program path
var Args = os.Args[1:]

type stack = []ArObject

func newscope() ArObject {
	return ArObject{
		TYPE: "map",
		obj:  make(anymap),
	}
}

func main() {
	ex, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	if len(Args) == 0 {
		shell()
		os.Exit(0)
	}
	_, err := importMod(Args[0], ex, true)
	if err.EXISTS {
		panicErr(err)
		os.Exit(1)
	}
}
