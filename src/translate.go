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

// returns (number | string | nil), success, error, step
func translateVal(code UNPARSEcode, index int, codelines []UNPARSEcode, isLine int) (any, bool, ArErr, int) {

	if isLine == 2 {
		if isDeleteVariable(code) {
			return parseDelete(code, index, codelines)
		} else if isComment(code) {
			resp, worked, err, step := parseComment(code, index, codelines)
			if worked {
				return resp, worked, err, step
			}
		} else if isReturn(code) {
			return parseReturn(code, index, codelines)
		} else if isBreak(code) {
			return parseBreak(code, index, codelines)
		} else if isIfStatement(code) {
			return parseIfStatement(code, index, codelines)
		} else if isWhileLoop(code) {
			return parseWhileLoop(code, index, codelines)
		} else if isForeverLoop(code) {
			return parseForeverLoop(code, index, codelines)
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

	if isBrackets(code) {
		bracket, worked, err, step := parseBrackets(code, index, codelines)
		if worked {
			return bracket, worked, err, step
		}
	} else if isnot(code) {
		return parseNot(code, index, codelines, isLine)
	}
	if isSetVariable(code) {
		setvar, worked, err, step := parseSetVariable(code, index, codelines, isLine)
		if worked {
			return setvar, worked, err, step
		}
	}
	if isAutoAsignVariable(code) {
		setvar, worked, err, step := parseAutoAsignVariable(code, index, codelines, isLine)
		if worked {
			return setvar, worked, err, step
		}
	}
	operation, worked, err, step := parseOperations(code, index, codelines)
	if worked {
		return operation, worked, err, step
	} else if err.EXISTS {
		return nil, worked, err, step
	}
	if isNumber(code) {
		return parseNumber(code)
	} else if isNegative(code) {
		return parseNegative(code, index, codelines)
	} else if isFactorial(code) {
		return parseFactorial(code, index, codelines)
	} else if isCall(code) {
		call, worked, err, step := parseCall(code, index, codelines)
		if worked {
			return call, worked, err, step
		}
	}
	if isBoolean(code) {
		return parseBoolean(code)
	} else if isVariable(code) {
		return parseVariable(code)
	} else if isMapGet(code) {
		return mapGetParse(code, index, codelines)
	} else if isIndexGet(code) {
		return indexGetParse(code, index, codelines)
	} else if isString(code) {
		return parseString(code)
	}
	return nil, false, ArErr{"Syntax Error", "invalid syntax", code.line, code.path, code.realcode, true}, 1
}

// returns [](number | string), error
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
