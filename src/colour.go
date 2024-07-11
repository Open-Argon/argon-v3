package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jwalton/go-supportscolor"
)

var ArColour = Map(
	anymap{
		"set": builtinFunc{"set", func(a ...any) (any, ArErr) {
			if len(a) != 2 {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "set() takes exactly 2 argument (" + fmt.Sprint(len(a)) + " given)",
					EXISTS:  true,
				}
			}
			var c *color.Color
			var s string
			if x, ok := a[0].(ArObject); ok {
				colour_int64, err := numberToInt64(x)
				if err != nil {
					return nil, ArErr{
						TYPE:    "Type Error",
						message: err.Error(),
						EXISTS:  true,
					}
				}
				c = color.New(color.Attribute(colour_int64))
			} else {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "set() argument 1 must be an number, not " + typeof(a[0]),
					EXISTS:  true,
				}
			}
			if typeof(a[1]) == "string" {
				s = ArValidToAny(a[1]).(string)
			} else {
				return nil, ArErr{
					TYPE:    "Type Error",
					message: "set() argument 2 must be a string, not " + typeof(a[1]),
					EXISTS:  true,
				}
			}
			if supportscolor.Stdout().SupportsColor {
				return c.Sprint(s), ArErr{}
			} else {
				return s, ArErr{}
			}
		}},
		"bg": Map(
			anymap{
				"black":     Number(int64(color.BgBlack)),
				"red":       Number(int64(color.BgRed)),
				"green":     Number(int64(color.BgGreen)),
				"yellow":    Number(int64(color.BgYellow)),
				"blue":      Number(int64(color.BgBlue)),
				"magenta":   Number(int64(color.BgMagenta)),
				"cyan":      Number(int64(color.BgCyan)),
				"white":     Number(int64(color.BgWhite)),
				"hiBlack":   Number(int64(color.BgHiBlack)),
				"hiRed":     Number(int64(color.BgHiRed)),
				"hiGreen":   Number(int64(color.BgHiGreen)),
				"hiYellow":  Number(int64(color.BgHiYellow)),
				"hiBlue":    Number(int64(color.BgHiBlue)),
				"hiMagenta": Number(int64(color.BgHiMagenta)),
				"hiCyan":    Number(int64(color.BgHiCyan)),
				"hiWhite":   Number(int64(color.BgHiWhite)),
			},
		),
		"fg": Map(
			anymap{
				"black":     Number(int64(color.FgBlack)),
				"red":       Number(int64(color.FgRed)),
				"green":     Number(int64(color.FgGreen)),
				"yellow":    Number(int64(color.FgYellow)),
				"blue":      Number(int64(color.FgBlue)),
				"magenta":   Number(int64(color.FgMagenta)),
				"cyan":      Number(int64(color.FgCyan)),
				"white":     Number(int64(color.FgWhite)),
				"hiBlack":   Number(int64(color.FgHiBlack)),
				"hiRed":     Number(int64(color.FgHiRed)),
				"hiGreen":   Number(int64(color.FgHiGreen)),
				"hiYellow":  Number(int64(color.FgHiYellow)),
				"hiBlue":    Number(int64(color.FgHiBlue)),
				"hiMagenta": Number(int64(color.FgHiMagenta)),
				"hiCyan":    Number(int64(color.FgHiCyan)),
				"hiWhite":   Number(int64(color.FgHiWhite)),
			},
		),
		"reset":        Number(int64(color.Reset)),
		"bold":         Number(int64(color.Bold)),
		"faint":        Number(int64(color.Faint)),
		"italic":       Number(int64(color.Italic)),
		"underline":    Number(int64(color.Underline)),
		"blinkSlow":    Number(int64(color.BlinkSlow)),
		"blinkRapid":   Number(int64(color.BlinkRapid)),
		"reverseVideo": Number(int64(color.ReverseVideo)),
		"concealed":    Number(int64(color.Concealed)),
		"crossedOut":   Number(int64(color.CrossedOut)),
	})
