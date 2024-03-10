package setting

import (
	"bytes"
	"encoding/json"
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
	Port                int
}

type nilableSettings struct {
	TargetDisplayIndex   *int            `json:"target_display_index,omitempty"`
	MeterWidth           *int            `json:"meter_width,omitempty"`
	MeterOpacity         *byte           `json:"meter_opacity,omitempty"`
	PastMinutes          *durationMinute `json:"past_minutes,omitempty"`
	FutureMinutes        *durationMinute `json:"future_minutes,omitempty"`
	ScaleIntervalMinutes *durationMinute `json:"scale_interval_minutes,omitempty"`
	ScheduleEditCommand  *string         `json:"schedule_edit_command,omitempty"`
	BackgroundColor      *colorHexString `json:"background_color,omitempty"`
	MainScaleColor       *colorHexString `json:"main_scale_color,omitempty"`
	SubScalesColor       *colorHexString `json:"sub_scales_color,omitempty"`
	ChartColor           *colorHexString `json:"chart_color,omitempty"`
	TipTextColor         *colorHexString `json:"tip_text_color,omitempty"`
	Port                 *int            `json:"port,omitempty"`
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
	s.Port = 50000
}

func (s *Settings) LoadFile(filename string) error {
	var nilable nilableSettings

	if jsonBytes, err := os.ReadFile(filename); err != nil {
		return err

	} else if err := json.NewDecoder(bytes.NewReader(jsonBytes)).Decode(&nilable); err != nil {
		return err
	}

	var settings Settings
	settings.Default()

	assignIfNotNil(&settings.TargetDisplayIndex, nilable.TargetDisplayIndex)
	assignIfNotNil(&settings.MeterWidth, nilable.MeterWidth)
	assignIfNotNil(&settings.MeterOpacity, nilable.MeterOpacity)
	assignIfNotNil(&settings.PastDuration, (*time.Duration)(nilable.PastMinutes))
	assignIfNotNil(&settings.FutureDuration, (*time.Duration)(nilable.FutureMinutes))
	assignIfNotNil(&settings.ScaleInterval, (*time.Duration)(nilable.ScaleIntervalMinutes))
	assignIfNotNil(&settings.ScheduleEditCommand, nilable.ScheduleEditCommand)
	assignIfNotNil(&settings.BackgroundColor, (*winapi.COLORREF)(nilable.BackgroundColor))
	assignIfNotNil(&settings.MainScaleColor, (*winapi.COLORREF)(nilable.MainScaleColor))
	assignIfNotNil(&settings.SubScalesColor, (*winapi.COLORREF)(nilable.SubScalesColor))
	assignIfNotNil(&settings.ChartColor, (*winapi.COLORREF)(nilable.ChartColor))
	assignIfNotNil(&settings.TipTextColor, (*winapi.COLORREF)(nilable.TipTextColor))
	assignIfNotNil(&settings.Port, nilable.Port)

	*s = settings

	return nil
}

func assignIfNotNil[T any](d *T, s *T) {
	if s != nil {
		*d = *s
	}
}
