package main

import (
	"fmt"
	"strings"
)

type ArMap = map[any]any
type ArArray = []any

type ArClass struct {
	value any
	MAP   ArMap
}

var mapGetCompile = makeRegex(`(.|\n)+\.([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*( *)`)
var indexGetCompile = makeRegex(`(.|\n)+\[(.|\n)+\]( *)`)

type ArMapGet struct {
	VAL           any
	start         any
	end           any
	step          any
	index         bool
	numberofindex int
	line          int
	code          string
	path          string
}

func mapGet(r ArMapGet, stack stack) (any, ArErr) {
	resp, err := runVal(r.VAL, stack)
	if err.EXISTS {
		return nil, err
	}
	switch m := resp.(type) {
	case ArMap:
		if r.numberofindex > 1 {
			return nil, ArErr{
				"IndexError",
				"index not found",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		key, err := runVal(r.start, stack)
		if err.EXISTS {
			return nil, err
		}
		if _, ok := m[key]; !ok {
			return nil, ArErr{
				"KeyError",
				"key '" + fmt.Sprint(key) + "' not found",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		return m[key], ArErr{}

	case ArArray:
		startindex := 0
		endindex := 1
		step := 1

		if !r.index {
			key, err := runVal(r.start, stack)
			if err.EXISTS {
				return nil, err
			}
			if key == "length" {
				return len(m), ArErr{}
			}
		}
		if r.start != nil {
			sindex, err := runVal(r.start, stack)
			if err.EXISTS {
				return nil, err
			}
			if typeof(sindex) != "number" {
				return nil, ArErr{
					"TypeError",
					"index must be a number",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			num := sindex.(number)
			if !num.IsInt() {
				return nil, ArErr{
					"TypeError",
					"index must be an integer",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			startindex = int(num.Num().Int64())
			endindex = startindex + 1
		}
		if r.end != nil {
			eindex, err := runVal(r.end, stack)
			if err.EXISTS {
				return nil, err
			}
			if typeof(eindex) != "number" {
				return nil, ArErr{
					"TypeError",
					"ending index must be a number",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			num := eindex.(number)
			if !num.IsInt() {
				return nil, ArErr{
					"TypeError",
					"ending index must be an integer",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			endindex = int(num.Num().Int64())
		} else if r.numberofindex > 1 {
			endindex = len(m)
		}
		if r.step != nil {
			step, err := runVal(r.step, stack)
			if err.EXISTS {
				return nil, err
			}
			if typeof(step) != "number" {
				return nil, ArErr{
					"TypeError",
					"step must be a number",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			num := step.(number)
			if !num.IsInt() {
				return nil, ArErr{
					"TypeError",
					"step must be an integer",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			step = int(num.Num().Int64())
		}
		if startindex < 0 {
			startindex = len(m) + startindex
		}
		if endindex < 0 {
			endindex = len(m) + endindex
		}
		if step < 0 {
			step = -step
			startindex, endindex = endindex, startindex
		}
		if startindex < 0 || startindex >= len(m) {
			return nil, ArErr{
				"IndexError",
				"index '" + fmt.Sprint(startindex) + "' out of range",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		if endindex < 0 || endindex > len(m) {
			return nil, ArErr{
				"IndexError",
				"index '" + fmt.Sprint(endindex) + "' out of range",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		if step == 0 {
			return nil, ArErr{
				"ValueError",
				"step cannot be 0",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		return m[startindex:endindex:step], ArErr{}
	case ArClass:
		if r.numberofindex > 1 {
			return nil, ArErr{
				"IndexError",
				"index not found",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		key, err := runVal(r.start, stack)
		if err.EXISTS {
			return nil, err
		}
		if _, ok := m.MAP[key]; !ok {
			return nil, ArErr{
				"KeyError",
				"key '" + fmt.Sprint(key) + "' not found",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		return m.MAP[key], ArErr{}
	case string:
		startindex := 0
		endindex := 1
		step := 1

		if !r.index {
			key, err := runVal(r.start, stack)
			if err.EXISTS {
				return nil, err
			}
			if key == "length" {
				return len(m), ArErr{}
			}
		}
		if r.start != nil {
			sindex, err := runVal(r.start, stack)
			if err.EXISTS {
				return nil, err
			}
			if typeof(sindex) != "number" {
				return nil, ArErr{
					"TypeError",
					"index must be a number",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			num := sindex.(number)
			if !num.IsInt() {
				return nil, ArErr{
					"TypeError",
					"index must be an integer",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			startindex = int(num.Num().Int64())
			endindex = startindex + 1
		}
		if r.end != nil {
			eindex, err := runVal(r.end, stack)
			if err.EXISTS {
				return nil, err
			}
			if typeof(eindex) != "number" {
				return nil, ArErr{
					"TypeError",
					"ending index must be a number",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			num := eindex.(number)
			if !num.IsInt() {
				return nil, ArErr{
					"TypeError",
					"ending index must be an integer",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			endindex = int(num.Num().Int64())
		} else if r.numberofindex > 1 {
			endindex = len(m)
		}
		if r.step != nil {
			step, err := runVal(r.step, stack)
			if err.EXISTS {
				return nil, err
			}
			if typeof(step) != "number" {
				return nil, ArErr{
					"TypeError",
					"step must be a number",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			num := step.(number)
			if !num.IsInt() {
				return nil, ArErr{
					"TypeError",
					"step must be an integer",
					r.line,
					r.path,
					r.code,
					true,
				}
			}
			step = int(num.Num().Int64())
		}
		if startindex < 0 {
			startindex = len(m) + startindex
		}
		if endindex < 0 {
			endindex = len(m) + endindex
		}
		if step < 0 {
			step = -step
			startindex, endindex = endindex, startindex
		}
		if startindex < 0 || startindex >= len(m) {
			return nil, ArErr{
				"IndexError",
				"index '" + fmt.Sprint(startindex) + "' out of range",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		if endindex < 0 || endindex > len(m) {
			return nil, ArErr{
				"IndexError",
				"index '" + fmt.Sprint(endindex) + "' out of range",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		if step == 0 {
			return nil, ArErr{
				"ValueError",
				"step cannot be 0",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		return string(([]byte(m))[startindex:endindex:step]), ArErr{}
	}

	key, err := runVal(r.start, stack)
	if err.EXISTS {
		return nil, err
	}
	return nil, ArErr{
		"TypeError",
		"cannot read " + anyToArgon(key, true, true, 3, 0, false, 0) + " from type '" + typeof(resp) + "'",
		r.line,
		r.path,
		r.code,
		true,
	}
}

func classVal(r any) any {
	if _, ok := r.(ArClass); ok {
		return r.(ArClass).value
	}
	return r
}

func isMapGet(code UNPARSEcode) bool {
	return mapGetCompile.MatchString(code.code)
}

func mapGetParse(code UNPARSEcode, index int, codelines []UNPARSEcode) (ArMapGet, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	split := strings.Split(trim, ".")
	start := strings.Join(split[:len(split)-1], ".")
	key := split[len(split)-1]
	resp, worked, err, i := translateVal(UNPARSEcode{code: start, realcode: code.realcode, line: code.line, path: code.path}, index, codelines, 0)
	if !worked {
		return ArMapGet{}, false, err, i
	}
	k := key
	return ArMapGet{resp, k, nil, nil, false, 1, code.line, code.realcode, code.path}, true, ArErr{}, 1
}

func isIndexGet(code UNPARSEcode) bool {
	return indexGetCompile.MatchString(code.code)
}

func indexGetParse(code UNPARSEcode, index int, codelines []UNPARSEcode) (ArMapGet, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	trim = trim[:len(trim)-1]
	split := strings.Split(trim, "[")
	var toindex any
	var start any
	var end any
	var step any
	numberofindexs := 0
	for i := 1; i < len(split); i++ {
		ti := strings.Join(split[:i], "[")
		innerbrackets := strings.Join(split[i:], "[")
		args, success, argserr := getValuesFromLetter(innerbrackets, ":", index, codelines, true)
		if !success {
			if i == len(split)-1 {
				return ArMapGet{}, false, argserr, 1
			}
			continue
		}
		if len(args) > 3 {
			return ArMapGet{}, false, ArErr{
				"SyntaxError",
				"too many arguments for index get",
				code.line,
				code.path,
				code.realcode,
				true,
			}, 1
		}
		tival, worked, err, i := translateVal(UNPARSEcode{code: ti, realcode: code.realcode, line: code.line, path: code.path}, index, codelines, 0)
		if !worked {
			if i == len(split)-1 {
				return ArMapGet{}, false, err, i
			}
			continue
		}
		numberofindexs = len(args)
		if len(args) >= 1 {
			toindex = tival
			start = args[0]
		}
		if len(args) >= 2 {
			end = args[1]
		}
		if len(args) >= 3 {
			step = args[2]
		}
	}
	if toindex == nil {
		return ArMapGet{}, false, ArErr{
			"SyntaxError",
			"invalid index get",
			code.line,
			code.path,
			code.realcode,
			true,
		}, 1
	}
	return ArMapGet{toindex, start, end, step, true, numberofindexs, code.line, code.realcode, code.path}, true, ArErr{}, 1
}
