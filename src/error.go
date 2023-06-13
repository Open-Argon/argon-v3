package main

import (
	"fmt"

	"github.com/jwalton/go-supportscolor"
)

type ArErr struct {
	TYPE    string
	message string
	line    int
	path    string
	code    string
	EXISTS  bool
}

func panicErr(err ArErr) {
	if err.code != "" && err.line != 0 && err.path != "" {
		fmt.Println("  File:", err.path+":"+fmt.Sprint(err.line))
		fmt.Println("    " + err.code)
		fmt.Println()
	}
	if supportscolor.Stdout().SupportsColor {
		fmt.Printf("\x1b[%dm%s\x1b[0m", 91, fmt.Sprint(err.TYPE, ": ", err.message, "\n"))
	} else {
		fmt.Println(fmt.Sprint(err.TYPE, ": ", err.message))
	}
}
