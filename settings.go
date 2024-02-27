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
	MeterWidth      int
	MeterOpacity    byte
	PastDuration    time.Duration
	FutureDuration  time.Duration
	ScaleInterval   time.Duration
	BackgroundColor winapi.COLORREF
	MainScaleColor  winapi.COLORREF
	SubScalesColor  winapi.COLORREF
	ChartColor      winapi.COLORREF
	TipTextColor    winapi.COLORREF
}

func (s *Settings) Default() {
	s.MeterWidth = 50
	s.MeterOpacity = 128
	s.PastDuration = time.Hour * 1
	s.FutureDuration = time.Hour * 3
	s.ScaleInterval = time.Hour * 1
	s.BackgroundColor = winapi.RGB(0, 0, 0)
	s.MainScaleColor = winapi.RGB(255, 255, 255)
	s.SubScalesColor = winapi.RGB(128, 128, 128)
	s.ChartColor = winapi.RGB(255, 128, 0)
	s.TipTextColor = winapi.RGB(255, 255, 255)
}

func (s *Settings) LoadFile(filename string) error {
	var rawSettings struct {
		MeterWidth            *int    `json:"meter_width,omitempty"`
		MeterOpacity          *byte   `json:"meter_opacity,omitempty"`
		PastMinutes           *int    `json:"past_minutes,omitempty"`
		FutureMinutes         *int    `json:"future_minutes,omitempty"`
		ScaleIntervalMinutes  *int    `json:"scale_interval_minutes,omitempty"`
		BackgroundColorString *string `json:"background_color,omitempty"`
		MainScaleColorString  *string `json:"main_scale_color,omitempty"`
		SubScalesColorString  *string `json:"sub_scales_color,omitempty"`
		ChartColorString      *string `json:"chart_color,omitempty"`
		TipTextColorString    *string `json:"tip_text_color,omitempty"`
	}

	if jsonBytes, err := os.ReadFile(filename); err != nil {
		return err

	} else if err := json.NewDecoder(bytes.NewReader(jsonBytes)).Decode(&rawSettings); err != nil {
		return err
	}

	if rawSettings.MeterWidth != nil {
		s.MeterWidth = *rawSettings.MeterWidth
	}

	if rawSettings.MeterOpacity != nil {
		s.MeterOpacity = *rawSettings.MeterOpacity
	}

	if rawSettings.PastMinutes != nil {
		s.PastDuration = time.Minute * time.Duration(*rawSettings.PastMinutes)
	}

	if rawSettings.FutureMinutes != nil {
		s.FutureDuration = time.Minute * time.Duration(*rawSettings.FutureMinutes)
	}

	if rawSettings.ScaleIntervalMinutes != nil {
		s.ScaleInterval = time.Minute * time.Duration(*rawSettings.ScaleIntervalMinutes)
	}

	if rawSettings.BackgroundColorString != nil {
		s.BackgroundColor = s.parseColorString(*rawSettings.BackgroundColorString)
	}

	if rawSettings.MainScaleColorString != nil {
		s.MainScaleColor = s.parseColorString(*rawSettings.MainScaleColorString)
	}

	if rawSettings.SubScalesColorString != nil {
		s.SubScalesColor = s.parseColorString(*rawSettings.SubScalesColorString)
	}

	if rawSettings.ChartColorString != nil {
		s.ChartColor = s.parseColorString(*rawSettings.ChartColorString)
	}

	if rawSettings.TipTextColorString != nil {
		s.TipTextColor = s.parseColorString(*rawSettings.TipTextColorString)
	}

	return nil
}

func (s *Settings) parseColorString(colorString string) winapi.COLORREF {
	var r, g, b int32
	fmt.Sscanf(colorString, "#%02x%02x%02x", &r, &g, &b)
	return winapi.RGB(r, g, b)
}
