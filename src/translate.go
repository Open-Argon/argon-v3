package main

import (
	"strings"
)

type UNPARSEcode struct {
	code     string
	realcode string
	line     int
	path     string
}

var QuickKnownFailures = map[string]bool{}

func translateVal(code UNPARSEcode, index int, codelines []UNPARSEcode, isLine int) (any, bool, ArErr, int) {
	var (
		resp   any   = nil
		worked bool  = false
		err    ArErr = ArErr{"Syntax Error", "invalid syntax", code.line, code.path, code.realcode, true}
		i      int   = 1
	)
	if isLine >= 3 {
		if isComment(code) {
			resp, worked, err, i = parseComment(code, index, codelines)
			if worked {
				return resp, worked, err, i
			}
		}
		if isIfStatement(code) {
			return parseIfStatement(code, index, codelines)
		} else if isWhileLoop(code) {
			return parseWhileLoop(code, index, codelines)
		} else if isForeverLoop(code) {
			return parseForeverLoop(code, index, codelines)

		} else if isForLoop(code) {
			return parseForLoop(code, index, codelines)

		} else if isGenericImport(code) {
			return parseGenericImport(code, index, codelines)

		} else if isTryCatch(code) {
			return parseTryCatch(code, index, codelines)
		}
	}

	if isLine >= 2 {
		if isReturn(code) {
			return parseReturn(code, index, codelines)

		} else if isBreak(code) {
			return parseBreak(code)

		} else if isContinue(code) {
			return parseContinue(code)

		} else if isDeleteVariable(code) {
			return parseDelete(code, index, codelines)

		}
	}

	if isLine > 1 {
		isLine = 1
	}

	if isDoWrap(code) {
		return parseDoWrap(code, index, codelines)
	} else if isBoolean(code) {
		return parseBoolean(code)
	} else if !QuickKnownFailures["brackets"+code.code] && isBrackets(code) {
		resp, worked, err, i = parseBrackets(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["brackets"+code.code] = true
	}
	if !QuickKnownFailures["abs"+code.code] && isAbs(code) {
		resp, worked, err, i = parseAbs(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["abs"+code.code] = true
	}
	if !QuickKnownFailures["autoasign"+code.code] && isAutoAsignVariable(code) {
		resp, worked, err, i = parseAutoAsignVariable(code, index, codelines, isLine)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["autoasign"+code.code] = true
	}
	if isSetVariable(code) {
		return parseSetVariable(code, index, codelines, isLine)
	} else if isNumber(code) {
		return parseNumber(code)
	} else if isString(code) {
		return parseString(code)
	} else if !QuickKnownFailures["squareroot"+code.code] && issquareroot(code) {
		resp, worked, err, i = parseSquareroot(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["squareroot"+code.code] = true
	}
	if !QuickKnownFailures["factorial"+code.code] && isFactorial(code) {
		resp, worked, err, i = parseFactorial(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["factorial"+code.code] = true
	}
	if isVariable(code) {
		return parseVariable(code)
	}
	if !QuickKnownFailures["array"+code.code] && isArray(code) {
		resp, worked, err, i = parseArray(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["array"+code.code] = true
	} else if isMap(code) {
		resp, worked, err, i = parseMap(code, index, codelines)
	}
	if !QuickKnownFailures["not"+code.code] && isnot(code) {
		resp, worked, err, i = parseNot(code, index, codelines, isLine)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["not"+code.code] = true
	}
	if !QuickKnownFailures["operations"+code.code] {
		operation, worked, err, step := parseOperations(code, index, codelines)
		if worked {
			return operation, worked, err, step
		}
		QuickKnownFailures["operations"+code.code] = true
		if err.EXISTS {
			return nil, worked, err, step
		}
	}
	if isNegative(code) {
		return parseNegative(code, index, codelines)
	}
	if !QuickKnownFailures["call"+code.code] && isCall(code) {
		resp, worked, err, i = parseCall(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["call"+code.code] = true
	}
	if isMapGet(code) {
		return mapGetParse(code, index, codelines)
	} else if !QuickKnownFailures["indexget"+code.code] && isIndexGet(code) {
		resp, worked, err, i = indexGetParse(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["indexget"+code.code] = true
	}

	return resp, worked, err, i
}

func translate(codelines []UNPARSEcode) ([]any, ArErr) {
	translated := []any{}
	for i := 0; i < len(codelines); {
		if isBlank(codelines[i]) {
			i++
			continue
		}
		currentindent := len(codelines[i].code) - len(strings.TrimLeft(codelines[i].code, " "))
		if currentindent != 0 {
			return nil, ArErr{"Syntax Error", "invalid indent", codelines[i].line, codelines[i].path, codelines[i].realcode, true}
		}
		val, _, err, step := translateVal(codelines[i], i, codelines, 4)
		i += step
		if err.EXISTS {
			return nil, err
		}
		err = translateThrowOnNonLoop(val)
		if err.EXISTS {
			return nil, err
		}

		translated = append(translated, val)
	}
	return translated, ArErr{}
}
