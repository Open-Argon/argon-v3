package main

type Callable struct {
	params []string
	run    any
	code   string
	stack  stack
	line   int
}
