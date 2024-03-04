package main

import (
	"path/filepath"
	"strings"
)

var genericImportCompiled = makeRegex(`import( )+(.|\n)+(( )+as( )+([a-zA-Z_]|(\p{L}\p{M}*))([a-zA-Z0-9_]|(\p{L}\p{M}*))*)?( *)`)

type ArImport struct {
	pretranslated bool
	translated    translatedImport
	FilePath      any
	Values        any
	Code          string
	Line          int
	Path          string
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
			} else if after == "*" {
				asStr = true
			} else if variableCompile.MatchString(after) {
				asStr = after
			} else {
				return ArImport{}, false, ArErr{"Syntax Error", "invalid variable name '" + after + "'", code.line, code.path, code.realcode, true}, i
			}
			toImport = toImportval
		}
	}

	importOBJ := ArImport{
		false,
		translatedImport{},
		toImport,
		asStr,
		code.realcode,
		code.line,
		code.path,
	}

	if str, ok := toImport.(string); ok {
		importOBJ.pretranslated = true
		var err ArErr
		importOBJ.translated, err = translateImport(str, filepath.Dir(filepath.ToSlash(code.path)))
		if err.EXISTS {
			if err.line == 0 {
				err.line = importOBJ.Line
			}
			if err.path == "" {
				err.path = importOBJ.Path
			}
			if err.code == "" {
				err.code = importOBJ.Code
			}
			return importOBJ, false, err, i
		}
	}

	return importOBJ, true, ArErr{}, i
}

func runImport(importOBJ ArImport, stack stack, stacklevel int) (any, ArErr) {
	var translated = importOBJ.translated
	if !importOBJ.pretranslated {
		val, err := runVal(importOBJ.FilePath, stack, stacklevel+1)
		val = ArValidToAny(val)
		if err.EXISTS {
			return nil, err
		}
		if typeof(val) != "string" {
			return nil, ArErr{"Type Error", "import requires a string, got type '" + typeof(val) + "'", importOBJ.Line, importOBJ.Path, importOBJ.Code, true}
		}
		parent := filepath.Dir(filepath.ToSlash(importOBJ.Path))
		translated, err = translateImport(val.(string), parent)
		if err.EXISTS {
			if err.line == 0 {
				err.line = importOBJ.Line
			}
			if err.path == "" {
				err.path = importOBJ.Path
			}
			if err.code == "" {
				err.code = importOBJ.Code
			}
			return nil, err
		}
	}
	stackMap, err := runTranslatedImport(translated, stack[0], false)
	if err.EXISTS {
		if err.line == 0 {
			err.line = importOBJ.Line
		}
		if err.path == "" {
			err.path = importOBJ.Path
		}
		if err.code == "" {
			err.code = importOBJ.Code
		}
		return nil, err
	}
	setindex, ok := stack[len(stack)-1].obj["__setindex__"]
	if !ok {
		return nil, ArErr{
			"Import Error",
			"could not find __setindex__ in module scope",
			importOBJ.Line,
			importOBJ.Path,
			importOBJ.Code,
			true,
		}
	}
	switch x := importOBJ.Values.(type) {
	case []string:
		for _, v := range x {
			val, ok := stackMap.obj[v]
			if !ok {
				return nil, ArErr{"Import Error", "could not find value " + anyToArgon(v, true, false, 3, 0, false, 0) + " in module " + anyToArgon(translated.path, true, false, 3, 0, false, 0), importOBJ.Line, importOBJ.Path, importOBJ.Code, true}
			}
			builtinCall(setindex, []any{v, val})
		}
	case string:
		builtinCall(setindex, []any{x, stackMap})
	case bool:
		keyGetter, ok := stackMap.obj["keys"]
		if !ok {
			return nil, ArErr{"Import Error", "could not find keys in module scope", importOBJ.Line, importOBJ.Path, importOBJ.Code, true}
		}
		valueGetter, ok := stackMap.obj["__getindex__"]
		if !ok {
			return nil, ArErr{"Import Error", "could not find __getindex__ in module scope", importOBJ.Line, importOBJ.Path, importOBJ.Code, true}
		}
		keys, err := builtinCall(keyGetter, []any{})
		if err.EXISTS {
			return nil, err
		}
		keys = ArValidToAny(keys)
		for _, v := range keys.([]any) {
			val, err := builtinCall(valueGetter, []any{v})
			if err.EXISTS {
				return nil, err
			}
			builtinCall(setindex, []any{v, val})
		}
	}
	return nil, ArErr{}
}
