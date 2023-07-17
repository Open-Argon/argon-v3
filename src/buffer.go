package main

import (
	"bytes"
	"fmt"
)

func ArByte(Byte byte) ArObject {
	obj := ArObject{
		obj: anymap{
			"__name__":  "byte",
			"__value__": Byte,
		},
	}
	obj.obj["__string__"] = builtinFunc{
		"__string__",
		func(a ...any) (any, ArErr) {
			return "<byte>", ArErr{}
		},
	}
	obj.obj["__repr__"] = builtinFunc{
		"__repr__",
		func(a ...any) (any, ArErr) {
			return "<byte>", ArErr{}
		},
	}
	obj.obj["number"] = builtinFunc{
		"number",
		func(a ...any) (any, ArErr) {
			return newNumber().SetInt64(int64(Byte)), ArErr{}
		},
	}
	obj.obj["from"] = builtinFunc{
		"from",
		func(a ...any) (any, ArErr) {
			if len(a) == 0 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected at least 1 argument, got 0",
					EXISTS:  true,
				}
			}
			a[0] = ArValidToAny(a[0])
			switch x := a[0].(type) {
			case number:
				if x.Denom().Cmp(one.Denom()) != 0 {
					return nil, ArErr{
						TYPE:    "TypeError",
						message: "expected integer, got " + fmt.Sprint(x),
						EXISTS:  true,
					}
				}
				n := x.Num().Int64()
				if n > 255 || n < 0 {
					return nil, ArErr{
						TYPE:    "ValueError",
						message: "expected number between 0 and 255, got " + fmt.Sprint(floor(x).Num().Int64()),
						EXISTS:  true,
					}
				}
				Byte = byte(n)
			case string:
				if len(x) != 1 {
					return nil, ArErr{
						TYPE:    "ValueError",
						message: "expected string of length 1, got " + fmt.Sprint(len(x)),
						EXISTS:  true,
					}
				}
				Byte = byte(x[0])
			default:
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected number or string, got " + typeof(x),
					EXISTS:  true,
				}
			}
			return obj, ArErr{}
		},
	}
	return obj
}

func ArBuffer(buf []byte) ArObject {
	obj := ArObject{
		obj: anymap{
			"__name__":  "buffer",
			"__value__": buf,
			"length":    newNumber().SetInt64(int64(len(buf))),
		},
	}
	obj.obj["__string__"] = builtinFunc{
		"__string__",
		func(a ...any) (any, ArErr) {
			return "<buffer>", ArErr{}
		},
	}
	obj.obj["__repr__"] = builtinFunc{
		"__repr__",
		func(a ...any) (any, ArErr) {
			return "<buffer>", ArErr{}
		},
	}
	obj.obj["from"] = builtinFunc{
		"from",
		func(a ...any) (any, ArErr) {
			if len(a) == 0 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected at least 1 argument, got 0",
					EXISTS:  true,
				}
			}
			a[0] = ArValidToAny(a[0])
			switch x := a[0].(type) {
			case string:
				buf = []byte(x)
			case []byte:
				buf = x
			case []any:
				outputbuf := []byte{}
				for _, v := range x {
					switch y := v.(type) {
					case number:
						if y.Denom().Cmp(one.Denom()) != 0 {
							return nil, ArErr{
								TYPE:    "TypeError",
								message: "Cannot convert non-integer to byte",
								EXISTS:  true,
							}
						}
						outputbuf = append(outputbuf, byte(y.Num().Int64()))
					default:
						return nil, ArErr{
							TYPE:    "TypeError",
							message: "Cannot convert " + typeof(v) + " to byte",
							EXISTS:  true,
						}
					}
				}
				buf = outputbuf
			default:
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected string or buffer, got " + typeof(x),
					EXISTS:  true,
				}
			}
			obj.obj["__value__"] = buf
			obj.obj["length"] = newNumber().SetInt64(int64(len(buf)))
			return obj, ArErr{}
		},
	}
	obj.obj["splitN"] = builtinFunc{
		"splitN",
		func(a ...any) (any, ArErr) {
			if len(a) != 2 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected 1 argument, got " + fmt.Sprint(len(a)),
					EXISTS:  true,
				}
			}
			splitVal := ArValidToAny(a[0])
			if typeof(splitVal) != "buffer" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected buffer, got " + typeof(splitVal),
					EXISTS:  true,
				}
			}
			var separator = splitVal.([]byte)
			nVal := ArValidToAny(a[1])
			if typeof(nVal) != "number" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected number, got " + typeof(nVal),
					EXISTS:  true,
				}
			}
			nNum := nVal.(number)
			if nNum.Denom().Cmp(one.Denom()) != 0 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected integer, got " + fmt.Sprint(nNum),
					EXISTS:  true,
				}
			}
			n := nNum.Num().Int64()
			var result [][]byte
			start := 0
			var i int64
			for i = 0; i < n; i++ {
				index := bytes.Index(buf[start:], separator)
				if index == -1 {
					result = append(result, buf[start:])
					break
				}
				end := start + index
				result = append(result, buf[start:end])
				start = end + len(separator)
			}

			if int64(len(result)) != n {
				result = append(result, buf[start:])
			}
			var bufoutput = []any{}
			for _, v := range result {
				bufoutput = append(bufoutput, ArBuffer(v))
			}
			return ArArray(bufoutput), ArErr{}
		},
	}
	obj.obj["split"] = builtinFunc{
		"split",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected 1 argument, got " + fmt.Sprint(len(a)),
					EXISTS:  true,
				}
			}
			splitVal := ArValidToAny(a[0])
			if typeof(splitVal) != "buffer" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected buffer, got " + typeof(splitVal),
					EXISTS:  true,
				}
			}
			var separator = splitVal.([]byte)
			var result [][]byte
			start := 0

			for {
				index := bytes.Index(buf[start:], separator)
				if index == -1 {
					result = append(result, buf[start:])
					break
				}
				end := start + index
				result = append(result, buf[start:end])
				start = end + len(separator)
			}
			var bufoutput = []any{}
			for _, v := range result {
				bufoutput = append(bufoutput, ArBuffer(v))
			}
			return ArArray(bufoutput), ArErr{}
		},
	}
	obj.obj["slice"] = builtinFunc{
		"slice",
		func(a ...any) (any, ArErr) {
			if len(a) != 2 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected 2 arguments, got " + fmt.Sprint(len(a)),
					EXISTS:  true,
				}
			}
			startVal := ArValidToAny(a[0])
			if typeof(startVal) != "number" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected number, got " + typeof(startVal),
					EXISTS:  true,
				}
			}
			start := startVal.(number)
			if start.Denom().Cmp(one.Denom()) != 0 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected integer, got " + fmt.Sprint(start),
					EXISTS:  true,
				}
			}
			endVal := ArValidToAny(a[1])
			if typeof(endVal) != "number" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected number, got " + typeof(endVal),
					EXISTS:  true,
				}
			}
			end := endVal.(number)
			if end.Denom().Cmp(one.Denom()) != 0 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected integer, got " + fmt.Sprint(end),
					EXISTS:  true,
				}
			}
			return ArBuffer(buf[floor(start).Num().Int64():floor(end).Num().Int64()]), ArErr{}
		},
	}
	obj.obj["to"] = builtinFunc{
		"to",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected 1 argument, got " + fmt.Sprint(len(a)),
					EXISTS:  true,
				}
			}
			if typeof(a[0]) != "string" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected string, got " + typeof(a[0]),
					EXISTS:  true,
				}
			}
			Type := ArValidToAny(a[0]).(string)
			switch Type {
			case "string":
				return ArString(string(buf)), ArErr{}
			case "bytes":
				output := []any{}
				for _, v := range buf {
					output = append(output, ArByte(v))
				}
				return ArArray(output), ArErr{}
			case "array":
				output := []any{}
				for _, v := range buf {
					output = append(output, newNumber().SetInt64(int64(v)))
				}
				return ArArray(output), ArErr{}
			default:
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected string, bytes or array, got '" + Type + "'",
					EXISTS:  true,
				}
			}
		},
	}
	obj.obj["append"] = builtinFunc{
		"append",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected 1 argument, got " + fmt.Sprint(len(a)),
					EXISTS:  true,
				}
			}
			a[0] = ArValidToAny(a[0])
			switch x := a[0].(type) {
			case number:
				if x.Denom().Cmp(one.Denom()) != 0 {
					return nil, ArErr{
						TYPE:    "TypeError",
						message: "Cannot convert non-integer to byte",
						EXISTS:  true,
					}
				}
				buf = append(buf, byte(x.Num().Int64()))
			case string:
				buf = append(buf, []byte(x)...)
			case []byte:
				buf = append(buf, x...)
			case []any:
				for _, v := range x {
					switch y := v.(type) {
					case number:
						if y.Denom().Cmp(one.Denom()) != 0 {
							return nil, ArErr{
								TYPE:    "TypeError",
								message: "Cannot convert non-integer to byte",
								EXISTS:  true,
							}
						}
						buf = append(buf, byte(y.Num().Int64()))
					default:
						return nil, ArErr{
							TYPE:    "TypeError",
							message: "Cannot convert " + typeof(v) + " to byte",
							EXISTS:  true,
						}
					}
				}
			default:
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected string, buffer or array, got " + typeof(x),
					EXISTS:  true,
				}
			}
			obj.obj["__value__"] = buf
			obj.obj["length"] = newNumber().SetInt64(int64(len(buf)))
			return obj, ArErr{}
		},
	}
	obj.obj["insert"] = builtinFunc{
		"insert",
		func(a ...any) (any, ArErr) {
			if len(a) != 2 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected 2 arguments, got " + fmt.Sprint(len(a)),
					EXISTS:  true,
				}
			}
			poss := ArValidToAny(a[0])
			values := ArValidToAny(a[1])
			if typeof(poss) != "number" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected number, got " + typeof(poss),
					EXISTS:  true,
				}
			}
			pos := poss.(number)
			if pos.Denom().Cmp(one.Denom()) != 0 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "position must be an integer",
					EXISTS:  true,
				}
			}
			posNum := pos.Num().Int64()
			switch x := values.(type) {
			case number:
				if x.Denom().Cmp(one.Denom()) != 0 {
					return nil, ArErr{
						TYPE:    "TypeError",
						message: "Cannot convert non-integer to byte",
						EXISTS:  true,
					}
				}
				buf = append(buf[:posNum], append([]byte{byte(x.Num().Int64())}, buf[posNum:]...)...)
			case string:
				buf = append(buf[:posNum], append([]byte(x), buf[posNum:]...)...)
			case []byte:
				buf = append(buf[:posNum], append(x, buf[posNum:]...)...)
			case []any:
				for _, v := range x {
					switch y := v.(type) {
					case number:
						if y.Denom().Cmp(one.Denom()) != 0 {
							return nil, ArErr{
								TYPE:    "TypeError",
								message: "Cannot convert non-integer to byte",
								EXISTS:  true,
							}
						}
						buf = append(buf[:posNum], append([]byte{byte(y.Num().Int64())}, buf[posNum:]...)...)
					default:
						return nil, ArErr{
							TYPE:    "TypeError",
							message: "Cannot convert " + typeof(v) + " to byte",
							EXISTS:  true,
						}
					}
				}
			default:
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected string or buffer, got " + typeof(x),
					EXISTS:  true,
				}
			}
			obj.obj["__value__"] = buf
			obj.obj["length"] = newNumber().SetInt64(int64(len(buf)))
			return obj, ArErr{}
		},
	}
	obj.obj["remove"] = builtinFunc{
		"remove",
		func(a ...any) (any, ArErr) {
			if len(a) != 1 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected 1 argument, got " + fmt.Sprint(len(a)),
					EXISTS:  true,
				}
			}
			poss := ArValidToAny(a[0])
			if typeof(poss) != "number" {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "expected number, got " + typeof(poss),
					EXISTS:  true,
				}
			}
			pos := poss.(number)
			if pos.Denom().Cmp(one.Denom()) != 0 {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "position must be an integer",
					EXISTS:  true,
				}
			}
			posNum := pos.Num().Int64()
			buf = append(buf[:posNum], buf[posNum+1:]...)
			obj.obj["__value__"] = buf
			obj.obj["length"] = newNumber().SetInt64(int64(len(buf)))
			return obj, ArErr{}
		},
	}

	return obj
}
