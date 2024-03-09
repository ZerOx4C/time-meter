package wrapped

import (
	"github.com/cwchiu/go-winapi"
)

type RECT winapi.RECT

func (r *RECT) Unwrap() *winapi.RECT {
	return (*winapi.RECT)(r)
}

func (r *RECT) Translate(x, y int32) {
	r.Left += x
	r.Right += x
	r.Top += y
	r.Bottom += y
}

func (r *RECT) Width() int32 {
	return r.Right - r.Left
}

func (r *RECT) Height() int32 {
	return r.Bottom - r.Top
}

func (r *RECT) Contains(p POINT) bool {
	if p.X < r.Left || r.Right < p.X {
		return false
	}

	if p.Y < r.Top || r.Bottom < p.Y {
		return false
	}

	return true
}
