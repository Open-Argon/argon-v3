package main

import (
	"fmt"
	"strings"
)

type ArObject struct {
	TYPE string
	obj  anymap
}

type anymap map[any]any

var mapGetCompile = makeRegex(`(.|\n)+\.([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*( *)`)
var indexGetCompile = makeRegex(`(.|\n)+\[(.|\n)+\]( *)`)

type ArMapGet struct {
	VAL   any
	args  []any
	index bool
	line  int
	code  string
	path  string
}

func mapGet(r ArMapGet, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(r.VAL, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	switch m := resp.(type) {
	case ArObject:
		if r.index && m.TYPE != "map" {
			if _, ok := m.obj["__getindex__"]; ok {
				callable := m.obj["__getindex__"]
				resp, err := runCall(call{
					callable: callable,
					args:     r.args,
					line:     r.line,
					path:     r.path,
					code:     r.code,
				}, stack, stacklevel+1)
				if !err.EXISTS {
					return resp, ArErr{}
				}
			}
		}
		if len(r.args) > 1 {
			return nil, ArErr{
				"IndexError",
				"index not found",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		key, err := runVal(r.args[0], stack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		key = ArValidToAny(key)
		if isUnhashable(key) {
			return nil, ArErr{
				"TypeError",
				"unhashable type: '" + typeof(key) + "'",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		if _, ok := m.obj[key]; !ok {
			return nil, ArErr{
				"KeyError",
				"key '" + fmt.Sprint(key) + "' not found",
				r.line,
				r.path,
				r.code,
				true,
			}
		}
		return AnyToArValid(m.obj[key]), ArErr{}
	}

	key, err := runVal(r.args[0], stack, stacklevel+1)
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
	return ArMapGet{resp, []any{key}, false, code.line, code.realcode, code.path}, true, ArErr{}, 1
}

func isIndexGet(code UNPARSEcode) bool {
	return indexGetCompile.MatchString(code.code)
}

func indexGetParse(code UNPARSEcode, index int, codelines []UNPARSEcode) (ArMapGet, bool, ArErr, int) {
	trim := strings.TrimSpace(code.code)
	trim = trim[:len(trim)-1]
	split := strings.Split(trim, "[")
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
		return ArMapGet{tival, args, true, code.line, code.realcode, code.path}, true, ArErr{}, 1
	}
	return ArMapGet{}, false, ArErr{
		"Syntax Error",
		"invalid index get",
		code.line,
		code.path,
		code.realcode,
		true,
	}, 1
}

func isUnhashable(val any) bool {
	keytype := typeof(val)
	return keytype == "array" || keytype == "map"
}
