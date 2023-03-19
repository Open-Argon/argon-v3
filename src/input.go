package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

func input(args ...any) string {
	output := []any{}
	for i := 0; i < len(args); i++ {
		output = append(output, anyToArgon(args[i], false, true, 3, 0, true, 0))
	}
	fmt.Print(output...)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	return input
}

func getPassword(args ...any) (string, error) {
	output := []any{}
	for i := 0; i < len(args); i++ {
		output = append(output, anyToArgon(args[i], false, true, 3, 0, true, 0))
	}
	fmt.Print(output...)
	password := []byte{}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	for {
		char := make([]byte, 1)
		_, err := os.Stdin.Read(char)
		if err != nil {
			fmt.Println(err)
			break
		}
		if char[0] == 3 || char[0] == 4 {
			return "", fmt.Errorf("User cancelled")
		} else if char[0] == '\r' || char[0] == '\n' {
			fmt.Println()
			break
		} else if char[0] == '\b' {
			if len(password) > 0 {
				password = password[:len(password)-1]
				fmt.Print("\b \b")
			}
		} else {
			password = append(password, char[0])
			fmt.Print("*")
		}
	}
	return string(password), nil
}
