package main

import "fmt"

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
					message: "expected string or []byte, got " + typeof(x),
					EXISTS:  true,
				}
			}
			obj.obj["__value__"] = buf
			obj.obj["length"] = newNumber().SetInt64(int64(len(buf)))
			return obj, ArErr{}
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
	return obj
}
