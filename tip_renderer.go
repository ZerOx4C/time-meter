package main

import (
	"fmt"
	"math"
	"syscall"
	"time"

	"github.com/cwchiu/go-winapi"
)

type TipRenderer struct {
	tasks []Task
	font  winapi.HFONT
}

func (tr *TipRenderer) initialize() error {
	tr.font = winapi.CreateFont(
		15, 0, 0, 0, winapi.FW_NORMAL, 0, 0, 0,
		winapi.ANSI_CHARSET, winapi.OUT_DEVICE_PRECIS,
		winapi.CLIP_DEFAULT_PRECIS, winapi.DEFAULT_QUALITY,
		winapi.VARIABLE_PITCH|winapi.FF_SWISS, nil)

	return nil
}

func (tr *TipRenderer) finalize() error {
	winapi.DeleteObject(winapi.HGDIOBJ(tr.font))

	return nil
}

func (tr *TipRenderer) draw(hWnd winapi.HWND) {
	var paint winapi.PAINTSTRUCT
	hdc := winapi.BeginPaint(hWnd, &paint)

	winapi.SetBkMode(hdc, winapi.TRANSPARENT)
	winapi.SelectObject(hdc, winapi.HGDIOBJ(tr.font))
	winapi.SetTextColor(hdc, winapi.RGB(255, 255, 255))

	subjectTextPtr, _ := syscall.UTF16PtrFromString(tr.createSubjectText(tr.tasks))
	timeTextPtr, _ := syscall.UTF16PtrFromString(tr.createTimeText(tr.tasks, time.Now()))

	var subjectRect winapi.RECT
	var timeRect winapi.RECT
	winapi.DrawText(hdc, subjectTextPtr, -1, &subjectRect, winapi.DT_CALCRECT)
	winapi.DrawText(hdc, timeTextPtr, -1, &timeRect, winapi.DT_RIGHT|winapi.DT_CALCRECT)

	const (
		PADDING_LEFT   = 5
		PADDING_RIGHT  = 5
		PADDING_TOP    = 2
		PADDING_BOTTOM = 5
		MARGIN         = 10
	)

	subjectRect.Left += PADDING_LEFT
	subjectRect.Right += PADDING_LEFT
	subjectRect.Top += PADDING_TOP
	subjectRect.Bottom += PADDING_TOP

	timeRect.Left += subjectRect.Right + MARGIN
	timeRect.Right += subjectRect.Right + MARGIN
	timeRect.Top += PADDING_TOP
	timeRect.Bottom += PADDING_TOP

	winapi.DrawText(hdc, subjectTextPtr, -1, &subjectRect, 0)
	winapi.DrawText(hdc, timeTextPtr, -1, &timeRect, winapi.DT_RIGHT)

	winapi.SetWindowPos(
		hWnd, winapi.HWND_TOPMOST,
		0, 0,
		PADDING_LEFT+subjectRect.Right-subjectRect.Left+MARGIN+timeRect.Right-timeRect.Left+PADDING_RIGHT,
		PADDING_TOP+subjectRect.Bottom-subjectRect.Top+PADDING_BOTTOM,
		winapi.SWP_NOACTIVATE|winapi.SWP_NOMOVE)

	winapi.EndPaint(hWnd, &paint)
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
			ret += fmt.Sprintf("%d分後", int(math.Ceil(task.BeginAt.Sub(now).Minutes())))

		} else if now.Before(task.EndAt) {
			ret += fmt.Sprintf("あと%d分", int(math.Ceil(task.EndAt.Sub(now).Minutes())))
		}
	}

	return ret
}
