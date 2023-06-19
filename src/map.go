package main

import (
	"fmt"
	"strings"
	"sync"
)

var mapCompiled = makeRegex(`( )*<( |\n)*(((.|\n)+)(,(.|\n)+)*)?( |\n)*>( )*`)

type createMap struct {
	body anymap
	code string
	line int
	path string
}

func isMap(code UNPARSEcode) bool {
	return mapCompiled.MatchString(code.code)
}

func parseMap(code UNPARSEcode, index int, codelines []UNPARSEcode) (any, bool, ArErr, int) {
	trimmed := strings.Trim(code.code, " ")
	trimmed = trimmed[1 : len(trimmed)-1]
	debugPrintln(trimmed)
	return Map(anymap{}), true, ArErr{}, 1
}

func Map(m anymap) ArObject {
	var mutex = sync.RWMutex{}
	obj := ArObject{
		obj: anymap{
			"__value__": m,
			"__name__":  "map",
			"get": builtinFunc{
				"get",
				func(args ...any) (any, ArErr) {
					if len(args) < 1 || len(args) > 2 {
						return nil, ArErr{
							TYPE:    "Runtime Error",
							message: "expected 1 or 2 argument, got " + fmt.Sprint(len(args)),
							EXISTS:  true,
						}
					}
					var DEFAULT any
					key := ArValidToAny(args[0])
					if isUnhashable(key) {
						return nil, ArErr{
							TYPE:    "Runtime Error",
							message: "unhashable type: " + typeof(key),
							EXISTS:  true,
						}
					}
					if len(args) == 2 {
						DEFAULT = (args[1])
					}
					mutex.RLock()
					if _, ok := m[key]; !ok {
						mutex.RUnlock()
						return DEFAULT, ArErr{}
					}
					v := m[key]
					mutex.RUnlock()
					return v, ArErr{}
				},
			},
			"__Contains__": builtinFunc{
				"__Contains__",
				func(args ...any) (any, ArErr) {
					if len(args) != 1 {
						return nil, ArErr{
							TYPE:    "TypeError",
							message: "expected 1 argument, got " + fmt.Sprint(len(args)),
							EXISTS:  true,
						}
					}
					key := ArValidToAny(args[0])
					if isUnhashable(key) {
						return false, ArErr{}
					}
					mutex.RLock()
					if _, ok := m[key]; !ok {
						mutex.RUnlock()
						return false, ArErr{}
					}
					mutex.RUnlock()
					return true, ArErr{}
				},
			},
			"__NotContains__": builtinFunc{
				"__NotContains__",
				func(args ...any) (any, ArErr) {
					if len(args) != 1 {
						return nil, ArErr{
							TYPE:    "TypeError",
							message: "expected 1 argument, got " + fmt.Sprint(len(args)),
							EXISTS:  true,
						}
					}
					key := ArValidToAny(args[0])
					if isUnhashable(key) {
						return true, ArErr{}
					}
					mutex.RLock()
					if _, ok := m[key]; !ok {
						mutex.RUnlock()
						return true, ArErr{}
					}
					mutex.RUnlock()
					return false, ArErr{}
				},
			},
			"__setindex__": builtinFunc{
				"__setindex__",
				func(args ...any) (any, ArErr) {
					if len(args) != 2 {
						return nil, ArErr{
							TYPE:    "TypeError",
							message: "expected 2 arguments, got " + fmt.Sprint(len(args)),
							EXISTS:  true,
						}
					}
					if isUnhashable(args[0]) {
						return nil, ArErr{
							TYPE:    "Runtime Error",
							message: "unhashable type: " + typeof(args[0]),
							EXISTS:  true,
						}
					}
					key := ArValidToAny(args[0])
					mutex.Lock()
					m[key] = args[1]
					mutex.Unlock()
					return nil, ArErr{}
				},
			},
			"__getindex__": builtinFunc{
				"__getindex__",
				func(args ...any) (any, ArErr) {
					if len(args) != 1 {
						return nil, ArErr{
							TYPE:    "TypeError",
							message: "expected 1 argument, got " + fmt.Sprint(len(args)),
							EXISTS:  true,
						}
					}
					key := ArValidToAny(args[0])
					if isUnhashable(key) {
						return nil, ArErr{
							TYPE:    "Runtime Error",
							message: "unhashable type: " + typeof(key),
							EXISTS:  true,
						}
					}
					mutex.RLock()
					if _, ok := m[key]; !ok {
						mutex.RUnlock()
						return nil, ArErr{
							TYPE:    "KeyError",
							message: "key " + fmt.Sprint(key) + " not found",
							EXISTS:  true,
						}
					}
					v := m[key]
					mutex.RUnlock()
					return v, ArErr{}
				},
			},
		},
	}
	obj.obj["__Equal__"] = builtinFunc{
		"__Equal__",
		func(args ...any) (any, ArErr) {
			debugPrintln("Equal", args)
			if len(args) != 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected 1 argument, got " + fmt.Sprint(len(args)),
					EXISTS:  true,
				}
			}
			if typeof(args[0]) != "map" {
				return false, ArErr{}
			}
			a := ArValidToAny(args[0]).(anymap)
			mutex.RLock()
			if len(m) != len(a) {
				mutex.RUnlock()
				return false, ArErr{}
			}
			for k, v := range m {
				debugPrintln(k, v)
				if _, ok := a[k]; !ok {
					mutex.RUnlock()
					return false, ArErr{}
				}
				val, err := runOperation(operationType{
					operation: 9,
					value1:    v,
					value2:    a[k],
				}, stack{}, 0)
				if err.EXISTS {
					return val, err
				}
				if !anyToBool(val) {
					mutex.RUnlock()
					return false, ArErr{}
				}
			}
			mutex.RUnlock()
			return true, ArErr{}
		},
	}
	obj.obj["__dir__"] = builtinFunc{
		"__dir__",
		func(args ...any) (any, ArErr) {
			x := []any{}
			for k := range m {
				x = append(x, k)
			}
			return x, ArErr{}
		}}
	return obj
}
