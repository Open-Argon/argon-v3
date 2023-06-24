package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

var ArFile = Map(anymap{
	"read":  builtinFunc{"read", ArRead},
	"write": builtinFunc{"write", ArWrite},
})

func readtext(file *os.File) (string, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func readbinary(file *os.File) ([]byte, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ArRead(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "read takes 1 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
	}
	if typeof(args[0]) != "string" {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "read takes a string not type '" + typeof(args[0]) + "'", EXISTS: true}
	}
	args[0] = ArValidToAny(args[0])
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
			return text, ArErr{}
		}},
		"json": builtinFunc{"json", func(...any) (any, ArErr) {
			text, err := readtext(file)
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return jsonparse(text), ArErr{}
		}},
		"contentType": builtinFunc{"contentType", func(...any) (any, ArErr) {
			mimetype, err := mimetype.DetectFile(filename)
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return mimetype.String(), ArErr{}
		}},
		"buffer": builtinFunc{"buffer", func(args ...any) (any, ArErr) {
			if len(args) > 1 {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "buffer takes 0 or 1 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			if len(args) == 1 {
				if typeof(args[0]) != "number" {
					return ArObject{}, ArErr{TYPE: "Runtime Error", message: "buffer takes a number not type '" + typeof(args[0]) + "'", EXISTS: true}
				}
				size := args[0].(number)
				if size.Denom().Int64() != 1 {
					return ArObject{}, ArErr{TYPE: "Runtime Error", message: "buffer takes an integer not type '" + typeof(args[0]) + "'", EXISTS: true}
				}
				buf := make([]byte, size.Num().Int64())
				n, err := file.Read(buf)
				if err != nil {
					return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
				}
				return ArBuffer(buf[:n]), ArErr{}
			}
			bytes, err := readbinary(file)
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return ArBuffer(bytes), ArErr{}
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
	args[0] = ArValidToAny(args[0])
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
		"buffer": builtinFunc{"buffer", func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "buffer takes 1 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			if typeof(args[0]) != "buffer" {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "buffer takes a buffer not type '" + typeof(args[0]) + "'", EXISTS: true}
			}
			args[0] = ArValidToAny(args[0])
			file.Write(args[0].([]byte))
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
