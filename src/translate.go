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

var knownFailures = []string{}
var knownFailuresErrs = []ArErr{}

func StringExists(arr []string, target string) (bool, ArErr) {
	for i, str := range arr {
		if str == target {
			return true, knownFailuresErrs[i]
		}
	}
	return false, ArErr{}
}

func translateVal(code UNPARSEcode, index int, codelines []UNPARSEcode, isLine int) (any, bool, ArErr, int) {
	known, knownErr := StringExists(knownFailures, code.code)
	if known {
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
	if isLine == 2 {
		if isDeleteVariable(code) {
			resp, worked, err, i = parseDelete(code, index, codelines)
			if !worked {
				knownFailures = append(knownFailures, code.code)
				knownFailuresErrs = append(knownFailuresErrs, err)
			}
			return resp, worked, err, i
		} else if isComment(code) {
			resp, worked, err, i = parseComment(code, index, codelines)
			if worked {
				return resp, worked, err, i
			}
		} else if isReturn(code) {
			resp, worked, err, i = parseReturn(code, index, codelines)
			if !worked {
				knownFailures = append(knownFailures, code.code)
				knownFailuresErrs = append(knownFailuresErrs, err)
			}
			return resp, worked, err, i
		} else if isBreak(code) {
			resp, worked, err, i = parseBreak(code)
			if !worked {
				knownFailures = append(knownFailures, code.code)
				knownFailuresErrs = append(knownFailuresErrs, err)
			}
			return resp, worked, err, i
		} else if isContinue(code) {
			resp, worked, err, i = parseContinue(code)
			if !worked {
				knownFailures = append(knownFailures, code.code)
				knownFailuresErrs = append(knownFailuresErrs, err)
			}
			return resp, worked, err, i
		} else if isIfStatement(code) {
			resp, worked, err, i = parseIfStatement(code, index, codelines)
			if !worked {
				knownFailures = append(knownFailures, code.code)
				knownFailuresErrs = append(knownFailuresErrs, err)
			}
			return resp, worked, err, i
		} else if isWhileLoop(code) {
			resp, worked, err, i = parseWhileLoop(code, index, codelines)
			if !worked {
				knownFailures = append(knownFailures, code.code)
				knownFailuresErrs = append(knownFailuresErrs, err)
			}
			return resp, worked, err, i
		} else if isForeverLoop(code) {
			resp, worked, err, i = parseForeverLoop(code, index, codelines)
			if !worked {
				knownFailures = append(knownFailures, code.code)
				knownFailuresErrs = append(knownFailuresErrs, err)
			}
			return resp, worked, err, i
		} else if isForLoop(code) {
			resp, worked, err, i = parseForLoop(code, index, codelines)
			if !worked {
				knownFailures = append(knownFailures, code.code)
				knownFailuresErrs = append(knownFailuresErrs, err)
			}
			return resp, worked, err, i
		} else if isGenericImport(code) {
			resp, worked, err, i = parseGenericImport(code, index, codelines)
			if !worked {
				knownFailures = append(knownFailures, code.code)
				knownFailuresErrs = append(knownFailuresErrs, err)
			}
			return resp, worked, err, i
		} else if isTryCatch(code) {
			resp, worked, err, i = parseTryCatch(code, index, codelines)
			if !worked {
				knownFailures = append(knownFailures, code.code)
				knownFailuresErrs = append(knownFailuresErrs, err)
			}
			return resp, worked, err, i
		}
	}

	if isLine >= 1 {
		if isDoWrap(code) {
			resp, worked, err, i = parseDoWrap(code, index, codelines)
			if !worked {
				knownFailures = append(knownFailures, code.code)
				knownFailuresErrs = append(knownFailuresErrs, err)
			}
			return resp, worked, err, i
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
		resp, worked, err, i = parseNumber(code)
		if !worked {
			knownFailures = append(knownFailures, code.code)
			knownFailuresErrs = append(knownFailuresErrs, err)
		}
		return resp, worked, err, i
	} else if isString(code) {
		resp, worked, err, i = parseString(code)
		if !worked {
			knownFailures = append(knownFailures, code.code)
			knownFailuresErrs = append(knownFailuresErrs, err)
		}
		return resp, worked, err, i
	} else if issquareroot(code) {
		resp, worked, err, i = parseSquareroot(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
	}
	if isFactorial(code) {
		resp, worked, err, i = parseFactorial(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
	}
	if isVariable(code) {
		resp, worked, err, i = parseVariable(code)
		if !worked {
			knownFailures = append(knownFailures, code.code)
			knownFailuresErrs = append(knownFailuresErrs, err)
		}
		return resp, worked, err, i
	}
	if isArray(code) {
		resp, worked, err, i = parseArray(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
	} else if isMap(code) {
		resp, worked, err, i = parseMap(code, index, codelines)
	}
	if isnot(code) {
		resp, worked, err, i = parseNot(code, index, codelines, isLine)
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
	if isCall(code) {
		resp, worked, err, i = parseCall(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
	}
	if isNegative(code) {
		resp, worked, err, i = parseNegative(code, index, codelines)
		if !worked {
			knownFailures = append(knownFailures, code.code)
			knownFailuresErrs = append(knownFailuresErrs, err)
		}
		return resp, worked, err, i
	} else if isMapGet(code) {
		resp, worked, err, i = mapGetParse(code, index, codelines)
		if !worked {
			knownFailures = append(knownFailures, code.code)
			knownFailuresErrs = append(knownFailuresErrs, err)
		}
		return resp, worked, err, i
	} else if isIndexGet(code) {
		resp, worked, err, i = indexGetParse(code, index, codelines)
		if worked {
			return resp, worked, err, i
		}
	}
	if !worked {
		knownFailures = append(knownFailures, code.code)
		knownFailuresErrs = append(knownFailuresErrs, err)
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
