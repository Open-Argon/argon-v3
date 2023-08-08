package main

import (
	"strings"
)

var ifstatmentCompile = makeRegex(`( *)if( )+\((.|\n)+\)( )+(.|\n)+`)
var elseifstatmentCompile = makeRegex(`( *)else( )+if( )+\((.|\n)+\)( )+(.|\n)+`)
var elseCompile = makeRegex(`( *)else( )+(.|\n)+`)

type statement struct {
	condition any
	THEN      any
	line      int
	code      string
	path      string
}

type ifstatement struct {
	conditions []statement
	ELSE       any
	line       int
	code       string
	path       string
}

func isIfStatement(code UNPARSEcode) bool {
	return ifstatmentCompile.MatchString(code.code)
}

func parseIfStatement(code UNPARSEcode, index int, codeline []UNPARSEcode) (ifstatement, bool, ArErr, int) {
	conditions := []statement{}
	var ELSE any
	i := index
	for i < len(codeline) && (elseifstatmentCompile.MatchString(codeline[i].code) || i == index) {
		trimmed := strings.TrimSpace(codeline[i].code)
		trimmed = strings.TrimSpace(trimmed[strings.Index(trimmed, "("):])
		trimmed = (trimmed[1:])
		split := strings.Split(trimmed, ")")
		for j := len(split) - 1; j > 0; j-- {
			conditionjoined := strings.Join(split[:j], ")")
			thenjoined := strings.Join(split[j:], ")")
			outindex := 0
			conditionval, worked, err, step := translateVal(
				UNPARSEcode{
					code:     conditionjoined,
					realcode: codeline[i].realcode,
					line:     code.line,
					path:     code.path,
				},
				i,
				codeline,
				0)
			if err.EXISTS || !worked {
				if j == 1 {
					return ifstatement{}, worked, err, step
				} else {
					continue
				}
			}

			outindex += step
			thenval, worked, err, step := translateVal(
				UNPARSEcode{
					code:     thenjoined,
					realcode: codeline[i].realcode,
					line:     code.line,
					path:     code.path,
				},
				i,
				codeline,
				2,
			)
			if err.EXISTS || !worked {
				return ifstatement{}, worked, err, step
			}
			outindex += step - 1
			conditions = append(conditions, statement{
				condition: conditionval,
				THEN:      thenval,
				line:      code.line,
				code:      code.realcode,
				path:      code.path,
			})
			i += outindex
			break
		}
	}
	if i < len(codeline) && elseCompile.MatchString(codeline[i].code) {
		trimmed := strings.TrimSpace(codeline[i].code)
		trimmed = strings.TrimSpace(trimmed[4:])
		ELSEval, _, err, step := translateVal(
			UNPARSEcode{
				code:     trimmed,
				realcode: codeline[i].realcode,
				line:     code.line,
				path:     code.path,
			},
			i,
			codeline,
			2,
		)
		if err.EXISTS {
			return ifstatement{}, false, err, step
		}
		ELSE = ELSEval
		i += step
	}
	return ifstatement{
		conditions: conditions,
		ELSE:       ELSE,
		line:       code.line,
		code:       code.realcode,
		path:       code.path,
	}, true, ArErr{}, i - index
}

func runIfStatement(code ifstatement, stack stack, stacklevel int) (any, ArErr) {
	for _, condition := range code.conditions {
		newstack := append(stack, newscope())
		resp, err := runVal(condition.condition, newstack, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		if anyToBool(resp) {
			return runVal(condition.THEN, newstack, stacklevel+1)
		}
	}
	if code.ELSE != nil {
		return runVal(code.ELSE, append(stack, newscope()), stacklevel+1)
	}
	return nil, ArErr{}
}
