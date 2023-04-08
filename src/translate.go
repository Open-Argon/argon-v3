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

func translateVal(code UNPARSEcode, index int, codelines []UNPARSEcode, isLine int) (any, bool, ArErr, int) {
	var (
		resp   any   = nil
		worked bool  = false
		err    ArErr = ArErr{"Syntax Error", "invalid syntax", code.line, code.path, code.realcode, true}
		i      int   = 1
	)
	if isLine == 2 {
		if isDeleteVariable(code) {
			return parseDelete(code, index, codelines)
		} else if isComment(code) {
			resp, worked, err, i = parseComment(code, index, codelines)
			if worked {
				return resp, worked, err, i
			}
		} else if isReturn(code) {
			return parseReturn(code, index, codelines)
		} else if isBreak(code) {
			return parseBreak(code)
		} else if isContinue(code) {
			return parseContinue(code)
		} else if isIfStatement(code) {
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
	if isLine >= 1 {
		if isDoWrap(code) {
			return parseDoWrap(code, index, codelines)
		}
	}

	if isLine == 2 {
		isLine = 1
	}

	if isBoolean(code) {
		return parseBoolean(code)
	} else if isBrackets(code) {
		resp, worked, err, i = parseBrackets(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
	}
	if isAbs(code) {
		resp, worked, err, i = parseAbs(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
	}
	if isnot(code) {
		return parseNot(code, index, codelines, isLine)
	}
	if isSetVariable(code) {
		resp, worked, err, i = parseSetVariable(code, index, codelines, isLine)
		if worked {
			return resp, worked, err, i
		}
	}
	if isAutoAsignVariable(code) {
		resp, worked, err, i = parseAutoAsignVariable(code, index, codelines, isLine)
		if worked {
			return resp, worked, err, i
		}
	}
	if isNumber(code) {
		return parseNumber(code)
	} else if isString(code) {
		return parseString(code)
	} else if issquareroot(code) {
		return parseSquareroot(code, index, codelines)
	} else if isFactorial(code) {
		return parseFactorial(code, index, codelines)
	}
	if isVariable(code) {
		return parseVariable(code)
	}
	if isArray(code) {
		resp, worked, err, i = parseArray(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
	}
	if isCall(code) {
		resp, worked, err, i = parseCall(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
	}
	{
		operation, worked, err, step := parseOperations(code, index, codelines)
		if worked {
			return operation, worked, err, step
		} else if err.EXISTS {
			return nil, worked, err, step
		}
	}
	if isNegative(code) {
		return parseNegative(code, index, codelines)
	} else if isMapGet(code) {
		return mapGetParse(code, index, codelines)
	} else if isIndexGet(code) {
		resp, worked, err, i = indexGetParse(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
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
		val, _, err, step := translateVal(codelines[i], i, codelines, 2)
		switch val.(type) {
		case CallReturn:
			return nil, ArErr{"Runtime Error", "Jump statment at top level", codelines[i].line, codelines[i].path, codelines[i].realcode, true}
		}
		i += step
		if err.EXISTS {
			return nil, err
		}
		translated = append(translated, val)
	}
	return translated, ArErr{}
}
