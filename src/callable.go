package main

type Callable struct {
	name   string
	params []string
	code   []any
	stack  stack
	line   int
}
