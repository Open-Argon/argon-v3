package main

import (
	"fmt"
	"reflect"
)

// returns (number|string|nil), error
func runVal(line any, stack stack, stacklevel int) (any, ArErr) {
	var (
		linenum       = 0
		path          = ""
		code          = ""
		stackoverflow = stacklevel >= 10000
	)
	switch x := line.(type) {
	case string:
		return ArString(x), ArErr{}
	case call:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runCall(x, stack, stacklevel+1)
	case factorial:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runFactorial(x, stack, stacklevel+1)
	case accessVariable:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return readVariable(x, stack)
	case ArMapGet:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return mapGet(x, stack, stacklevel+1)
	case setVariable:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return setVariableValue(x, stack, stacklevel+1)
	case negative:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		resp, err := runVal(x.VAL, stack, stacklevel+1)
		resp = AnyToArValid(resp)
		if err.EXISTS {
			return nil, err
		}
		switch y := resp.(type) {
		case number:
			if !x.sign {
				return newNumber().Neg(y), ArErr{}
			}
			return y, ArErr{}
		}
		return nil, ArErr{
			TYPE:    "TypeError",
			message: "cannot negate a non-number",
			EXISTS:  true,
		}
	case brackets:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runVal(x.VAL, stack, stacklevel+1)
	case operationType:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runOperation(x, stack, stacklevel+1)
	case dowrap:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runDoWrap(x, stack, stacklevel+1)
	case CallReturn:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runReturn(x, stack, stacklevel+1)
	case Break:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return x, ArErr{}
	case Continue:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return x, ArErr{}
	case ArDelete:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runDelete(x, stack, stacklevel+1)
	case not:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runNot(x, stack, stacklevel+1)
	case ifstatement:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runIfStatement(x, stack, stacklevel+1)
	case whileLoop:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runWhileLoop(x, stack, stacklevel+1)
	case forLoop:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runForLoop(x, stack, stacklevel+1)
	case CreateArray:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runArray(x, stack, stacklevel+1)
	case squareroot:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runSquareroot(x, stack, stacklevel+1)
	case createMap:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runCreateMap(x, stack, stacklevel+1)
	case ArImport:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runImport(x, stack, stacklevel+1)
	case ABS:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runAbs(x, stack, stacklevel+1)
	case TryCatch:
		if stackoverflow {
			linenum = x.line
			path = x.path
			code = x.code
			break
		}
		return runTryCatch(x, stack, stacklevel+1)
	case bool, ArObject, number, nil, Callable, builtinFunc, anymap:
		return x, ArErr{}
	}
	if stackoverflow {
		return nil, ArErr{
			TYPE:    "RuntimeError",
			message: "stack overflow",
			line:    linenum,
			path:    path,
			code:    code,
			EXISTS:  true,
		}
	}
	fmt.Println("unreachable", reflect.TypeOf(line))
	panic("unreachable")
}

// returns error
func run(translated []any, stack stack) (any, ArErr) {
	var output any = nil
	for _, val := range translated {
		val, err := runVal(val, stack, 0)
		output = val
		if err.EXISTS {
			return output, err
		}
	}
	return output, ArErr{}
}
