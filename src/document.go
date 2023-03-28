package main

import (
	"syscall/js"
)

func windowElement(element js.Value) ArObject {
	return ArObject{
		TYPE: "map",
		obj: anymap{
			"innerHTML": builtinFunc{"innerHTML", func(args ...any) (any, ArErr) {
				if len(args) > 0 {
					if typeof(args[0]) != "string" {
						return nil, ArErr{"Argument Error", "innerHTML only accepts strings", 0, "", "", true}
					}
					element.Set("innerHTML", args[0].(string))
				}
				return element.Get("innerHTML").String(), ArErr{}
			}},
			"innerText": builtinFunc{"innerText", func(args ...any) (any, ArErr) {
				if len(args) > 0 {
					if typeof(args[0]) != "string" {
						return nil, ArErr{"Argument Error", "innerText only accepts strings", 0, "", "", true}
					}
					element.Set("innerText", args[0].(string))
				}
				return element.Get("innerText").String(), ArErr{}
			}},
			"addEventListener": builtinFunc{"addEventListener", func(args ...any) (any, ArErr) {
				if len(args) < 2 {
					return nil, ArErr{"Argument Error", "Not enough arguments for addEventListener", 0, "", "", true}
				}
				if typeof(args[0]) != "string" {
					return nil, ArErr{"Argument Error", "addEventListener's first argument must be a string", 0, "", "", true}
				}
				event := args[0].(string)
				if typeof(args[1]) != "function" {
					return nil, ArErr{"Argument Error", "addEventListener's second argument must be a function", 0, "", "", true}
				}
				callable := args[1]
				element.Call("addEventListener", event, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					runCall(call{
						callable: callable,
						args:     []any{},
					}, stack{}, 0)
					return nil
				}))
				return nil, ArErr{}
			}},
			"removeEventListener": builtinFunc{"removeEventListener", func(args ...any) (any, ArErr) {
				if len(args) < 2 {
					return nil, ArErr{"Argument Error", "Not enough arguments for removeEventListener", 0, "", "", true}
				}
				if typeof(args[0]) != "string" {
					return nil, ArErr{"Argument Error", "removeEventListener's first argument must be a string", 0, "", "", true}
				}
				event := args[0].(string)
				if typeof(args[1]) != "function" {
					return nil, ArErr{"Argument Error", "removeEventListener's second argument must be a function", 0, "", "", true}
				}
				callable := args[1]
				element.Call("removeEventListener", event, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					runCall(call{
						callable: callable,
						args:     []any{},
					}, stack{}, 0)
					return nil
				}))
				return nil, ArErr{}
			}},
			"appendChild": builtinFunc{"appendChild", func(args ...any) (any, ArErr) {
				if len(args) < 1 {
					return nil, ArErr{"Argument Error", "Not enough arguments for appendChild", 0, "", "", true}
				}
				if typeof(args[0]) != "map" {
					return nil, ArErr{"Argument Error", "appendChild's first argument must be a map", 0, "", "", true}
				}
				child := args[0].(anymap)
				if child["__TYPE__"] != "windowElement" {
					return nil, ArErr{"Argument Error", "appendChild's first argument must be an element", 0, "", "", true}
				}
				element.Call("appendChild", child["__element__"])
				return nil, ArErr{}
			}},
			"removeChild": builtinFunc{"removeChild", func(args ...any) (any, ArErr) {
				if len(args) < 1 {
					return nil, ArErr{"Argument Error", "Not enough arguments for removeChild", 0, "", "", true}
				}
				if typeof(args[0]) != "map" {
					return nil, ArErr{"Argument Error", "removeChild's first argument must be a map", 0, "", "", true}
				}
				child := args[0].(anymap)
				if child["__TYPE__"] != "windowElement" {
					return nil, ArErr{"Argument Error", "removeChild's first argument must be an element", 0, "", "", true}
				}
				element.Call("removeChild", child["__element__"])
				return nil, ArErr{}
			}},
			"setAttribute": builtinFunc{"setAttribute", func(args ...any) (any, ArErr) {
				if len(args) < 2 {
					return nil, ArErr{"Argument Error", "Not enough arguments for setAttribute", 0, "", "", true}
				}
				if typeof(args[0]) != "string" {
					return nil, ArErr{"Argument Error", "setAttribute's first argument must be a string", 0, "", "", true}
				}
				element.Call("setAttribute", args[0].(string), anyToArgon(args[1], false, false, 3, 0, false, 0))
				return nil, ArErr{}
			}},
			"__element__": element,
			"__TYPE__":    "windowElement",
		},
	}
}

var ArDocument = Map(
	anymap{
		"body": builtinFunc{"getElementById", func(args ...any) (any, ArErr) {
			return windowElement(js.Global().Get("document").Get("body")), ArErr{}
		}},
		"head": builtinFunc{"getElementById", func(args ...any) (any, ArErr) {
			return windowElement(js.Global().Get("document").Get("head")), ArErr{}
		}},
		"getElementById": builtinFunc{"getElementById", func(args ...any) (any, ArErr) {
			if len(args) < 1 {
				return nil, ArErr{"Argument Error", "Not enough arguments for getElementById", 0, "", "", true}
			}
			if typeof(args[0]) != "string" {
				return nil, ArErr{"Argument Error", "getElementById's first argument must be a string", 0, "", "", true}
			}
			id := args[0].(string)
			result := js.Global().Get("document").Call("getElementById", id)
			if js.Null().Equal(result) {
				return nil, ArErr{}
			}
			return windowElement(result), ArErr{}
		}},
		"createElement": builtinFunc{"createElement", func(args ...any) (any, ArErr) {
			if len(args) < 1 {
				return nil, ArErr{"Argument Error", "Not enough arguments for createElement", 0, "", "", true}
			}
			if typeof(args[0]) != "string" {
				return nil, ArErr{"Argument Error", "createElement's first argument must be a string", 0, "", "", true}
			}
			tag := args[0].(string)
			return windowElement(js.Global().Get("document").Call("createElement", tag)), ArErr{}
		}},
		"createTextNode": builtinFunc{"createTextNode", func(args ...any) (any, ArErr) {
			if len(args) < 1 {
				return nil, ArErr{"Argument Error", "Not enough arguments for createTextNode", 0, "", "", true}
			}
			if typeof(args[0]) != "string" {
				return nil, ArErr{"Argument Error", "createTextNode's first argument must be a string", 0, "", "", true}
			}
			text := args[0].(string)
			return windowElement(js.Global().Get("document").Call("createTextNode", text)), ArErr{}
		}},
		"createComment": builtinFunc{"createComment", func(args ...any) (any, ArErr) {
			if len(args) < 1 {
				return nil, ArErr{"Argument Error", "Not enough arguments for createComment", 0, "", "", true}
			}
			if typeof(args[0]) != "string" {
				return nil, ArErr{"Argument Error", "createComment's first argument must be a string", 0, "", "", true}
			}
			text := args[0].(string)
			return windowElement(js.Global().Get("document").Call("createComment", text)), ArErr{}
		}},
		"createDocumentFragment": builtinFunc{"createDocumentFragment", func(args ...any) (any, ArErr) {
			return windowElement(js.Global().Get("document").Call("createDocumentFragment")), ArErr{}
		}},
		"addEventListener": builtinFunc{"addEventListener", func(args ...any) (any, ArErr) {
			if len(args) < 2 {
				return nil, ArErr{"Argument Error", "Not enough arguments for addEventListener", 0, "", "", true}
			}
			if typeof(args[0]) != "string" {
				return nil, ArErr{"Argument Error", "addEventListener's first argument must be a string", 0, "", "", true}
			}
			event := args[0].(string)
			if typeof(args[1]) != "function" {
				return nil, ArErr{"Argument Error", "addEventListener's second argument must be a function", 0, "", "", true}
			}
			callable := args[1]
			js.Global().Get("document").Call("addEventListener", event, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				runCall(call{
					callable: callable,
					args:     []any{},
				}, stack{}, 0)
				return nil
			}))
			return nil, ArErr{}
		}},
		"removeEventListener": builtinFunc{"removeEventListener", func(args ...any) (any, ArErr) {
			if len(args) < 2 {
				return nil, ArErr{"Argument Error", "Not enough arguments for removeEventListener", 0, "", "", true}
			}
			if typeof(args[0]) != "string" {
				return nil, ArErr{"Argument Error", "removeEventListener's first argument must be a string", 0, "", "", true}
			}
			event := args[0].(string)
			if typeof(args[1]) != "function" {
				return nil, ArErr{"Argument Error", "removeEventListener's second argument must be a function", 0, "", "", true}
			}
			callable := args[1]
			js.Global().Get("document").Call("removeEventListener", event, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				runCall(call{
					callable: callable,
					args:     []any{},
				}, stack{}, 0)
				return nil
			}))
			return nil, ArErr{}
		}},
		"__TYPE__": "document",
	},
)
