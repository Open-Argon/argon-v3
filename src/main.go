package main

import (
	"os"
	"syscall/js"
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
	c := make(chan ArObject)
	obj := js.Global().Get("Object").New()
	obj.Set("eval", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		code := ""
		if len(args) > 0 {
			code = args[0].String()
		}
		val, err := wasmRun(code)
		if err.EXISTS {
			panicErr(err)
			return js.Null()
		}

		return js.ValueOf(argonToJsValid(val))
	}))
	js.Global().Set("Ar", obj)
	<-c
}
