package setting

import (
	"encoding/json"
	"fmt"

	"github.com/cwchiu/go-winapi"
)

type colorHexString winapi.COLORREF

func (crw *colorHexString) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("#%06x", *crw)), nil
}

func (crw *colorHexString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	_, err := fmt.Sscanf(str, "#%06x", crw)
	return err
}
