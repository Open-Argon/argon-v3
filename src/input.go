package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chzyer/readline"
	"golang.org/x/term"
)

var tempFilePath = os.TempDir() + "/argon_input_history.tmp"

func input(args ...any) (string, error) {
	output := []any{}
	for i := 0; i < len(args); i++ {
		output = append(output, anyToArgon(args[i], false, true, 3, 0, true, 0))
	}
	message := fmt.Sprint(output...)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:            message,
		HistoryFile:       tempFilePath,
		HistorySearchFold: true,
	})
	if err != nil {
		log.Fatalf("Failed to create readline instance: %v", err)
	}
	defer rl.Close()
	line, err := rl.Readline()
	if err != nil { // io.EOF or other error
		return "", err
	}
	return line, nil
}

func getPassword(args ...any) (string, error) {
	output := []any{}
	for i := 0; i < len(args); i++ {
		output = append(output, anyToArgon(args[i], false, true, 3, 0, true, 0))
	}
	message := fmt.Sprint(output...)
	rl, err := readline.NewEx(&readline.Config{
		Prompt:     message,
		MaskRune:   '*',
		EnableMask: true,
	})
	if err != nil {
		log.Fatalf("Failed to create readline instance: %v", err)
	}
	defer rl.Close()
	line, err := rl.Readline()
	if err != nil { // io.EOF or other error
		return "", err
	}
	return line, nil
}

func pause() {
	fmt.Print("Press Enter to continue...")
	term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
}
