package main

import "fmt"

func runLine(line any) (any, string) {
	switch line.(type) {
	case translateNumber:
		return (numberToString(line.(translateNumber).number, 0)), ""
	case translateString:
		return (line.(translateString).str), ""
	}
	return nil, "Error: invalid code on line " + fmt.Sprint(line.(translateNumber).line) + ": " + line.(translateNumber).code
}

// returns error
func run(translated []any) (any, string) {
	for _, val := range translated {
		_, err := runLine(val)
		if err != "" {
			return nil, err
		}
	}
	return nil, ""
}
