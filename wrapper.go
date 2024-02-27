package main

import (
	"github.com/cwchiu/go-winapi"
)

type POINT winapi.POINT
type RECT winapi.RECT

func (p *POINT) unwrap() *winapi.POINT {
	return (*winapi.POINT)(p)
}

func (r *RECT) unwrap() *winapi.RECT {
	return (*winapi.RECT)(r)
}

func (r *RECT) translate(x, y int32) {
	r.Left += x
	r.Right += x
	r.Top += y
	r.Bottom += y
}

func (r *RECT) width() int32 {
	return r.Right - r.Left
}

func (r *RECT) height() int32 {
	return r.Bottom - r.Top
}

func (r *RECT) contains(p POINT) bool {
	if p.X < r.Left || r.Right < p.X {
		return false
	}

	if p.Y < r.Top || r.Bottom < p.Y {
		return false
	}

	return true
}
