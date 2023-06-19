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
					TYPE:    "TypeError",
					message: "set() takes exactly 2 argument (" + fmt.Sprint(len(a)) + " given)",
					EXISTS:  true,
				}
			}
			var c *color.Color
			var s string
			if x, ok := a[0].(number); ok {
				c = color.Set(color.Attribute(x.Num().Int64()))
			} else {
				return nil, ArErr{
					TYPE:    "TypeError",
					message: "set() argument 1 must be an number, not " + typeof(a[0]),
					EXISTS:  true,
				}
			}
			if typeof(a[1]) == "string" {
				s = ArValidToAny(a[1]).(string)
			} else {
				return nil, ArErr{
					TYPE:    "TypeError",
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
				"black":     newNumber().SetInt64(int64(color.BgBlack)),
				"red":       newNumber().SetInt64(int64(color.BgRed)),
				"green":     newNumber().SetInt64(int64(color.BgGreen)),
				"yellow":    newNumber().SetInt64(int64(color.BgYellow)),
				"blue":      newNumber().SetInt64(int64(color.BgBlue)),
				"magenta":   newNumber().SetInt64(int64(color.BgMagenta)),
				"cyan":      newNumber().SetInt64(int64(color.BgCyan)),
				"white":     newNumber().SetInt64(int64(color.BgWhite)),
				"hiBlack":   newNumber().SetInt64(int64(color.BgHiBlack)),
				"hiRed":     newNumber().SetInt64(int64(color.BgHiRed)),
				"hiGreen":   newNumber().SetInt64(int64(color.BgHiGreen)),
				"hiYellow":  newNumber().SetInt64(int64(color.BgHiYellow)),
				"hiBlue":    newNumber().SetInt64(int64(color.BgHiBlue)),
				"hiMagenta": newNumber().SetInt64(int64(color.BgHiMagenta)),
				"hiCyan":    newNumber().SetInt64(int64(color.BgHiCyan)),
				"hiWhite":   newNumber().SetInt64(int64(color.BgHiWhite)),
			},
		),
		"fg": Map(
			anymap{
				"black":     newNumber().SetInt64(int64(color.FgBlack)),
				"red":       newNumber().SetInt64(int64(color.FgRed)),
				"green":     newNumber().SetInt64(int64(color.FgGreen)),
				"yellow":    newNumber().SetInt64(int64(color.FgYellow)),
				"blue":      newNumber().SetInt64(int64(color.FgBlue)),
				"magenta":   newNumber().SetInt64(int64(color.FgMagenta)),
				"cyan":      newNumber().SetInt64(int64(color.FgCyan)),
				"white":     newNumber().SetInt64(int64(color.FgWhite)),
				"hiBlack":   newNumber().SetInt64(int64(color.FgHiBlack)),
				"hiRed":     newNumber().SetInt64(int64(color.FgHiRed)),
				"hiGreen":   newNumber().SetInt64(int64(color.FgHiGreen)),
				"hiYellow":  newNumber().SetInt64(int64(color.FgHiYellow)),
				"hiBlue":    newNumber().SetInt64(int64(color.FgHiBlue)),
				"hiMagenta": newNumber().SetInt64(int64(color.FgHiMagenta)),
				"hiCyan":    newNumber().SetInt64(int64(color.FgHiCyan)),
				"hiWhite":   newNumber().SetInt64(int64(color.FgHiWhite)),
			},
		),
		"reset":        newNumber().SetInt64(int64(color.Reset)),
		"bold":         newNumber().SetInt64(int64(color.Bold)),
		"faint":        newNumber().SetInt64(int64(color.Faint)),
		"italic":       newNumber().SetInt64(int64(color.Italic)),
		"underline":    newNumber().SetInt64(int64(color.Underline)),
		"blinkSlow":    newNumber().SetInt64(int64(color.BlinkSlow)),
		"blinkRapid":   newNumber().SetInt64(int64(color.BlinkRapid)),
		"reverseVideo": newNumber().SetInt64(int64(color.ReverseVideo)),
		"concealed":    newNumber().SetInt64(int64(color.Concealed)),
		"crossedOut":   newNumber().SetInt64(int64(color.CrossedOut)),
	})
