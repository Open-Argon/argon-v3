package main

import (
	"bytes"
	"io"
	"os"
)

var ArFile = ArMap{
	"read": builtinFunc{"read", ArRead},
}

func readtext(file *os.File) (string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return "", err
	}
	return string(buf.Bytes()), nil
}

func ArRead(args ...any) (any, ArErr) {
	if len(args) == 0 {
		return ArMap{}, ArErr{TYPE: "Runtime Error", message: "open takes 1 argument", EXISTS: true}
	}
	if typeof(args[0]) != "string" {
		return ArMap{}, ArErr{TYPE: "Runtime Error", message: "open takes a string not a '" + typeof(args[0]) + "'", EXISTS: true}
	}
	filename := args[0].(string)
	file, err := os.Open(filename)
	if err != nil {
		return ArMap{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
	}
	return ArMap{
		"text": builtinFunc{"text", func(...any) (any, ArErr) {
			text, err := readtext(file)
			if err != nil {
				return ArMap{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return text, ArErr{}
		}},
		"json": builtinFunc{"json", func(...any) (any, ArErr) {
			text, err := readtext(file)
			if err != nil {
				return ArMap{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return parse(text), ArErr{}
		}},
		"line": builtinFunc{"line", func(...any) (any, ArErr) {

			return "", ArErr{}
		}},
	}, ArErr{}
}
