package setting

import (
	"encoding/json"
	"fmt"

	"github.com/cwchiu/go-winapi"
)

type colorHexString winapi.COLORREF

func (crw *colorHexString) MarshalJSON() ([]byte, error) {
	var colorRef winapi.COLORREF = winapi.COLORREF(*crw)
	r := winapi.GetRValue(colorRef)
	g := winapi.GetGValue(colorRef)
	b := winapi.GetBValue(colorRef)

	str := fmt.Sprintf("#%02x%02x%02x", r, g, b)
	return []byte(str), nil
}

func (crw *colorHexString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	var r, g, b int32
	if _, err := fmt.Sscanf(str, "#%02x%02x%02x", &r, &g, &b); err != nil {
		return err
	}

	*crw = colorHexString(winapi.RGB(r, g, b))
	return nil
}
