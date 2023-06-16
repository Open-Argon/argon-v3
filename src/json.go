package main

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
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
		return x
	case float64:
		return newNumber().SetFloat64(x)
	case bool:
		return x
	case nil:
		return nil
	}
	return nil
}

func jsonparse(str string) any {
	var jsonMap any
	json.Unmarshal([]byte(str), &jsonMap)
	return convertToArgon(jsonMap)
}

func jsonstringify(obj any, level int) (string, error) {
	if level > 100 {
		return "", errors.New("json stringify error: too many levels")
	}
	output := []string{}
	obj = ArValidToAny(obj)
	switch x := obj.(type) {
	case anymap:
		for key, value := range x {
			str, err := jsonstringify(value, level+1)
			if err != nil {
				return "", err
			}
			output = append(output, ""+strconv.Quote(anyToArgon(key, false, true, 3, 0, false, 0))+": "+str)
		}
		return "{" + strings.Join(output, ", ") + "}", nil
	case []any:
		for _, value := range x {
			str, err := jsonstringify(value, level+1)
			if err != nil {
				return "", err
			}
			output = append(output, str)
		}
		return "[" + strings.Join(output, ", ") + "]", nil
	case string:
		return strconv.Quote(x), nil
	case number:
		num, _ := x.Float64()
		if math.IsNaN(num) || math.IsInf(num, 0) {
			return "null", nil
		}
		return numberToString(x, false), nil
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
		return jsonparse(args[0].(string)), ArErr{}
	}},
	"stringify": builtinFunc{"stringify", func(args ...any) (any, ArErr) {
		if len(args) == 0 {
			return nil, ArErr{TYPE: "Runtime Error", message: "stringify takes 1 argument", EXISTS: true}
		}
		str, err := jsonstringify(args[0], 0)
		if err != nil {
			return nil, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
		}
		return ArString(str), ArErr{}
	}},
})
