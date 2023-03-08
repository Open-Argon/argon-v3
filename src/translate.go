package main

type UNPARSEcode struct {
	code     string
	realcode string
	line     int
	path     string
}

// returns (number | string | nil), success, error, step
func translateVal(code UNPARSEcode, index int, codelines []UNPARSEcode, isLine bool) (any, bool, ArErr, int) {

	if isLine {
		if isBlank(code) {
			return nil, true, ArErr{}, 1
		} else if isDeleteVariable(code) {
			return parseDelete(code, index, codelines)
		} else if isComment(code) {
			resp, worked, err, step := parseComment(code, index, codelines)
			if worked {
				return resp, worked, err, step
			}
		}
	}

	if isBrackets(code) {
		bracket, worked, err, step := parseBrackets(code, index, codelines)
		if worked {
			return bracket, worked, err, step
		}
	}
	if isSetVariable(code) {
		setvar, worked, err, step := parseSetVariable(code, index, codelines)
		if worked {
			return setvar, worked, err, step
		}
	}
	if isAutoAsignVariable(code) {
		setvar, worked, err, step := parseAutoAsignVariable(code, index, codelines)
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
	} else if isCall(code) {
		call, worked, err, step := parseCall(code, index, codelines)
		if worked {
			return call, worked, err, step
		}
	}
	if isVariable(code) {
		return parseVariable(code)
	} else if isMapGet(code) {
		return mapGetParse(code, index, codelines)
	} else if isString(code) {
		return parseString(code)
	}
	return nil, false, ArErr{"Syntax Error", "invalid syntax", code.line, code.path, code.realcode, true}, 1
}

// returns [](number | string), error
func translate(codelines []UNPARSEcode) ([]any, ArErr) {
	translated := []any{}
	for i := 0; i < len(codelines); {
		val, _, err, step := translateVal(codelines[i], i, codelines, true)
		i += step
		if err.EXISTS {
			return nil, err
		}
		if val == nil {
			continue
		}
		translated = append(translated, val)
	}
	return translated, ArErr{}
}
