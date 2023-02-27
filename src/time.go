package main

import (
	"time"
)

var MicroSeconds = newNumber().SetInt64(1000000)

func ArTimeClass(N time.Time) ArClass {
	return ArClass{
		newNumber().Quo(newNumber().SetInt64(N.UnixMicro()), MicroSeconds),
		ArMap{
			"year": builtinFunc{
				"year",
				func(a ...any) (any, ArErr) {
					return newNumber().SetInt64(int64(N.Year())), ArErr{}
				},
			},
			"month": builtinFunc{
				"month",
				func(a ...any) (any, ArErr) {
					return N.Month().String(), ArErr{}
				},
			},
			"day": builtinFunc{
				"day",
				func(a ...any) (any, ArErr) {
					return newNumber().SetInt64(int64(N.Day())), ArErr{}
				},
			},
			"hour": builtinFunc{
				"hour",
				func(a ...any) (any, ArErr) {
					return newNumber().SetInt64(int64(N.Hour())), ArErr{}
				},
			},
			"minute": builtinFunc{
				"minute",
				func(a ...any) (any, ArErr) {
					return newNumber().SetInt64(int64(N.Minute())), ArErr{}
				},
			},
			"second": builtinFunc{
				"second",
				func(a ...any) (any, ArErr) {
					return newNumber().SetInt64(int64(N.Second())), ArErr{}
				},
			},
			"nanosecond": builtinFunc{
				"nanosecond",
				func(a ...any) (any, ArErr) {
					return newNumber().SetInt64(int64(N.Nanosecond())), ArErr{}
				},
			},
			"weekday": builtinFunc{
				"weekday",
				func(a ...any) (any, ArErr) {
					return N.Weekday().String(), ArErr{}
				},
			},
			"yearDay": builtinFunc{
				"yearDay",
				func(a ...any) (any, ArErr) {
					return newNumber().SetInt64(int64(N.YearDay())), ArErr{}
				},
			},
			"unix": builtinFunc{
				"unix",
				func(a ...any) (any, ArErr) {
					return newNumber().SetInt64(N.Unix()), ArErr{}
				},
			},
			"unixNano": builtinFunc{
				"unixNano",
				func(a ...any) (any, ArErr) {
					return newNumber().SetInt64(N.UnixNano()), ArErr{}
				},
			},
			"unixMilli": builtinFunc{
				"unixMilli",
				func(a ...any) (any, ArErr) {
					return newNumber().SetInt64(N.UnixMilli()), ArErr{}
				},
			},
			"unixMicro": builtinFunc{
				"unixMicro",
				func(a ...any) (any, ArErr) {
					return newNumber().SetInt64(N.UnixMicro()), ArErr{}
				},
			},
			"format": builtinFunc{
				"date",
				func(a ...any) (any, ArErr) {
					if len(a) == 0 {
						return N.Format(time.UnixDate), ArErr{}
					}
					return N.Format(a[0].(string)), ArErr{}
				},
			},
		}}
}

var ArTime = map[any]any{
	"snooze": builtinFunc{"snooze", func(a ...any) (any, ArErr) {
		if len(a) > 0 {
			float, _ := a[0].(number).Float64()
			time.Sleep(time.Duration(float*1000000000) * time.Nanosecond)
		}
		return nil, ArErr{}
	}},
	"now": builtinFunc{"now", func(a ...any) (any, ArErr) {
		return ArTimeClass(time.Now()), ArErr{}
	}},
	"parse": builtinFunc{"parse", func(a ...any) (any, ArErr) {
		if len(a) == 1 {
			N, err := time.Parse(time.UnixDate, a[0].(string))
			if err != nil {
				return nil, ArErr{
					TYPE:    "ArErr",
					message: err.Error(),
				}
			}
			return ArTimeClass(N), ArErr{}
		} else if len(a) > 1 {
			N, err := time.Parse(a[0].(string), a[1].(string))
			if err != nil {
				return nil, ArErr{
					TYPE:    "ArErr",
					message: err.Error(),
					EXISTS:  true,
				}
			}
			return ArTimeClass(N), ArErr{}
		}
		return nil, ArErr{
			TYPE:    "ArErr",
			message: "parse requires 2 arguments",
			EXISTS:  true,
		}
	}},
	"parseInLocation": builtinFunc{"parseInLocation", func(a ...any) (any, ArErr) {
		if len(a) > 2 {
			N, err := time.ParseInLocation(a[0].(string), a[1].(string), time.Local)
			if err != nil {
				return nil, ArErr{
					TYPE:    "ArErr",
					message: err.Error(),
					EXISTS:  true,
				}
			}
			return ArTimeClass(N), ArErr{}
		}
		return nil, ArErr{
			TYPE:    "ArErr",
			message: "parseInLocation requires 3 arguments",
			EXISTS:  true,
		}
	},
	},
	"date": builtinFunc{"date", func(a ...any) (any, ArErr) {
		if len(a) > 0 {
			N, err := time.Parse(time.UnixDate, a[0].(string))
			if err != nil {
				return nil, ArErr{
					TYPE:    "ArErr",
					message: err.Error(),
					EXISTS:  true,
				}
			}
			return ArTimeClass(N), ArErr{}
		}
		return nil, ArErr{
			TYPE:    "ArErr",
			message: "date requires 1 argument",
			EXISTS:  true,
		}
	},
	},
	"Unix": builtinFunc{"Unix", func(a ...any) (any, ArErr) {
		if len(a) > 1 {
			sec, _ := a[0].(number).Float64()
			nsec, _ := a[1].(number).Float64()
			return ArTimeClass(time.Unix(int64(sec), int64(nsec))), ArErr{}
		}
		return nil, ArErr{
			TYPE:    "ArErr",
			message: "Unix requires 2 arguments",
			EXISTS:  true,
		}
	},
	},
	"UnixMilli": builtinFunc{"UnixMilli", func(a ...any) (any, ArErr) {
		if len(a) > 0 {
			msec, _ := a[0].(number).Float64()
			return ArTimeClass(time.UnixMilli(int64(msec))), ArErr{}
		}
		return nil, ArErr{
			TYPE:    "ArErr",
			message: "UnixMilli requires 1 argument",
			EXISTS:  true,
		}
	},
	},
	"UnixMicro": builtinFunc{"UnixMicro", func(a ...any) (any, ArErr) {
		if len(a) > 0 {
			usec, _ := a[0].(number).Float64()
			return ArTimeClass(time.UnixMicro(int64(usec))), ArErr{}
		}
		return nil, ArErr{
			TYPE:    "ArErr",
			message: "UnixMicro requires 1 argument",
			EXISTS:  true,
		}
	},
	},
}
