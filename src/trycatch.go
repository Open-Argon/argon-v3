package main

import (
	"strings"
)

var tryCompiled = makeRegex(`\s*try\s(.|\n)+`)
var catchCompiled = makeRegex(`\s*catch\s*\(\s*` + spacelessVariable + `\s*\)\s*(.|\n)+`)

type TryCatch struct {
	Try       any
	Catch     any
	errorName string
	line      int
	path      string
	code      string
}

func isTryCatch(code UNPARSEcode) bool {
	return tryCompiled.MatchString(code.code)
}

func parseTryCatch(code UNPARSEcode, index int, codelines []UNPARSEcode) (TryCatch, bool, ArErr, int) {
	trytrimmed := strings.TrimSpace(code.code)
	totalIndex := 0
	tryparsed, worked, err, i := translateVal(UNPARSEcode{trytrimmed[4:], code.realcode, code.line, code.path}, index, codelines, 1)
	if !worked {
		return TryCatch{}, false, err, i
	}
	totalIndex += i

	if index+totalIndex >= len(codelines) {
		return TryCatch{}, false, ArErr{"Syntax Error", "expected catch statement", code.line, code.path, code.realcode, true}, i
	}
	catchtrimmed := strings.TrimSpace(codelines[index+totalIndex].code)
	if !catchCompiled.MatchString(catchtrimmed) {
		return TryCatch{}, false, ArErr{"Syntax Error", "invalid syntax", code.line, code.path, code.realcode, true}, i
	}
	catchtrimmed = catchtrimmed[6:]
	catchbracketSplit := strings.SplitN(catchtrimmed, ")", 2)
	errorName := strings.TrimSpace(strings.TrimSpace(catchbracketSplit[0])[1:])
	errcode := catchbracketSplit[1]
	catchparsed, worked, err, i := translateVal(UNPARSEcode{errcode, code.realcode, code.line, code.path}, index+totalIndex, codelines, 1)
	if !worked {
		return TryCatch{}, false, err, i
	}
	totalIndex += i

	return TryCatch{
		tryparsed,
		catchparsed,
		errorName,
		code.line,
		code.path,
		code.realcode,
	}, true, ArErr{}, totalIndex
}

func runTryCatch(t TryCatch, stack stack, stacklevel int) (any, ArErr) {
	val, err := runVal(t.Try, stack, stacklevel+1)
	if err.EXISTS {
		vars := anymap{}
		vars[t.errorName] = Map(anymap{
			"type":    err.TYPE,
			"message": err.message,
			"line":    newNumber().SetInt64(int64(err.line)),
			"path":    err.path,
			"code":    err.code,
		})
		val, err = runVal(t.Catch, append(stack, Map(vars)), stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
	}
	return val, ArErr{}
}
