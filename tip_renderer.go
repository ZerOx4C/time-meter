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

	backBuffer := new(BackBuffer)
	backDc := backBuffer.begin(hWnd, hdc)

	winapi.SetBkMode(backDc, winapi.TRANSPARENT)
	winapi.SelectObject(backDc, winapi.HGDIOBJ(tr.font))
	winapi.SetTextColor(backDc, winapi.RGB(255, 255, 255))

	subjectTextPtr, _ := syscall.UTF16PtrFromString(tr.createSubjectText(tr.tasks))
	timeTextPtr, _ := syscall.UTF16PtrFromString(tr.createTimeText(tr.tasks, time.Now()))

	var subjectRect RECT
	var timeRect RECT
	winapi.DrawText(backDc, subjectTextPtr, -1, subjectRect.unwrap(), winapi.DT_CALCRECT)
	winapi.DrawText(backDc, timeTextPtr, -1, timeRect.unwrap(), winapi.DT_RIGHT|winapi.DT_CALCRECT)

	const (
		PADDING_LEFT   = 5
		PADDING_RIGHT  = 5
		PADDING_TOP    = 2
		PADDING_BOTTOM = 5
		MARGIN         = 10
	)

	subjectRect.translate(PADDING_LEFT, PADDING_TOP)
	timeRect.translate(subjectRect.Right+MARGIN, PADDING_TOP)

	winapi.DrawText(backDc, subjectTextPtr, -1, subjectRect.unwrap(), 0)
	winapi.DrawText(backDc, timeTextPtr, -1, timeRect.unwrap(), winapi.DT_RIGHT)

	winapi.SetWindowPos(
		hWnd, winapi.HWND_TOPMOST,
		0, 0,
		PADDING_LEFT+subjectRect.width()+MARGIN+timeRect.width()+PADDING_RIGHT,
		PADDING_TOP+subjectRect.height()+PADDING_BOTTOM,
		winapi.SWP_NOACTIVATE|winapi.SWP_NOMOVE)

	backBuffer.end()

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