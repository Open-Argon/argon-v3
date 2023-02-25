package main

import (
	"fmt"
	"strconv"
	"strings"
)

var stringCompile = makeRegex("(( *)\"((\\\\([a-z\\\"'`]))|[^\\\"])*\"( *))|(( *)'((\\\\([a-z\\'\"`]))|[^\\'])*'( *))")

func isString(code UNPARSEcode) bool {
	return stringCompile.MatchString(code.code)
}

func unquoted(
	str string,
) (string, error) {
	str = strings.Trim(str, " ")
	if str[0] == '\'' {
		str = strings.Replace(str, "\\\"", "\"", -1)
		str = strings.Replace(str, "\"", "\\\"", -1)
	}
	str = str[1 : len(str)-1]
	str = strings.Replace(str, "\\'", "'", -1)
	str = "\"" + str + "\""
	return strconv.Unquote(str)
}

// returns translateString, success, error
func parseString(code UNPARSEcode) (translateString, bool, string) {
	trim := strings.Trim(code.code, " ")

	unquoted, err := unquoted(trim)
	if err != nil {
		return translateString{}, false, "Syntax Error: invalid string on line " + fmt.Sprint(code.line) + ": " + code.code
	}

	return translateString{
		str:  unquoted,
		line: code.line,
	}, true, ""
}
