package main

import (
	"fmt"
	"os"
	"syscall/js"

	"golang.org/x/term"
)

func input(args ...any) string {
	output := []any{}
	for i := 0; i < len(args); i++ {
		output = append(output, anyToArgon(args[i], false, true, 3, 0, true, 0))
		if i != len(args)-1 {
			output = append(output, " ")
		}
	}
	text := fmt.Sprint(output...)
	fmt.Print(text)
	inputData := js.Global().Get("getInput").Invoke(fmt.Sprint(output...)).String()
	fmt.Println(inputData)
	return (inputData)
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
			return "", fmt.Errorf("keyboard interupt")
		} else if char[0] == '\r' || char[0] == '\n' {
			fmt.Println()
			break
		} else if char[0] == '\b' || char[0] == 127 {
			if len(password) > 0 {
				password = password[:len(password)-1]
				fmt.Print("\b \b")
			}
		} else {
			password = append(password, char[0])
			fmt.Print("*")
		}
	}
	fmt.Print("\r")
	return string(password), nil
}
