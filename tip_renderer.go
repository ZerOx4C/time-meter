package main

import (
	"math"
	"syscall"
	"time"
	"time-meter/textmap"

	"github.com/cwchiu/go-winapi"
)

type TipRenderer struct {
	textMap              textmap.TextMap
	settings             *Settings
	tasks                []Task
	errorMessage         string
	backgroundBrush      winapi.HBRUSH
	errorBackgroundBrush winapi.HBRUSH
	font                 winapi.HFONT
}

func (tr *TipRenderer) initialize() error {
	tr.backgroundBrush = winapi.CreateSolidBrush(tr.settings.BackgroundColor)
	tr.errorBackgroundBrush = winapi.CreateSolidBrush(winapi.RGB(160, 0, 0))
	tr.font = winapi.CreateFont(
		15, 0, 0, 0, winapi.FW_NORMAL, 0, 0, 0,
		winapi.ANSI_CHARSET, winapi.OUT_DEVICE_PRECIS,
		winapi.CLIP_DEFAULT_PRECIS, winapi.DEFAULT_QUALITY,
		winapi.VARIABLE_PITCH|winapi.FF_SWISS, nil)

	return nil
}

func (tr *TipRenderer) finalize() error {
	winapi.DeleteObject(winapi.HGDIOBJ(tr.backgroundBrush))
	winapi.DeleteObject(winapi.HGDIOBJ(tr.errorBackgroundBrush))
	winapi.DeleteObject(winapi.HGDIOBJ(tr.font))

	return nil
}

func (tr *TipRenderer) draw(hWnd winapi.HWND) {
	var paint winapi.PAINTSTRUCT
	hdc := winapi.BeginPaint(hWnd, &paint)

	backBuffer := new(BackBuffer)
	backDc := backBuffer.begin(hWnd, hdc)

	if tr.errorMessage == "" {
		tr.drawAsTasks(hWnd, backDc, tr.tasks, tr.settings.TipTextColor)

	} else {
		tr.drawAsMessage(hWnd, backDc, tr.errorMessage)
	}

	backBuffer.end()

	winapi.EndPaint(hWnd, &paint)
}

func (tr *TipRenderer) drawAsTasks(hWnd winapi.HWND, hdc winapi.HDC, tasks []Task, tipTextColor winapi.COLORREF) {
	var clientRect RECT
	winapi.GetClientRect(hWnd, clientRect.unwrap())
	winapi.FillRect(hdc, clientRect.unwrap(), tr.backgroundBrush)

	winapi.SetBkMode(hdc, winapi.TRANSPARENT)
	winapi.SelectObject(hdc, winapi.HGDIOBJ(tr.font))
	winapi.SetTextColor(hdc, tipTextColor)

	subjectTextPtr, _ := syscall.UTF16PtrFromString(tr.createSubjectText(tasks))
	timeTextPtr, _ := syscall.UTF16PtrFromString(tr.createTimeText(tasks, time.Now()))

	var subjectRect RECT
	var timeRect RECT
	winapi.DrawText(hdc, subjectTextPtr, -1, subjectRect.unwrap(), winapi.DT_CALCRECT)
	winapi.DrawText(hdc, timeTextPtr, -1, timeRect.unwrap(), winapi.DT_RIGHT|winapi.DT_CALCRECT)

	const (
		PADDING_LEFT   = 5
		PADDING_RIGHT  = 5
		PADDING_TOP    = 2
		PADDING_BOTTOM = 5
		MARGIN         = 10
	)

	subjectRect.translate(PADDING_LEFT, PADDING_TOP)
	timeRect.translate(subjectRect.Right+MARGIN, PADDING_TOP)

	winapi.DrawText(hdc, subjectTextPtr, -1, subjectRect.unwrap(), 0)
	winapi.DrawText(hdc, timeTextPtr, -1, timeRect.unwrap(), winapi.DT_RIGHT)

	winapi.SetWindowPos(
		hWnd, winapi.HWND_TOPMOST,
		0, 0,
		PADDING_LEFT+subjectRect.width()+MARGIN+timeRect.width()+PADDING_RIGHT,
		PADDING_TOP+subjectRect.height()+PADDING_BOTTOM,
		winapi.SWP_NOACTIVATE|winapi.SWP_NOMOVE)
}

func (tr *TipRenderer) createSubjectText(sourceTasks []Task) string {
	var ret string

	for index, task := range sourceTasks {
		if 0 < index {
			ret += "\n"
		}

		ret += task.Subject
	}

	return ret
}

func (tr *TipRenderer) createTimeText(sourceTasks []Task, now time.Time) string {
	var ret string

	for index, task := range sourceTasks {
		if 0 < index {
			ret += "\n"
		}

		if now.Before(task.BeginAt) {
			ret += tr.textMap.Of("INDICATOR_AFTER_MINUTES").
				SetInt("minutes", int(math.Ceil(task.BeginAt.Sub(now).Minutes()))).
				String()

		} else if now.Before(task.EndAt) {
			ret += tr.textMap.Of("INDICATOR_REMAINING_MINUTES").
				SetInt("minutes", int(math.Ceil(task.EndAt.Sub(now).Minutes()))).
				String()
		}
	}

	return ret
}

func (tr *TipRenderer) drawAsMessage(hWnd winapi.HWND, hdc winapi.HDC, message string) {
	var clientRect RECT
	winapi.GetClientRect(hWnd, clientRect.unwrap())
	winapi.FillRect(hdc, clientRect.unwrap(), tr.errorBackgroundBrush)

	winapi.SetBkMode(hdc, winapi.TRANSPARENT)
	winapi.SelectObject(hdc, winapi.HGDIOBJ(tr.font))
	winapi.SetTextColor(hdc, winapi.RGB(255, 255, 255))

	messagePtr, _ := syscall.UTF16PtrFromString(message)

	var messageRect RECT
	winapi.DrawText(hdc, messagePtr, -1, messageRect.unwrap(), winapi.DT_CALCRECT)

	const (
		PADDING_LEFT   = 5
		PADDING_RIGHT  = 5
		PADDING_TOP    = 2
		PADDING_BOTTOM = 5
	)

	messageRect.translate(PADDING_LEFT, PADDING_TOP)

	winapi.DrawText(hdc, messagePtr, -1, messageRect.unwrap(), 0)

	winapi.SetWindowPos(
		hWnd, winapi.HWND_TOPMOST,
		0, 0,
		PADDING_LEFT+messageRect.width()+PADDING_RIGHT,
		PADDING_TOP+messageRect.height()+PADDING_BOTTOM,
		winapi.SWP_NOACTIVATE|winapi.SWP_NOMOVE)
}
