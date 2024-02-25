package main

import (
	"time"

	"github.com/cwchiu/go-winapi"
)

type MeterRenderer struct {
	hWnd           winapi.HWND
	width          int32
	height         int32
	futureDuration time.Duration
	pastDuration   time.Duration
	headPen        winapi.HPEN
	hourPen        winapi.HPEN
	chartBrush     winapi.HBRUSH
}

func (mr *MeterRenderer) initialize() error {
	mr.headPen = winapi.CreatePen(winapi.PS_SOLID, 1, winapi.RGB(255, 255, 255))
	mr.hourPen = winapi.CreatePen(winapi.PS_SOLID, 1, winapi.RGB(64, 64, 64))
	mr.chartBrush = winapi.CreateSolidBrush(winapi.RGB(255, 128, 0))

	return nil
}

func (mr *MeterRenderer) finalize() error {
	winapi.DeleteObject(winapi.HGDIOBJ(mr.headPen))
	winapi.DeleteObject(winapi.HGDIOBJ(mr.hourPen))
	winapi.DeleteObject(winapi.HGDIOBJ(mr.chartBrush))

	return nil
}

func (mr *MeterRenderer) draw() {
	var paint winapi.PAINTSTRUCT
	hdc := winapi.BeginPaint(mr.hWnd, &paint)

	mr.drawAllScaleLines(hdc)

	winapi.EndPaint(mr.hWnd, &paint)
}

func (mr *MeterRenderer) drawAllScaleLines(hdc winapi.HDC) {
	offset := mr.futureDuration
	totalDuration := mr.futureDuration + mr.pastDuration
	totalSeconds := int32(totalDuration / time.Second)

	winapi.SelectObject(hdc, winapi.HGDIOBJ(mr.hourPen))

	for time.Hour < offset {
		offset -= time.Hour
	}

	for offset < totalDuration {
		if offset != mr.futureDuration {
			mr.drawScaleLine(hdc, mr.height*int32(offset/time.Second)/totalSeconds)
		}
		offset += time.Hour
	}

	winapi.SelectObject(hdc, winapi.HGDIOBJ(mr.headPen))
	mr.drawScaleLine(hdc, mr.height*int32(mr.futureDuration/time.Second)/totalSeconds)
}

func (mr *MeterRenderer) drawScaleLine(hdc winapi.HDC, y int32) {
	winapi.MoveToEx(hdc, 0, y, nil)
	winapi.LineTo(hdc, mr.width, y)
}
