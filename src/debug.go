package main

import (
	"fmt"
	"os"
)

var debug = os.Getenv("__ARGON_DEBUG__") == "true"

func debugPrintln(a ...interface{}) {
	if debug {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("debugPrintln: panic:", r)
				}
			}()
			fmt.Println(a...)
		}()
	}
}
