package main

type Callable struct {
	name   string
	params []string
	run    any
	code   string
	stack  stack
	line   int
}
