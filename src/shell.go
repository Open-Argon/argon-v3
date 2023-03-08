package main

import (
	"fmt"
	"os"
	"os/signal"
)

func shell() {
	global := stack{vars, scope{}}
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
		code := input("\x1b[38;5;240m>>> \x1b[0m\x1b[1;5;240m")
		fmt.Print("\x1b[0m")
		translated, translationerr := translate([]UNPARSEcode{{code, code, 1, "<shell>"}})
		if translationerr.EXISTS {
			panicErr(translationerr)
		}
		_, runimeErr, count, output := run(translated, global)
		if runimeErr.EXISTS {
			panicErr(runimeErr)
		}
		if count == 0 {
			fmt.Println(anyToArgon(output, true, true, 3, 0, true, 1))
		}
	}
}
