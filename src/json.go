package main

import (
	"encoding/json"
	"errors"
	"strconv"
)

func convertToArgon(obj any) any {
	switch x := obj.(type) {
	case map[string]interface{}:
		newmap := anymap{}
		for key, value := range x {
			newmap[key] = convertToArgon(value)
		}
		return Map(newmap)
	case []any:
		for i, value := range x {
			x[i] = convertToArgon(value)
		}
		return ArArray(x)
	case string:
		return ArString(x)
	case float64:
		return Number(x)
	case bool:
		return x
	case nil:
		return nil
	}
	return nil
}

func jsonparse(str string) (any, ArErr) {
	var jsonMap any
	var err = json.Unmarshal([]byte(str), &jsonMap)
	if err != nil {
		return nil, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
	}
	return convertToArgon(jsonMap), ArErr{}
}

func jsonstringify(obj any, level int64) (string, error) {
	if level > 100 {
		return "", errors.New("json stringify error: too many levels")
	}
	switch x := obj.(type) {
	case ArObject:
		if callable, ok := x.obj["__json__"]; ok {
			val, err := runCall(
				call{
					Callable: callable,
					Args:     []any{Int64ToNumber(level)},
				},
				stack{},
				0,
			)
			if err.EXISTS {
				return "", errors.New(err.message)
			}
			val = ArValidToAny(val)
			if x, ok := val.(string); ok {
				return x, nil
			} else {
				return "", errors.New("json stringify error: __json__ must return a string")
			}
		}
	case bool:
		return strconv.FormatBool(x), nil
	case nil:
		return "null", nil
	}
	err := errors.New("Cannot stringify '" + typeof(obj) + "'")
	return "", err
}

var ArJSON = Map(anymap{
	"parse": builtinFunc{"parse", func(args ...any) (any, ArErr) {
		if len(args) == 0 {
			return nil, ArErr{TYPE: "Runtime Error", message: "parse takes 1 argument", EXISTS: true}
		}
		if typeof(args[0]) != "string" {
			return nil, ArErr{TYPE: "Runtime Error", message: "parse takes a string not a '" + typeof(args[0]) + "'", EXISTS: true}
		}
		args[0] = ArValidToAny(args[0])
		return jsonparse(args[0].(string))
	}},
	"stringify": builtinFunc{"stringify", func(args ...any) (any, ArErr) {
		if len(args) != 1 && len(args) != 2 {
			return nil, ArErr{TYPE: "Runtime Error", message: "stringify takes 1 or 2 arguments", EXISTS: true}
		}
		var level int64 = 0
		if len(args) == 2 {
			if typeof(args[1]) != "number" {
				return nil, ArErr{TYPE: "Runtime Error", message: "stringify takes a number not a '" + typeof(args[1]) + "'", EXISTS: true}
			}
			var err error
			level, err = numberToInt64(args[1].(ArObject))
			if err != nil {
				return nil, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
		}
		str, err := jsonstringify(args[0], level)
		if err != nil {
			return nil, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
		}
		return ArString(str), ArErr{}
	}},
})
