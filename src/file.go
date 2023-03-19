package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func ArOpen(args ...any) (any, ArErr) {
	if len(args) > 2 {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "open takes 1 or 2 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
	}
	if typeof(args[0]) != "string" {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "open takes a string not type '" + typeof(args[0]) + "'", EXISTS: true}
	}
	path := args[0].(string)
	mode := "r"
	if len(args) == 2 {
		if typeof(args[1]) != "string" {
			return ArObject{}, ArErr{TYPE: "Runtime Error", message: "open takes a string not type '" + typeof(args[1]) + "'", EXISTS: true}
		}
		mode = args[1].(string)
	}
	if mode != "r" && mode != "w" {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "open mode must be 'r', or 'w'", EXISTS: true}
	}
	if mode == "r" {
		file, err := os.Open(path)
		if err != nil {
			return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
		}
		return Map(anymap{
			"text": builtinFunc{"text", func(...any) (any, ArErr) {
				text, err := readtext(file)
				if err != nil {
					return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
				}
				return ArString(text), ArErr{}
			},
			},
			"json": builtinFunc{"json", func(...any) (any, ArErr) {
				text, err := readtext(file)
				if err != nil {
					return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
				}
				return jsonparse(text), ArErr{}
			},
			},
		}), ArErr{}
	}
	file, err := os.Create(path)
	if err != nil {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
	}
	return Map(anymap{
		"text": builtinFunc{"text", func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "text takes 1 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			if typeof(args[0]) != "string" {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "text takes a string not type '" + typeof(args[0]) + "'", EXISTS: true}
			}
			file.Write([]byte(args[0].(string)))
			return nil, ArErr{}
		}},
		"json": builtinFunc{"json", func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "json takes 1 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			jsonstr, err := jsonstringify(args[0], 0)
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			file.Write([]byte(jsonstr))
			return nil, ArErr{}
		}},
		"append": builtinFunc{"append", func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "append takes 1 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			if typeof(args[0]) != "string" {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "append takes a string not type '" + typeof(args[0]) + "'", EXISTS: true}
			}
			file.Write([]byte(args[0].(string)))
			return nil, ArErr{}
		}},
	}), ArErr{}

}

func readtext(file *os.File) (string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ArRead(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "read takes 1 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
	}
	if typeof(args[0]) != "string" {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "read takes a string not type '" + typeof(args[0]) + "'", EXISTS: true}
	}
	filename := args[0].(string)
	file, err := os.Open(filename)
	if err != nil {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
	}
	return Map(anymap{
		"text": builtinFunc{"text", func(...any) (any, ArErr) {
			text, err := readtext(file)
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return ArString(text), ArErr{}
		}},
		"json": builtinFunc{"json", func(...any) (any, ArErr) {
			text, err := readtext(file)
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return jsonparse(text), ArErr{}
		}},
	}), ArErr{}
}

func ArWrite(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "write takes 1 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
	}
	if typeof(args[0]) != "string" {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "write takes a string not type '" + typeof(args[0]) + "'", EXISTS: true}
	}
	filename := args[0].(string)
	file, err := os.Create(filename)
	if err != nil {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
	}
	return Map(anymap{
		"text": builtinFunc{"text", func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "text takes 1 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			if typeof(args[0]) != "string" {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "text takes a string not type '" + typeof(args[0]) + "'", EXISTS: true}
			}
			file.Write([]byte(args[0].(string)))
			return nil, ArErr{}
		}},
		"json": builtinFunc{"json", func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "json takes 1 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			jsonstr, err := jsonstringify(args[0], 0)
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			file.Write([]byte(jsonstr))
			return nil, ArErr{}
		}},
	}), ArErr{}

}
