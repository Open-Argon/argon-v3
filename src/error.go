package main

import (
	"fmt"
	"os"
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
	fmt.Println("  File:", err.path+":"+fmt.Sprint(err.line))
	fmt.Println("    " + err.code)
	fmt.Println()
	fmt.Println(err.TYPE+":", err.message)
	os.Exit(1)
}
