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

var knownFailures = map[string]ArErr{}
var QuickKnownFailures = map[string]bool{}

func translateVal(code UNPARSEcode, index int, codelines []UNPARSEcode, isLine int) (any, bool, ArErr, int) {
	if knownErr, ok := knownFailures[code.code]; ok {
		return nil, false, ArErr{
			knownErr.TYPE,
			knownErr.message,
			code.line,
			code.path,
			code.realcode,
			true,
		}, 1
	}
	var (
		resp   any   = nil
		worked bool  = false
		err    ArErr = ArErr{"Syntax Error", "invalid syntax", code.line, code.path, code.realcode, true}
		i      int   = 1
	)
	if isLine == 3 {
		if isComment(code) {
			resp, worked, err, i = parseComment(code, index, codelines)
			if worked {
				return resp, worked, err, i
			}
		}
		if isIfStatement(code) {
			resp, worked, err, i = parseIfStatement(code, index, codelines)
			if !worked {
				knownFailures[code.code] = err
			}
			return resp, worked, err, i
		} else if isWhileLoop(code) {
			resp, worked, err, i = parseWhileLoop(code, index, codelines)
			if !worked {
				knownFailures[code.code] = err
			}
			return resp, worked, err, i
		} else if isForeverLoop(code) {
			resp, worked, err, i = parseForeverLoop(code, index, codelines)
			if !worked {
				knownFailures[code.code] = err
			}
			return resp, worked, err, i
		} else if isForLoop(code) {
			resp, worked, err, i = parseForLoop(code, index, codelines)
			if !worked {
				knownFailures[code.code] = err
			}
			return resp, worked, err, i
		} else if isGenericImport(code) {
			resp, worked, err, i = parseGenericImport(code, index, codelines)
			if !worked {
				knownFailures[code.code] = err
			}
			return resp, worked, err, i
		} else if isTryCatch(code) {
			resp, worked, err, i = parseTryCatch(code, index, codelines)
			if !worked {
				knownFailures[code.code] = err
			}
			return resp, worked, err, i
		}
	}

	if isLine >= 2 {
		if isReturn(code) {
			resp, worked, err, i = parseReturn(code, index, codelines)
			if !worked {
				knownFailures[code.code] = err
			}
			return resp, worked, err, i
		} else if isBreak(code) {
			resp, worked, err, i = parseBreak(code)
			if !worked {
				knownFailures[code.code] = err
			}
			return resp, worked, err, i
		} else if isContinue(code) {
			resp, worked, err, i = parseContinue(code)
			if !worked {
				knownFailures[code.code] = err
			}
			return resp, worked, err, i
		} else if isDeleteVariable(code) {
			resp, worked, err, i = parseDelete(code, index, codelines)
			if !worked {
				knownFailures[code.code] = err
			}
			return resp, worked, err, i
		}
	}

	if isLine >= 1 {
		if isDoWrap(code) {
			resp, worked, err, i = parseDoWrap(code, index, codelines)
			if !worked {
				knownFailures[code.code] = err
			}
			return resp, worked, err, i
		}
	}

	if isLine > 1 {
		isLine = 1
	}

	if isBoolean(code) {
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
	if !QuickKnownFailures["setvar"+code.code] && isSetVariable(code) {
		resp, worked, err, i = parseSetVariable(code, index, codelines, isLine)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["setvar"+code.code] = true
	}
	if !QuickKnownFailures["autoasign"+code.code] && isAutoAsignVariable(code) {
		resp, worked, err, i = parseAutoAsignVariable(code, index, codelines, isLine)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["autoasign"+code.code] = true
	}
	if isNumber(code) {
		resp, worked, err, i = parseNumber(code)
		if !worked {
			knownFailures[code.code] = err
		}
		return resp, worked, err, i
	} else if isString(code) {
		resp, worked, err, i = parseString(code)
		if !worked {
			knownFailures[code.code] = err
		}
		return resp, worked, err, i
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
		resp, worked, err, i = parseVariable(code)
		if !worked {
			knownFailures[code.code] = err
		}
		return resp, worked, err, i
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
	if !QuickKnownFailures["call"+code.code] && isCall(code) {
		resp, worked, err, i = parseCall(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["call"+code.code] = true
	}
	if isNegative(code) {
		resp, worked, err, i = parseNegative(code, index, codelines)
		if !worked {
			knownFailures[code.code] = err
		}
		return resp, worked, err, i
	} else if isMapGet(code) {
		resp, worked, err, i = mapGetParse(code, index, codelines)
		if !worked {
			knownFailures[code.code] = err
		}
		return resp, worked, err, i
	} else if !QuickKnownFailures["indexget"+code.code] && isIndexGet(code) {
		resp, worked, err, i = indexGetParse(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
		QuickKnownFailures["indexget"+code.code] = true
	}
	if !worked {
		knownFailures[code.code] = err
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
		val, _, err, step := translateVal(codelines[i], i, codelines, 3)
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
