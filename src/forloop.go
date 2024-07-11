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
			stepval = _one_Number
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
		bodyval, worked, err, bodystep := translateVal(UNPARSEcode{code: body, realcode: code.realcode, line: code.line, path: code.path}, index, codelines, 3)
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
	from := fromval.(ArObject)
	toval, err := runVal(loop.to, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	if typeof(toval) != "number" {
		return nil, ArErr{"Type Error", "for loop to value must be a number", loop.line, loop.path, loop.code, true}
	}
	to := toval.(ArObject)
	stepval, err := runVal(loop.step, stack, stacklevel+1)
	if err.EXISTS {
		return nil, err
	}
	if typeof(stepval) != "number" {
		return nil, ArErr{"Type Error", "for loop step value must be a number", loop.line, loop.path, loop.code, true}
	}
	step := stepval.(ArObject)
	if isNumberInt64(from) && isNumberInt64(to) && isNumberInt64(step) {
		i, _ := numberToInt64(from)
		to_, _ := numberToInt64(to)
		step_, _ := numberToInt64(step)
		layer := anymap{}
		stacks := append(stack, Map(layer))
		for i < to_ {
			layer[loop.variable] = Number(i)
			resp, err := runVal(loop.body, stacks, stacklevel+1)
			if err.EXISTS {
				return nil, err
			}
			switch x := resp.(type) {
			case Return:
				return x, ArErr{}
			case Break:
				return nil, ArErr{}
			case Continue:
			}
			i += step_
		}
		return nil, ArErr{}
	}
	i := from
	direction_obj, err := CompareObjects(step, _zero_Number)
	if err.EXISTS {
		return nil, err
	}
	currentDirection_obj, err := CompareObjects(to, i)
	if err.EXISTS {
		return nil, err
	}
	currentDirection, error := numberToInt64(currentDirection_obj)
	if error != nil {
		return nil, ArErr{"Type Error", error.Error(), loop.line, loop.path, loop.code, true}
	}
	direction, error := numberToInt64(direction_obj)
	if error != nil {
		return nil, ArErr{"Type Error", error.Error(), loop.line, loop.path, loop.code, true}
	}
	layer := anymap{}
	stacks := append(stack, Map(layer))
	for currentDirection == direction {
		layer[loop.variable] = i
		resp, err := runVal(loop.body, stacks, stacklevel+1)
		if err.EXISTS {
			return nil, err
		}
		switch x := resp.(type) {
		case Return:
			return x, ArErr{}
		case Break:
			return nil, ArErr{}
		case Continue:
		}
		i, err = AddObjects(i, step)
		if err.EXISTS {
			return nil, err
		}
		currentDirection_obj, err = CompareObjects(to, i)
		if err.EXISTS {
			return nil, err
		}
		currentDirection, error = numberToInt64(currentDirection_obj)
		if error != nil {
			return nil, ArErr{"Type Error", error.Error(), loop.line, loop.path, loop.code, true}
		}
		direction, error = numberToInt64(direction_obj)
		if error != nil {
			return nil, ArErr{"Type Error", error.Error(), loop.line, loop.path, loop.code, true}
		}
	}
	return nil, ArErr{}
}
