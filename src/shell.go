package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
)

func shell(global ArObject) {
	stack := stack{global, newscope()}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				fmt.Println("\x1b[0m\n\x1b[32;5;240mBye :)\x1b[0m")
				os.Exit(0)
			}
		}
	}()
	fmt.Print("\x1b[32;240mWelcome to the Argon v3!\x1b[0m\n\n")
	for {
		indent := 0
		previous := 0
		totranslate := []UNPARSEcode{}
		textBefore := ">>>"
		for i := 1; indent > 0 || (previous != indent && indent >= 0) || i == 1; i++ {
			indentStr := strings.Repeat("    ", indent)
			inp, err := input("\x1b[38;240m" + textBefore + indentStr + " \x1b[0m\x1b[1;240m")
			if err != nil {
				fmt.Println("\x1b[0m\n\x1b[32;240mBye :)\x1b[0m")
				os.Exit(0)
			}
			code := indentStr + inp
			fmt.Print("\x1b[0m")
			totranslate = append(totranslate, UNPARSEcode{code, code, i, "<shell>"})
			trimmed := strings.TrimSpace(code)
			previous = indent
			if len(trimmed) >= 2 && trimmed[len(trimmed)-2:] == "do" {
				indent++
			} else if trimmed == "" {
				indent--
			}
			textBefore = "..."
		}
		translated, translationerr := translate(totranslate)
		count := len(translated)
		if translationerr.EXISTS {
			panicErr(translationerr)
		}
		output, runimeErr := ThrowOnNonLoop(run(translated, stack))
		output = openReturn(output)

		if runimeErr.EXISTS {
			panicErr(runimeErr)
		} else if count == 1 {
			fmt.Println(anyToArgon(output, true, true, 3, 0, true, 1))
		}
	}
}
