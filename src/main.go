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
	initRandom()
	garbageCollect()
	global := makeGlobal()
	if len(Args) == 0 {
		shell(global)
		os.Exit(0)
	}
	ex, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	_, err := importMod(Args[0], ex, true, global)
	if err.EXISTS {
		panicErr(err)
		os.Exit(1)
	}
}
