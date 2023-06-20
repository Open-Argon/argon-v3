package main

import (
	"os"
	"strings"
)

var genericImportCompiled = makeRegex(`import( )+(.|\n)+(( )+as( )+([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*)?( *)`)

type ArImport struct {
	filePath any
	values   any
	code     string
	line     int
	path     string
}

func isGenericImport(code UNPARSEcode) bool {
	return genericImportCompiled.MatchString(code.code)
}

func parseGenericImport(code UNPARSEcode, index int, codeline []UNPARSEcode) (ArImport, bool, ArErr, int) {
	trim := strings.Trim(code.code, " ")
	pathAndAs := trim[6:]
	split := strings.Split(pathAndAs, " as ")
	var toImport any
	var asStr any
	var i = 1
	if len(split) == 1 {
		toImportval, worked, err, I := translateVal(UNPARSEcode{
			code:     strings.Trim(split[0], " "),
			realcode: code.realcode,
			line:     code.line,
			path:     code.path,
		}, index, codeline, 0)
		if !worked || err.EXISTS {
			return ArImport{}, worked, err, I
		}
		toImport = toImportval
		i = I
	} else {
		for i := 1; i < len(split); i++ {
			before := strings.Trim(strings.Join(split[:i], " as "), " ")
			after := strings.Trim(strings.Join(split[i:], " as "), " ")
			toImportval, worked, err, I := translateVal(UNPARSEcode{
				code:     before,
				realcode: code.realcode,
				line:     code.line,
				path:     code.path,
			}, index, codeline, 0)
			i = I
			if !worked || err.EXISTS {
				if i == len(split)-1 {
					return ArImport{}, worked, err, i
				}
				continue
			}
			if after == "" {
			} else if variableCompile.MatchString(after) {
				asStr = after
			} else {
				return ArImport{}, false, ArErr{"Syntax Error", "invalid variable name '" + after + "'", code.line, code.path, code.realcode, true}, i
			}
			toImport = toImportval
		}
	}

	return ArImport{
		toImport,
		asStr,
		code.realcode,
		code.line,
		code.path,
	}, true, ArErr{}, i
}

func runImport(importOBJ ArImport, stack stack, stacklevel int) (any, ArErr) {
	val, err := runVal(importOBJ.filePath, stack, stacklevel+1)
	val = ArValidToAny(val)
	if err.EXISTS {
		return nil, err
	}
	if typeof(val) != "string" {
		return nil, ArErr{"Type Error", "import requires a string, got type '" + typeof(val) + "'", importOBJ.line, importOBJ.path, importOBJ.code, true}
	}
	path := val.(string)
	ex, e := os.Getwd()
	if e != nil {
		return nil, ArErr{"File Error", "could not get current working directory", importOBJ.line, importOBJ.path, importOBJ.code, true}
	}
	stackMap, err := importMod(path, ex, false, stack[0])
	if err.EXISTS {
		if err.line == 0 {
			err.line = importOBJ.line
		}
		if err.path == "" {
			err.path = importOBJ.path
		}
		if err.code == "" {
			err.code = importOBJ.code
		}
		return nil, err
	}
	setindex, ok := stack[len(stack)-1].obj["__setindex__"]
	if !ok {
		return nil, ArErr{
			"Import Error",
			"could not find __setindex__ in module scope",
			importOBJ.line,
			importOBJ.path,
			importOBJ.code,
			true,
		}
	}
	switch x := importOBJ.values.(type) {
	case []string:
		for _, v := range x {
			val, ok := stackMap.obj[v]
			if !ok {
				return nil, ArErr{"Import Error", "could not find value " + anyToArgon(v, true, false, 3, 0, false, 0) + " in module " + anyToArgon(path, true, false, 3, 0, false, 0), importOBJ.line, importOBJ.path, importOBJ.code, true}
			}
			builtinCall(setindex, []any{v, val})
		}
	case string:
		builtinCall(setindex, []any{x, stackMap})
	}
	return nil, ArErr{}
}
