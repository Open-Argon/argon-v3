package main

import "github.com/wadey/go-rounding"

func floor(x number) number {

	n := newNumber().Set(x)
	if n.Sign() < 0 {
		return rounding.Round(n, 0, rounding.Up)
	}
	return rounding.Round(n, 0, rounding.Down)
}

func ceil(x number) number {
	n := newNumber().Set(x)
	if n.Sign() < 0 {
		return rounding.Round(n, 0, rounding.Down)
	}
	return rounding.Round(n, 0, rounding.Up)
}

func round(x number, precision int) number {
	return rounding.Round(newNumber().Set(x), precision, rounding.HalfUp)
}
