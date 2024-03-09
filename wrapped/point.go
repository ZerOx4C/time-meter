package wrapped

import (
	"github.com/cwchiu/go-winapi"
)

type POINT winapi.POINT

func (p *POINT) Unwrap() *winapi.POINT {
	return (*winapi.POINT)(p)
}
