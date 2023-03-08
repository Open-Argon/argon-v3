package main

import (
	"strings"
)

func getValuesFromCommas(str string, index int, codelines []UNPARSEcode) ([]any, bool, ArErr) {
	// make a function which takes a string of code and returns a translated values
	str = strings.Trim(str, " ")
	commasplit := strings.Split(str, ",")
	temp := []string{}
	arguments := []any{}
	if str != "" {
		for i, arg := range commasplit {
			temp = append(temp, arg)
			test := strings.TrimSpace(strings.Join(temp, ","))
			resp, worked, _, _ := translateVal(UNPARSEcode{code: test, realcode: codelines[index].realcode, line: index + 1, path: codelines[index].path}, index, codelines, false)
			if worked {
				arguments = append(arguments, resp)
				temp = []string{}
			} else if i == len(commasplit)-1 {
				return nil, false, ArErr{"Syntax Error", "invalid argument", codelines[index].line, codelines[index].path, codelines[index].realcode, true}
			}
		}
	}
	return arguments, true, ArErr{}
}