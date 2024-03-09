package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cwchiu/go-winapi"
)

type Settings struct {
	TargetDisplayIndex  int
	MeterWidth          int
	MeterOpacity        byte
	PastDuration        time.Duration
	FutureDuration      time.Duration
	ScaleInterval       time.Duration
	ScheduleEditCommand string
	BackgroundColor     winapi.COLORREF
	MainScaleColor      winapi.COLORREF
	SubScalesColor      winapi.COLORREF
	ChartColor          winapi.COLORREF
	TipTextColor        winapi.COLORREF
}

type ColorRefWrapper winapi.COLORREF
type DurationMinute time.Duration

func (crw *ColorRefWrapper) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("#%06x", *crw)), nil
}

func (crw *ColorRefWrapper) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	_, err := fmt.Sscanf(str, "#%06x", crw)
	return err
}

func (dm *DurationMinute) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Duration(*dm)/time.Minute)), nil
}

func (dm *DurationMinute) UnmarshalJSON(data []byte) error {
	var minutes time.Duration
	if err := json.Unmarshal(data, &minutes); err != nil {
		return err
	}

	*dm = DurationMinute(minutes * time.Minute)
	return nil
}

func (s *Settings) Default() {
	s.TargetDisplayIndex = 0
	s.MeterWidth = 50
	s.MeterOpacity = 128
	s.PastDuration = time.Hour * 1
	s.FutureDuration = time.Hour * 3
	s.ScaleInterval = time.Hour * 1
	s.ScheduleEditCommand = "notepad"
	s.BackgroundColor = winapi.RGB(0, 0, 0)
	s.MainScaleColor = winapi.RGB(255, 255, 255)
	s.SubScalesColor = winapi.RGB(128, 128, 128)
	s.ChartColor = winapi.RGB(255, 128, 0)
	s.TipTextColor = winapi.RGB(255, 255, 255)
}

func (s *Settings) LoadFile(filename string) error {
	var rawSettings struct {
		TargetDisplayIndex   *int             `json:"target_display_index,omitempty"`
		MeterWidth           *int             `json:"meter_width,omitempty"`
		MeterOpacity         *byte            `json:"meter_opacity,omitempty"`
		PastMinutes          *DurationMinute  `json:"past_minutes,omitempty"`
		FutureMinutes        *DurationMinute  `json:"future_minutes,omitempty"`
		ScaleIntervalMinutes *DurationMinute  `json:"scale_interval_minutes,omitempty"`
		ScheduleEditCommand  *string          `json:"schedule_edit_command,omitempty"`
		BackgroundColor      *ColorRefWrapper `json:"background_color,omitempty"`
		MainScaleColor       *ColorRefWrapper `json:"main_scale_color,omitempty"`
		SubScalesColor       *ColorRefWrapper `json:"sub_scales_color,omitempty"`
		ChartColor           *ColorRefWrapper `json:"chart_color,omitempty"`
		TipTextColor         *ColorRefWrapper `json:"tip_text_color,omitempty"`
	}

	if jsonBytes, err := os.ReadFile(filename); err != nil {
		return err

	} else if err := json.NewDecoder(bytes.NewReader(jsonBytes)).Decode(&rawSettings); err != nil {
		return err
	}

	tryAssign(&s.TargetDisplayIndex, rawSettings.TargetDisplayIndex)
	tryAssign(&s.MeterWidth, rawSettings.MeterWidth)
	tryAssign(&s.MeterOpacity, rawSettings.MeterOpacity)
	tryAssign(&s.PastDuration, (*time.Duration)(rawSettings.PastMinutes))
	tryAssign(&s.FutureDuration, (*time.Duration)(rawSettings.FutureMinutes))
	tryAssign(&s.ScaleInterval, (*time.Duration)(rawSettings.ScaleIntervalMinutes))
	tryAssign(&s.ScheduleEditCommand, rawSettings.ScheduleEditCommand)
	tryAssign(&s.BackgroundColor, (*winapi.COLORREF)(rawSettings.BackgroundColor))
	tryAssign(&s.MainScaleColor, (*winapi.COLORREF)(rawSettings.MainScaleColor))
	tryAssign(&s.SubScalesColor, (*winapi.COLORREF)(rawSettings.SubScalesColor))
	tryAssign(&s.ChartColor, (*winapi.COLORREF)(rawSettings.ChartColor))
	tryAssign(&s.TipTextColor, (*winapi.COLORREF)(rawSettings.TipTextColor))

	return nil
}

func tryAssign[T any](d *T, s *T) {
	if s != nil {
		*d = *s
	}
}
