package main

import (
	"strings"
)

type ArObject struct {
	obj anymap
}

type anymap map[any]any

var mapGetCompile = makeRegex(`(.|\n)+\.([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*( *)`)
var indexGetCompile = makeRegex(`(.|\n)+\[(.|\n)+\]( *)`)

type ArMapGet struct {
	VAL                any
	Args               []any
	IncludeConstuctors bool
	Line               int
	Code               string
	Path               string
}

func isObject(val any) bool {
	if _, ok := val.(ArObject); ok {
		return true
	}
	return false
}

func hashableObject(obj ArObject) (string, ArErr) {
	if callable, ok := obj.obj["__hash__"]; ok {
		resp, err := runCall(call{
			Callable: callable,
			Args:     []any{},
		}, stack{}, 0)
		if err.EXISTS {
			return "", err
		}
		resp = ArValidToAny(resp)
		if str, ok := resp.(string); ok {
			return str, ArErr{}
		}
		return "", ArErr{
			TYPE:    "TypeError",
			EXISTS:  true,
			message: "expected string from __hash__ method, got " + typeof(resp),
		}
	}
	return "", ArErr{
		TYPE:    "TypeError",
		EXISTS:  true,
		message: "cannot hash object",
	}
}

func mapGet(r ArMapGet, stack stack, stacklevel int) (any, ArErr) {
	resp, err := runVal(r.VAL, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	switch m := resp.(type) {
	case ArObject:
		if r.IncludeConstuctors {
			if obj, ok := m.obj[r.Args[0]]; ok {
				return obj, ArErr{}
			}
		}
		if callable, ok := m.obj["__getindex__"]; ok {
			resp, err := runCall(call{
				Callable: callable,
				Args:     r.Args,
				Line:     r.Line,
				Path:     r.Path,
				Code:     r.Code,
			}, stack, stacklevel+1)
			return resp, err
		}
	}

	key, err := runVal(r.Args[0], stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	return nil, ArErr{
		"Type Error",
		"cannot read " + anyToArgon(key, true, true, 3, 0, false, 0) + " from type '" + typeof(resp) + "'",
		r.Line,
		r.Path,
		r.Code,
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
	return ArMapGet{resp, []any{key}, true, code.line, code.realcode, code.path}, true, ArErr{}, 1
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
		return ArMapGet{tival, args, false, code.line, code.realcode, code.path}, true, ArErr{}, 1
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

var hashabletypes = []string{
	"number",
	"string",
	"bool",
	"null",
}

func isUnhashable(val any) bool {
	keytype := typeof(val)
	for _, v := range hashabletypes {
		if v == keytype {
			return false
		}
	}
	return true
}
