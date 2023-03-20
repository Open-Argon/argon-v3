package main

import (
	"strings"
)

var forloopCompile = makeRegex(`( *)for( +)\(( *)` + spacelessVariable + `( +)from( +)(\n|.)+to( +)(\n|.)+(( +)step( +)(\n|.)+)?( *)\)( *)(\n|.)+`)

type forLoop struct {
	variable string
	from     any
	to       any
	step     any
	body     any
	line     int
	code     string
	path     string
}

func isForLoop(code UNPARSEcode) bool {
	return forloopCompile.MatchString(code.code)
}

func parseForLoop(code UNPARSEcode, index int, codelines []UNPARSEcode) (forLoop, bool, ArErr, int) {
	totalstep := 0
	trimmed := strings.TrimSpace(strings.TrimSpace(code.code)[3:])[1:]
	split := strings.SplitN(trimmed, " from ", 2)
	name := strings.TrimSpace(split[0])
	split = strings.SplitN(split[1], " to ", 2)
	from := strings.TrimSpace(split[0])
	fromval, worked, err, fromstep := translateVal(UNPARSEcode{code: from, realcode: code.realcode, line: index + 1, path: code.path}, index, codelines, 0)
	if !worked {
		return forLoop{}, worked, err, fromstep
	}
	totalstep += fromstep
	tosplit := strings.Split(split[1], ")")
	for i := len(tosplit) - 1; i >= 0; i-- {
		innertotalstep := 0
		val := strings.Join(tosplit[:i], ")")
		valsplit := strings.SplitN(val, " step ", 2)
		var stepval any
		if len(valsplit) == 2 {
			step := strings.TrimSpace(valsplit[1])
			stepval_, worked, err, stepstep := translateVal(UNPARSEcode{code: step, realcode: code.realcode, line: code.line, path: code.path}, index, codelines, 0)
			if !worked {
				if i == 0 {
					return forLoop{}, worked, err, stepstep
				}
				continue
			}
			innertotalstep += stepstep - 1
			stepval = stepval_
		} else {
			stepval = newNumber().SetInt64(1)
		}
		to := strings.TrimSpace(valsplit[0])
		toval, worked, err, tostep := translateVal(UNPARSEcode{code: to, realcode: code.realcode, line: code.line, path: code.path}, index, codelines, 0)
		if !worked {
			if i == 0 {
				return forLoop{}, worked, err, tostep
			}
			continue
		}
		innertotalstep += tostep - 1
		body := strings.Join(tosplit[i:], ")")
		bodyval, worked, err, bodystep := translateVal(UNPARSEcode{code: body, realcode: code.realcode, line: code.line, path: code.path}, index, codelines, 1)
		if !worked {
			if i == 0 {
				return forLoop{}, worked, err, bodystep
			}
			continue
		}
		innertotalstep += bodystep - 1
		return forLoop{variable: name, from: fromval, to: toval, step: stepval, body: bodyval, line: code.line, code: code.code, path: code.path}, true, ArErr{}, totalstep + innertotalstep
	}
	return forLoop{}, false, ArErr{"Syntax Error", "invalid for loop", code.line, code.path, code.realcode, true}, 1
}
func runForLoop(loop forLoop, stack stack, stacklevel int) (any, ArErr) {
	fromval, err := runVal(loop.from, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	if typeof(fromval) != "number" {
		return nil, ArErr{"Type Error", "for loop from value must be a number", loop.line, loop.path, loop.code, true}
	}
	from := fromval.(number)
	toval, err := runVal(loop.to, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	if typeof(toval) != "number" {
		return nil, ArErr{"Type Error", "for loop to value must be a number", loop.line, loop.path, loop.code, true}
	}
	to := toval.(number)
	stepval, err := runVal(loop.step, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	if typeof(stepval) != "number" {
		return nil, ArErr{"Type Error", "for loop step value must be a number", loop.line, loop.path, loop.code, true}
	}
	step := stepval.(number)
	for i := newNumber().Set(from); i.Cmp(to) == -1; i = i.Add(i, step) {
		resp, err := runVal(loop.body, append(stack, Map(anymap{
			loop.variable: newNumber().Set(i),
		})), stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		switch x := resp.(type) {
		case Return:
			return x, ArErr{}
		case Break:
			return nil, ArErr{}
		case Continue:
			continue
		}
	}
	return nil, ArErr{}
}
