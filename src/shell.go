package main

import (
	"fmt"
	"os"
	"os/signal"
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
	for {
		indo := false
		totranslate := []UNPARSEcode{}
		code := input("\x1b[38;5;240m>>> \x1b[0m\x1b[1;5;240m")
		fmt.Print("\x1b[0m")
		if code == "" {
			continue
		}
		indo = true
		totranslate = append(totranslate, UNPARSEcode{code, code, 1, "<shell>"})
		for i := 2; indo; i++ {
			code := input("\x1b[38;5;240m... \x1b[0m\x1b[1;5;240m")
			fmt.Print("\x1b[0m")
			totranslate = append(totranslate, UNPARSEcode{code, code, i, "<shell>"})
			if code == "" {
				indo = false
			}
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
