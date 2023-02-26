package main

type UNPARSEcode struct {
	code     string
	realcode string
	line     int
	path     string
}

// returns (translateNumber | translateString| nil), success, error
func translateVal(code UNPARSEcode, index int, codelines []UNPARSEcode, isLine bool) (any, bool, ArErr, int) {
	if isLine {
		if isBlank(code) {
			return nil, true, ArErr{}, 1
		} else if isComment(code) {
			resp, worked, err := parseComment(code, index, codelines)
			if worked {
				return resp, worked, err, 1
			}
		}
	}
	if isCall(code) {
		return parseCall(code, index, codelines)
	} else if isVariable(code) {
		return parseVariable(code)
	} else if isMapGet(code) {
		return mapGetParse(code, index, codelines)
	} else if isNumber(code) {
		return parseNumber(code)
	} else if isString(code) {
		return parseString(code)
	}
	return nil, false, ArErr{"Syntax Error", "invalid syntax", code.line, code.path, code.realcode, true}, 1
}

// returns [](translateNumber | translateString), error
func translate(codelines []UNPARSEcode) ([]any, ArErr) {
	translated := []any{}
	for i, code := range codelines {
		val, _, err, _ := translateVal(code, i, codelines, true)

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
