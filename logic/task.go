package logic

import (
	"bytes"
	"encoding/json"
	"os"
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

func LoadTasksFromFile(filename string) ([]Task, error) {
	ret := []Task{}

	if jsonBytes, err := os.ReadFile(filename); err != nil {
		return nil, err

	} else if err := json.NewDecoder(bytes.NewReader(jsonBytes)).Decode(&ret); err != nil {
		return nil, err

	} else {
		return ret, nil
	}
}

func SaveTasksFromFile(filename string, tasks []Task) error {
	jsonBuffer := bytes.NewBuffer(nil)

	encoder := json.NewEncoder(jsonBuffer)
	encoder.SetIndent("", "\t")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(tasks); err != nil {
		return err

	} else if err := os.WriteFile(filename, jsonBuffer.Bytes(), os.ModePerm); err != nil {
		return err

	} else {
		return nil
	}
}
