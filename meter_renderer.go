package main

import (
	"github.com/cwchiu/go-winapi"
)

type MeterRenderer struct {
	hWnd          winapi.HWND
	width         int32
	height        int32
	futureMinutes int32
	pastMinutes   int32
	headPen       winapi.HPEN
	hourPen       winapi.HPEN
	chartBrush    winapi.HBRUSH
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
	minutes := mr.futureMinutes
	totalMinutes := mr.futureMinutes + mr.pastMinutes

	winapi.SelectObject(hdc, winapi.HGDIOBJ(mr.hourPen))

	for 60 < minutes {
		minutes -= 60
	}

	for minutes < totalMinutes {
		if minutes != mr.futureMinutes {
			mr.drawScaleLine(hdc, mr.height*minutes/totalMinutes)
		}
		minutes += 60
	}

	winapi.SelectObject(hdc, winapi.HGDIOBJ(mr.headPen))
	mr.drawScaleLine(hdc, mr.height*mr.futureMinutes/totalMinutes)
}

func (mr *MeterRenderer) drawScaleLine(hdc winapi.HDC, y int32) {
	winapi.MoveToEx(hdc, 0, y, nil)
	winapi.LineTo(hdc, mr.width, y)
}
