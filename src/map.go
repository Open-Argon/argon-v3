package main

import (
	"fmt"
	"strings"
)

type ArMap = map[any]any

type ArClass struct {
	value any
	MAP   ArMap
}

var mapGetCompile = makeRegex(`(.|\n)+\.([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*( *)`)

type ArMapGet struct {
	VAL  any
	key  any
	line int
	code string
	path string
}

func mapGet(r ArMapGet, stack stack) (any, ArErr) {
	resp, err := runVal(r.VAL, stack)
	if err.EXISTS {
		return nil, err
	}
	key, err := runVal(r.key, stack)
	if err.EXISTS {
		return nil, err
	}
	switch m := resp.(type) {
	case ArMap:
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
	case ArClass:
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

	}
	return nil, ArErr{
		"TypeError",
		"cannot read " + anyToArgon(key, true, true, 3, 0) + " from type '" + typeof(resp) + "'",
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
	resp, worked, err, i := translateVal(UNPARSEcode{code: start, realcode: code.realcode, line: code.line, path: code.path}, index, codelines, false)
	if !worked {
		return ArMapGet{}, false, err, i
	}
	k := key
	return ArMapGet{resp, k, code.line, code.realcode, code.path}, true, ArErr{}, 1
}
