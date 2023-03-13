package main

import (
	"math/big"
	"strings"
)

var squarerootcompiled = makeRegex(`(.|\n)*√(.|\n)+`)

type squareroot struct {
	first  any
	second any
	code   string
	line   int
	path   string
}

func issquareroot(code UNPARSEcode) bool {
	return squarerootcompiled.MatchString(code.code)
}

func parseSquareroot(code UNPARSEcode, index int, lines []UNPARSEcode) (squareroot, bool, ArErr, int) {
	split := strings.SplitN(code.code, "√", 2)
	first := strings.TrimSpace(split[0])
	second := strings.TrimSpace(split[1])
	var firstparsed any = newNumber().SetInt64(1)
	outputlen := 0
	if first != "" {
		val, worked, err, i := translateVal(UNPARSEcode{code: first, realcode: code.realcode, line: code.line, path: code.path}, index, lines, 0)
		if !worked {
			return squareroot{}, false, err, i
		}
		outputlen += i - 1
		firstparsed = val
	}
	secondparsed, worked, err, i := translateVal(UNPARSEcode{code: second, realcode: code.realcode, line: code.line, path: code.path}, index, lines, 0)
	if !worked {
		return squareroot{}, false, err, i
	}
	outputlen += i
	return squareroot{
		firstparsed,
		secondparsed,
		code.realcode,
		code.line,
		code.path,
	}, true, ArErr{}, outputlen
}

func runSquareroot(squareroot squareroot, stack stack, stacklevel int) (number, ArErr) {
	val1, err := runVal(squareroot.first, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	val2, err := runVal(squareroot.second, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	if typeof(val1) != "number" || typeof(val2) != "number" {
		return nil, ArErr{"Type Error", "Cannot take the square root of a non-number", squareroot.line, squareroot.path, squareroot.code, true}
	}

	var x big.Float
	x.SetPrec(30)
	x.SetRat(val2.(number))

	var s big.Float
	s.SetPrec(15)
	s.Sqrt(&x)

	r, _ := s.Rat(nil)
	r.Mul(r, val1.(number))
	return r, ArErr{}
}
