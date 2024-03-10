package logic

import (
	"time"
)

type Task struct {
	Subject string    `json:"subject"`
	BeginAt time.Time `json:"begin_at"`
	EndAt   time.Time `json:"end_at"`
}

func (t *Task) OverlapWith(beginAt time.Time, endAt time.Time) bool {
	if beginAt.Before(t.EndAt) && t.BeginAt.Before(endAt) {
		return true
	}

	return false
}
