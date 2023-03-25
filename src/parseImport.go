package main

import (
	"os"
	"strings"
)

var genericImportCompiled = makeRegex(`import( )+(.|\n)+( )+as( )+([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*( *)`)

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
	split := strings.SplitN(pathAndAs, " as ", 2)
	toImportstr := strings.TrimSpace(split[0])
	asStr := strings.TrimSpace(split[1])
	toImport, worked, err, i := translateVal(UNPARSEcode{
		code:     toImportstr,
		realcode: code.realcode,
		line:     code.line,
		path:     code.path,
	}, index, codeline, 0)
	if !worked {
		return ArImport{}, false, err, i
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
	return nil, ArErr{"Import Error", "importing in WASM is currently not supported", importOBJ.line, importOBJ.path, importOBJ.code, true}
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
	stackMap, err := importMod(path, ex, false)
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
	switch x := importOBJ.values.(type) {
	case []string:
		for _, v := range x {
			val, ok := stackMap.obj[v]
			if !ok {
				return nil, ArErr{"Import Error", "could not find value " + anyToArgon(v, true, false, 3, 0, false, 0) + " in module " + anyToArgon(path, true, false, 3, 0, false, 0), importOBJ.line, importOBJ.path, importOBJ.code, true}
			}
			stack[len(stack)-1].obj[v] = val
		}
	case string:
		stack[len(stack)-1].obj[x] = stackMap
	}
	return nil, ArErr{}
}
