package main

import (
	"fmt"
	"strings"
	"sync"
)

var mapCompiled = makeRegex(`( )*{( |\n)*(((.|\n)+)(,(.|\n)+)*)?( |\n)*}( )*`)

type createMap struct {
	body [][2]any
	code string
	line int
	path string
}

func isMap(code UNPARSEcode) bool {
	return mapCompiled.MatchString(code.code)
}

func runCreateMap(m createMap, stack stack, stacklevel int) (any, ArErr) {
	var body = m.body
	var newmap = anymap{}
	for _, pair := range body {
		key := pair[0]
		val := pair[1]
		keyVal, err := runVal(key, stack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		keyVal = ArValidToAny(keyVal)
		valVal, err := runVal(val, stack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		if isUnhashable(keyVal) {
			return nil, ArErr{
				"TypeError",
				"unhashable type: '" + typeof(keyVal) + "'",
				m.line,
				m.path,
				m.code,
				true,
			}
		}
		newmap[keyVal] = valVal
	}
	return newmap, ArErr{}
}

func parseMap(code UNPARSEcode, index int, codelines []UNPARSEcode) (any, bool, ArErr, int) {
	trimmed := strings.TrimSpace(code.code)
	trimmed = trimmed[1 : len(trimmed)-1]
	if len(trimmed) == 0 {
		return createMap{
			body: [][2]any{},
			code: code.realcode,
			line: code.line,
			path: code.path,
		}, true, ArErr{}, 1
	}
	var body [][2]any
	var LookingAtKey bool = true
	var current int
	var currentKey any
	var countIndex int = 1
	for i := 0; i < len(trimmed); i++ {
		var str string
		if LookingAtKey {
			if trimmed[i] != ':' {
				continue
			}
			str = trimmed[current:i]
		} else {
			if trimmed[i] != ',' && i != len(trimmed)-1 {
				continue
			}
			if i == len(trimmed)-1 {
				str = trimmed[current:]
			} else {
				str = trimmed[current:i]
			}
		}
		var value any
		if LookingAtKey && variableCompile.MatchString(str) {
			value = strings.TrimSpace(str)
		} else {
			val1, worked, err, indexcounted := translateVal(UNPARSEcode{code: str, realcode: code.realcode, line: code.line, path: code.path}, index, codelines, 0)
			if !worked || err.EXISTS {
				if i == len(trimmed)-1 {
					return val1, worked, err, i
				}
				continue
			}
			value = val1
			countIndex += indexcounted - 1
		}
		if LookingAtKey {
			currentKey = value
			current = i + 1
			LookingAtKey = false
		} else {
			body = append(body, [2]any{currentKey, value})
			current = i + 1
			LookingAtKey = true
		}
	}
	return createMap{
		body: body,
		code: code.realcode,
		line: code.line,
		path: code.path,
	}, true, ArErr{}, countIndex
}

var mutex = sync.RWMutex{}

func Map(m anymap) ArObject {
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
	obj.obj["keys"] = builtinFunc{
		"keys",
		func(args ...any) (any, ArErr) {
			if len(args) != 0 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected 0 arguments, got " + fmt.Sprint(len(args)),
					EXISTS:  true,
				}
			}
			mutex.RLock()
			keys := []any{}
			for k := range m {
				keys = append(keys, k)
			}
			mutex.RUnlock()
			return keys, ArErr{}
		},
	}
	obj.obj["__Boolean__"] = builtinFunc{
		"__Boolean__",
		func(args ...any) (any, ArErr) {
			mutex.RLock()
			if len(m) == 0 {
				mutex.RUnlock()
				return false, ArErr{}
			}
			mutex.RUnlock()
			return true, ArErr{}
		},
	}
	obj.obj["object"] = builtinFunc{
		"object",
		func(args ...any) (any, ArErr) {
			if len(args) != 0 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected 0 arguments, got " + fmt.Sprint(len(args)),
					EXISTS:  true,
				}
			}
			return ArObject{
				obj: m,
			}, ArErr{}
		},
	}
	return obj
}
