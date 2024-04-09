package main

import (
	"fmt"
)

type ArErr struct {
	TYPE    string
	message string
	line    int
	path    string
	code    string
	EXISTS  bool
}

func ArThrowError(a ...any) (any, ArErr) {
	if len(a) != 2 {
		return nil, ArErr{
			TYPE:    "Type Error",
			message: "throwError takes 2 arguments, " + fmt.Sprint(len(a)) + " given",
			EXISTS:  true,
		}
	} else if typeof(a[0]) != "string" {
		return nil, ArErr{
			TYPE:    "Type Error",
			message: "throwError type must be a string",
			EXISTS:  true,
		}
	} else if typeof(a[1]) != "string" {
		return nil, ArErr{
			TYPE:    "Type Error",
			message: "throwError message must be a string",
			EXISTS:  true,
		}
	}
	return nil, ArErr{
		TYPE:    ArValidToAny(a[0]).(string),
		message: ArValidToAny(a[1]).(string),
		EXISTS:  true,
	}
}

func panicErr(err ArErr) {
	if err.code != "" && err.line != 0 && err.path != "" {
		fmt.Println("  File:", err.path+":"+fmt.Sprint(err.line))
		fmt.Println("    " + err.code)
		fmt.Println()
	}
	fmt.Printf("\x1b[%dm%s\x1b[0m", 91, fmt.Sprint(err.TYPE, ": ", err.message, "\n"))
}
