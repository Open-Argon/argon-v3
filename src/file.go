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
	"move":  builtinFunc{"move", ARmoveFile},
	"copy":  builtinFunc{"copy", ARcopyFile},
})

func ARmoveFile(args ...any) (any, ArErr) {
	if len(args) != 2 {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "move takes 2 arguments, got " + fmt.Sprint(len(args)), EXISTS: true}
	}
	if typeof(args[0]) != "string" {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "move takes a string not type '" + typeof(args[0]) + "'", EXISTS: true}
	}
	if typeof(args[1]) != "string" {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "move takes a string not type '" + typeof(args[1]) + "'", EXISTS: true}
	}
	args[0] = ArValidToAny(args[0])
	args[1] = ArValidToAny(args[1])
	err := os.Rename(args[0].(string), args[1].(string))
	if err != nil {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
	}
	return nil, ArErr{}
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func ARcopyFile(args ...any) (any, ArErr) {
	if len(args) != 2 {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "copy takes 2 arguments, got " + fmt.Sprint(len(args)), EXISTS: true}
	}
	if typeof(args[0]) != "string" {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "copy takes a string not type '" + typeof(args[0]) + "'", EXISTS: true}
	}
	if typeof(args[1]) != "string" {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "copy takes a string not type '" + typeof(args[1]) + "'", EXISTS: true}
	}
	args[0] = ArValidToAny(args[0])
	args[1] = ArValidToAny(args[1])
	_, err := copyFile(args[0].(string), args[1].(string))
	if err != nil {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
	}
	return nil, ArErr{}
}

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

	fileInfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return ArObject{}, ArErr{TYPE: "Runtime Error", message: "File does not exist: " + filename, EXISTS: true}
		} else {
			return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
		}
	} else if fileInfo.Mode().IsDir() {
		return ArObject{}, ArErr{TYPE: "Runtime Error", message: "path goes to a directory, not a file: " + filename, EXISTS: true}
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
			return jsonparse(text)
		}},
		"contentType": builtinFunc{"contentType", func(...any) (any, ArErr) {
			file.Seek(0, io.SeekStart)
			mimetype, err := mimetype.DetectReader(file)
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
		"seek": builtinFunc{"seek", func(args ...any) (any, ArErr) {
			if len(args) != 1 {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "seek takes 1 argument, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			if typeof(args[0]) != "number" {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "seek takes a number not type '" + typeof(args[0]) + "'", EXISTS: true}
			}
			offset := args[0].(number)
			if offset.Denom().Int64() != 1 {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: "seek takes an integer not type '" + typeof(args[0]) + "'", EXISTS: true}
			}
			_, err := file.Seek(offset.Num().Int64(), io.SeekStart)
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return nil, ArErr{}
		}},
		"size": builtinFunc{"size", func(...any) (any, ArErr) {
			info, err := file.Stat()
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return newNumber().SetInt64(info.Size()), ArErr{}
		}},
		"ModTime": builtinFunc{"ModTime", func(...any) (any, ArErr) {
			info, err := file.Stat()
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return ArTimeClass(info.ModTime()), ArErr{}
		}},
		"close": builtinFunc{"close", func(...any) (any, ArErr) {
			err := file.Close()
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return nil, ArErr{}
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
			args[0] = ArValidToAny(args[0])
			_, err = file.Write([]byte(args[0].(string)))
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
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
			_, err = file.Write(args[0].([]byte))
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
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
			_, err = file.Write([]byte(jsonstr))
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return nil, ArErr{}
		}},
		"close": builtinFunc{"close", func(...any) (any, ArErr) {
			err := file.Close()
			if err != nil {
				return ArObject{}, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return nil, ArErr{}
		}},
	}), ArErr{}

}
