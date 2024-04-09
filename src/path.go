package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var ArPath = Map(
	anymap{
		"ReadDir": builtinFunc{
			"ReadDir",
			func(args ...any) (any, ArErr) {
				if len(args) != 1 {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "ReadDir takes exactly 1 argument, got " + fmt.Sprint(len(args)),
						EXISTS:  true,
					}
				}
				args[0] = ArValidToAny(args[0])
				if typeof(args[0]) != "string" {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "ReadDir argument must be a string, got " + typeof(args[0]),
						EXISTS:  true,
					}
				}
				files, err := os.ReadDir(args[0].(string))
				if err != nil {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: err.Error(),
						EXISTS:  true,
					}
				}
				var ret []any
				for _, file := range files {
					ret = append(ret, file.Name())
				}
				return ret, ArErr{}
			}},
		"exists": builtinFunc{
			"exists",
			func(args ...any) (any, ArErr) {
				if len(args) != 1 {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "exists takes exactly 1 argument, got " + fmt.Sprint(len(args)),
						EXISTS:  true,
					}
				}
				args[0] = ArValidToAny(args[0])
				if typeof(args[0]) != "string" {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "exists argument must be a string, got " + typeof(args[0]),
						EXISTS:  true,
					}
				}
				_, err := os.Stat(args[0].(string))
				if err != nil {
					if os.IsNotExist(err) {
						return false, ArErr{}
					}
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: err.Error(),
						EXISTS:  true,
					}
				}
				return true, ArErr{}
			}},
		"mkAllDir": builtinFunc{
			"mkAllDir",
			func(args ...any) (any, ArErr) {
				if len(args) != 1 {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "mkAllDir takes exactly 1 argument, got " + fmt.Sprint(len(args)),
						EXISTS:  true,
					}
				}
				args[0] = ArValidToAny(args[0])
				if typeof(args[0]) != "string" {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "mkAllDir argument must be a string, got " + typeof(args[0]),
						EXISTS:  true,
					}
				}
				err := os.MkdirAll(args[0].(string), os.ModePerm)
				if err != nil {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: err.Error(),
						EXISTS:  true,
					}
				}
				return nil, ArErr{}
			}},
		"mkDir": builtinFunc{
			"mkDir",
			func(args ...any) (any, ArErr) {
				if len(args) != 1 {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "mkDir takes exactly 1 argument, got " + fmt.Sprint(len(args)),
						EXISTS:  true,
					}
				}
				args[0] = ArValidToAny(args[0])
				if typeof(args[0]) != "string" {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "mkDir argument must be a string, got " + typeof(args[0]),
						EXISTS:  true,
					}
				}
				err := os.Mkdir(args[0].(string), os.ModePerm)
				if err != nil {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: err.Error(),
						EXISTS:  true,
					}
				}
				return nil, ArErr{}
			}},
		"remove": builtinFunc{
			"remove",
			func(args ...any) (any, ArErr) {
				if len(args) != 1 {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "remove takes exactly 1 argument, got " + fmt.Sprint(len(args)),
						EXISTS:  true,
					}
				}
				args[0] = ArValidToAny(args[0])
				if typeof(args[0]) != "string" {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "remove argument must be a string, got " + typeof(args[0]),
						EXISTS:  true,
					}
				}
				err := os.Remove(args[0].(string))
				if err != nil {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: err.Error(),
						EXISTS:  true,
					}
				}
				return nil, ArErr{}
			}},
		"isDir": builtinFunc{
			"isDir",
			func(args ...any) (any, ArErr) {
				if len(args) != 1 {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "isDir takes exactly 1 argument, got " + fmt.Sprint(len(args)),
						EXISTS:  true,
					}
				}
				args[0] = ArValidToAny(args[0])
				if typeof(args[0]) != "string" {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "isDir argument must be a string, got " + typeof(args[0]),
						EXISTS:  true,
					}
				}
				stat, err := os.Stat(args[0].(string))
				if err != nil {
					if os.IsNotExist(err) {
						return false, ArErr{}
					}
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: err.Error(),
						EXISTS:  true,
					}
				}
				return stat.IsDir(), ArErr{}
			}},
		"join": builtinFunc{
			"join",
			func(args ...any) (any, ArErr) {
				if len(args) != 1 {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "join takes exactly 1 argument, got " + fmt.Sprint(len(args)),
						EXISTS:  true,
					}
				}
				args[0] = ArValidToAny(args[0])
				switch arr := args[0].(type) {
				case []any:
					var Path []string
					for _, x := range arr {
						x = ArValidToAny(x)
						if typeof(x) != "string" {
							return nil, ArErr{
								TYPE:    "Runtime Error",
								message: "join argument must be an array of strings, got " + typeof(x),
								EXISTS:  true,
							}
						}
						Path = append(Path, x.(string))
					}
					return filepath.Join(Path...), ArErr{}
				}
				return nil, ArErr{
					TYPE:    "Runtime Error",
					message: "join argument must be an array, got " + typeof(args[0]),
					EXISTS:  true,
				}
			}},
		"parent": builtinFunc{
			"parent",
			func(args ...any) (any, ArErr) {
				if len(args) != 1 {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "parent takes exactly 1 argument, got " + fmt.Sprint(len(args)),
						EXISTS:  true,
					}
				}
				args[0] = ArValidToAny(args[0])
				if typeof(args[0]) != "string" {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "parent argument must be a string, got " + typeof(args[0]),
						EXISTS:  true,
					}
				}
				return path.Dir(filepath.ToSlash(args[0].(string))), ArErr{}
			},
		},
	})
