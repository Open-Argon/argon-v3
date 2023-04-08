package main

import (
	"fmt"
	"strings"
	"sync"
)

var mapCompiled = makeRegex(`( *)\{(((( *).+( *):( *).+( *))|(` + spacelessVariable + `))(( *)\,(( *).+( *):( *).+( *))|(` + spacelessVariable + `)))*\}( *)`)

type createMap struct {
	body anymap
	code string
	line int
	path string
}

func isMap(code UNPARSEcode) bool {
	return mapCompiled.MatchString(code.code)
}

func parseMap(code UNPARSEcode) (any, UNPARSEcode) {
	trimmed := strings.Trim(code.code, " ")
	trimmed = trimmed[1 : len(trimmed)-1]
	fmt.Println(trimmed)
	return nil, UNPARSEcode{}
}

func Map(m anymap) ArObject {
	var mutex = sync.RWMutex{}
	return ArObject{
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
}
