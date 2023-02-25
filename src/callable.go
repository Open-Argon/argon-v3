package main

type Callable struct {
	name   string
	params []string
	code   []any
	stack  []map[string]variableValue
	line   int
}
