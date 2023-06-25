package main

import (
	"fmt"
	"os"
	"path"
)

var ArPath = Map(
	anymap{
		"ReadDir": builtinFunc{
			"ReadDir",
			func(args ...any) (any, ArErr) {
				if len(args) != 1 {
					return nil, ArErr{
						TYPE:    "runtime",
						message: "ReadDir takes exactly 1 argument, got " + fmt.Sprint(len(args)),
						EXISTS:  true,
					}
				}
				args[0] = ArValidToAny(args[0])
				if typeof(args[0]) != "string" {
					return nil, ArErr{
						TYPE:    "runtime",
						message: "ReadDir argument must be a string, got " + typeof(args[0]),
						EXISTS:  true,
					}
				}
				files, err := os.ReadDir(args[0].(string))
				if err != nil {
					return nil, ArErr{
						TYPE:    "runtime",
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
		"join": builtinFunc{
			"join",
			func(args ...any) (any, ArErr) {
				if len(args) != 1 {
					return nil, ArErr{
						TYPE:    "runtime",
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
								TYPE:    "runtime",
								message: "join argument must be an array of strings, got " + typeof(x),
								EXISTS:  true,
							}
						}
						Path = append(Path, x.(string))
					}
					return path.Join(Path...), ArErr{}
				}
				return nil, ArErr{
					TYPE:    "runtime",
					message: "join argument must be an array, got " + typeof(args[0]),
					EXISTS:  true,
				}
			}},
	})
