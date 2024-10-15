package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

var tempFilePath = os.TempDir() + "/argon_input_history.tmp"

func pause() {
	fmt.Print("Press Enter to continue...")
	term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
}
