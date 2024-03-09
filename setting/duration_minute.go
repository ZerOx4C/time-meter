package setting

import (
	"encoding/json"
	"fmt"
	"time"
)

type durationMinute time.Duration

func (dm *durationMinute) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Duration(*dm)/time.Minute)), nil
}

func (dm *durationMinute) UnmarshalJSON(data []byte) error {
	var minutes time.Duration
	if err := json.Unmarshal(data, &minutes); err != nil {
		return err
	}

	*dm = durationMinute(minutes * time.Minute)
	return nil
}
