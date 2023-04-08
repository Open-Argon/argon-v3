package main

import "fmt"

func ArSequence(a ...any) (any, ArErr) {
	if len(a) < 1 || len(a) > 2 {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("sequence expected 1 or 2 arguments, got %d", len(a)),
			EXISTS:  true,
		}
	}
	f := a[0]
	initial := newNumber()
	if typeof(f) != "function" {
		return nil, ArErr{TYPE: "Runtime Error",
			message: fmt.Sprintf("sequence expected function, got %s", typeof(f)),
			EXISTS:  true,
		}
	}
	if len(a) == 2 {
		if typeof(a[1]) != "number" {
			return nil, ArErr{TYPE: "Runtime Error",
				message: fmt.Sprintf("sequence expected number, got %s", typeof(a[1])),
				EXISTS:  true,
			}
		}
		initial.Set(a[1].(number))
	}
	return ArObject{
		obj: map[any]any{
			"__name__":  "sequence",
			"__value__": "test",
		},
	}, ArErr{}
}
