package main

import (
	"bytes"
	"encoding/json"
	"os"
	"time"
)

type Settings struct {
	MeterWidth     int
	MeterOpacity   byte
	PastDuration   time.Duration
	FutureDuration time.Duration
	ScaleInterval  time.Duration
}

func (s *Settings) Default() {
	s.MeterWidth = 50
	s.MeterOpacity = 128
	s.PastDuration = time.Hour * 1
	s.FutureDuration = time.Hour * 3
	s.ScaleInterval = time.Hour * 1
}

func (s *Settings) LoadFile(filename string) error {
	var rawSettings struct {
		MeterWidth           *int  `json:"meter_width,omitempty"`
		MeterOpacity         *byte `json:"meter_opacity,omitempty"`
		PastMinutes          *int  `json:"past_minutes,omitempty"`
		FutureMinutes        *int  `json:"future_minutes,omitempty"`
		ScaleIntervalMinutes *int  `json:"scale_interval_minutes,omitempty"`
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

	return nil
}
