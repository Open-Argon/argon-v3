package main

import (
	"fmt"
	"os"
)

// args without the program path
var Args = os.Args[1:]

type stack = []ArObject

func newscope() ArObject {
	return Map(anymap{})
}

func main() {
	debugInit()

	if !debug {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("There was a fundamental error in argon v3 that caused it to crash.")
				fmt.Println()
				fmt.Println("website:", website)
				fmt.Println("docs:", docs)
				fmt.Println()
				if fork {
					fmt.Println("This is a fork of Open-Argon. Please report this to the fork's maintainer.")
					fmt.Println("Fork repo:", forkrepo)
					fmt.Println("Fork issue page:", forkissuesPage)
					fmt.Println()
				} else {
					fmt.Println("Please report this to the Open-Argon team.")
					fmt.Println("Main repo:", mainrepo)
					fmt.Println("Issue page:", mainissuesPage)
					fmt.Println()
				}
				fmt.Println("please include the following information:")
				fmt.Println("panic:", r)
				os.Exit(1)
			}
		}()
	}
	initRandom()
	garbageCollect()
	global := makeGlobal()
	if len(Args) == 0 {
		shell(global)
		os.Exit(0)
	}
	ex, e := os.Getwd()
	if e != nil {
		panic(e)
	}
	_, err := importMod(Args[0], ex, true, global)
	if err.EXISTS {
		panicErr(err)
		os.Exit(1)
	}
}
