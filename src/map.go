package main

import (
	"fmt"
	"strings"
)

var mapCompiled = makeRegex(`( *)\{(((( *).+( *):( *).+( *))|(` + spacelessVariable + `))(( *)\,(( *).+( *):( *).+( *))|(` + spacelessVariable + `)))*\}( *)`)

type createMap struct {
	body anymap
	code string
	line int
	path string
}

func isMap(code UNPARSEcode) bool {
	return mapCompiled.MatchString(code.code)
}

func parseMap(code UNPARSEcode) (any, UNPARSEcode) {
	trimmed := strings.Trim(code.code, " ")
	trimmed = trimmed[1 : len(trimmed)-1]
	fmt.Println(trimmed)
	return nil, UNPARSEcode{}
}

func Map(val anymap) ArObject {
	return ArObject{
		TYPE: "map",
		obj:  val,
	}
}

func Class(val anymap) ArObject {
	return ArObject{
		TYPE: "class",
		obj:  val,
	}
}
