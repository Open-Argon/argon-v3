package main

import (
	"fmt"
	"strings"
	"syscall/js"
)

func argonToJsValid(argon any) any {
	switch x := argon.(type) {
	case number:
		f, _ := x.Float64()
		return f
	case ArObject:
		if typeof(x) == "array" {
			arr := js.Global().Get("Array").New()
			for i, v := range x.obj["__value__"].([]any) {
				arr.SetIndex(i, argonToJsValid(v))
			}
			return arr
		} else if typeof(x) == "string" {
			return x.obj["__value__"].(string)
		}

		obj := js.Global().Get("Object").New()
		for k, v := range x.obj {
			obj.Set(anyToArgon(k, false, false, 3, 0, false, 0), argonToJsValid(v))
		}
		return obj
	case bool, string:
		return x
	default:
		return nil
	}
}

func wasmRun(code string, allowDocument bool) (any, ArErr) {
	JSclearTimers()
	initRandom()
	global := makeGlobal(allowDocument)
	lines := strings.Split(code, "\n")
	codelines := []UNPARSEcode{}
	for i := 0; i < len(lines); i++ {
		codelines = append(codelines, UNPARSEcode{
			lines[i],
			lines[i],
			i + 1,
			"<wasm>",
		})
	}

	translated, translationerr := translate(codelines)
	if translationerr.EXISTS {
		return nil, translationerr
	}
	local := newscope()
	localvars := Map(anymap{
		"program": Map(anymap{
			"args":   []any{},
			"origin": "",
			"import": builtinFunc{"import", func(args ...any) (any, ArErr) {
				return nil, ArErr{"Import Error", "Cannot Import in WASM", 0, "<wasm>", "", true}
			}},
			"cwd": "",
			"exc": "",
			"file": Map(anymap{
				"name": "<wasm>",
				"path": "",
			}),
			"main":  true,
			"scope": global,
		}),
	})
	return ThrowOnNonLoop(run(translated, stack{global, localvars, local}))
}

func await(awaitable js.Value) ([]js.Value, []js.Value) {
	then := make(chan []js.Value)
	defer close(then)
	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		then <- args
		return nil
	})
	defer thenFunc.Release()

	catch := make(chan []js.Value)
	defer close(catch)
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		catch <- args
		return nil
	})
	defer catchFunc.Release()

	awaitable.Call("then", thenFunc).Call("catch", catchFunc)

	select {
	case result := <-then:
		return result, nil
	case err := <-catch:
		return nil, err
	}
}

var IntervalList = []int{}
var TimeoutList = []int{}

func JSclearTimers() {
	for _, v := range IntervalList {
		js.Global().Call("clearInterval", v)
	}
	for _, v := range TimeoutList {
		js.Global().Call("clearTimeout", v)
	}
}

var ArJS = Map(anymap{
	"setTimeout": builtinFunc{"setTimeout", func(args ...any) (any, ArErr) {
		if len(args) > 2 || len(args) < 1 {
			return nil, ArErr{"TypeError", "Expected 1 or 2 argument, got " + fmt.Sprint(len(args)), 0, "<wasm>", "", true}
		}
		if typeof(args[0]) != "function" {
			return nil, ArErr{"TypeError", "Expected function, got " + typeof(args[0]), 0, "<wasm>", "", true}
		}
		var ms int64 = 0
		if len(args) == 2 {
			if typeof(args[1]) != "number" {
				return nil, ArErr{"TypeError", "Expected number, got " + typeof(args[1]), 0, "<wasm>", "", true}
			}
			if !args[1].(number).IsInt() {
				return nil, ArErr{"TypeError", "Expected integer, got float", 0, "<wasm>", "", true}
			}
			ms = args[1].(number).Num().Int64()
		}
		f := js.FuncOf(func(this js.Value, a []js.Value) interface{} {
			runCall(
				call{
					callable: args[0],
					args:     []any{},
				},
				stack{},
				0,
			)
			return nil
		})
		n := js.Global().Call("setTimeout", f, ms).Int()
		TimeoutList = append(TimeoutList, n)
		return newNumber().SetInt64(int64(n)), ArErr{}
	}},
	"setInterval": builtinFunc{"setInterval", func(args ...any) (any, ArErr) {
		if len(args) > 2 || len(args) < 1 {
			return nil, ArErr{"TypeError", "Expected 1 or 2 argument, got " + fmt.Sprint(len(args)), 0, "<wasm>", "", true}
		}
		if typeof(args[0]) != "function" {
			return nil, ArErr{"TypeError", "Expected function, got " + typeof(args[0]), 0, "<wasm>", "", true}
		}
		var ms int64 = 0
		if len(args) == 2 {
			if typeof(args[1]) != "number" {
				return nil, ArErr{"TypeError", "Expected number, got " + typeof(args[1]), 0, "<wasm>", "", true}
			}
			if !args[1].(number).IsInt() {
				return nil, ArErr{"TypeError", "Expected integer, got float", 0, "<wasm>", "", true}
			}
			ms = args[1].(number).Num().Int64()
		}
		f := js.FuncOf(func(this js.Value, a []js.Value) interface{} {
			runCall(
				call{
					callable: args[0],
					args:     []any{},
				},
				stack{},
				0,
			)
			return nil
		})
		n := js.Global().Call("setInterval", f, ms).Int()
		IntervalList = append(IntervalList, n)
		return newNumber().SetInt64(int64(n)), ArErr{}
	}},
	"clearTimeout": builtinFunc{"clearTimeout", func(args ...any) (any, ArErr) {
		if len(args) != 1 {
			return nil, ArErr{"TypeError", "Expected 1 argument, got " + fmt.Sprint(len(args)), 0, "<wasm>", "", true}
		}
		if typeof(args[0]) != "number" {
			return nil, ArErr{"TypeError", "Expected number, got " + typeof(args[0]), 0, "<wasm>", "", true}
		}
		if !args[0].(number).IsInt() {
			return nil, ArErr{"TypeError", "Expected integer, got float", 0, "<wasm>", "", true}
		}
		n := args[0].(number).Num().Int64()
		for i, v := range TimeoutList {
			if v == int(n) {
				TimeoutList = append(TimeoutList[:i], TimeoutList[i+1:]...)
				break
			}
		}
		js.Global().Call("clearTimeout", n)
		return nil, ArErr{}
	}},
	"clearInterval": builtinFunc{"clearInterval", func(args ...any) (any, ArErr) {
		if len(args) != 1 {
			return nil, ArErr{"TypeError", "Expected 1 argument, got " + fmt.Sprint(len(args)), 0, "<wasm>", "", true}
		}
		if typeof(args[0]) != "number" {
			return nil, ArErr{"TypeError", "Expected number, got " + typeof(args[0]), 0, "<wasm>", "", true}
		}
		if !args[0].(number).IsInt() {
			return nil, ArErr{"TypeError", "Expected integer, got float", 0, "<wasm>", "", true}
		}
		n := args[0].(number).Num().Int64()
		for i, v := range IntervalList {
			if v == int(n) {
				IntervalList = append(IntervalList[:i], IntervalList[i+1:]...)
				break
			}
		}
		js.Global().Call("clearInterval", n)
		return nil, ArErr{}
	}},
})
