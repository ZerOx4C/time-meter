package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/cwchiu/go-winapi"
)

type App struct {
	meterWindow   *MeterWindow
	tipWindow     *TipWindow
	meterRenderer *MeterRenderer
	tipRenderer   *TipRenderer
	contextMenu   *PopupMenu
}

type MenuId int16

const (
	MID_ZERO MenuId = iota
	MID_QUIT
)

func (a *App) Run() error {
	settings := new(Settings)
	settings.Default()
	if err := settings.LoadFile("settings.json"); err != nil {
		println(err.Error())
	}

	tasks, err := a.loadTasks("schedule.json")
	if err != nil {
		return err
	}

	a.meterWindow = new(MeterWindow)
	a.meterWindow.settings = settings
	if err := a.meterWindow.initialize(); err != nil {
		return err
	}
	defer a.meterWindow.finalize()

	a.tipWindow = new(TipWindow)
	a.tipWindow.settings = settings
	if err := a.tipWindow.initialize(); err != nil {
		return err
	}
	defer a.tipWindow.finalize()

	a.meterRenderer = new(MeterRenderer)
	a.meterRenderer.settings = settings
	if err := a.meterRenderer.initialize(); err != nil {
		return err
	}
	defer a.meterRenderer.finalize()

	a.tipRenderer = new(TipRenderer)
	a.tipRenderer.settings = settings
	if err := a.tipRenderer.initialize(); err != nil {
		return err
	}
	defer a.tipRenderer.finalize()

	a.contextMenu = new(PopupMenu)
	if err := a.contextMenu.initialize(); err != nil {
		return err
	}
	defer a.contextMenu.finalize()

	a.contextMenu.appendStringItem(MID_QUIT, "終了")

	a.meterRenderer.tasks = tasks

	a.meterWindow.onPaint = func() {
		a.meterRenderer.width = a.meterWindow.bound.width()
		a.meterRenderer.height = a.meterWindow.bound.height()
		a.meterRenderer.draw(a.meterWindow.hWnd)
	}

	a.meterWindow.onMouseMove = func() {
		var cursorPos POINT
		winapi.GetCursorPos(cursorPos.unwrap())

		focusRatio := 1 - float64(cursorPos.Y-a.meterWindow.bound.Top)/float64(a.meterWindow.bound.height())
		totalDuration := settings.FutureDuration + settings.PastDuration
		focusAt := time.Now().Add(-settings.PastDuration + time.Duration(focusRatio*float64(totalDuration)))

		var focusTasks []Task
		for _, task := range tasks {
			if task.overlapWith(focusAt, focusAt) {
				focusTasks = append(focusTasks, task)
			}
		}

		a.tipRenderer.tasks = focusTasks

		if 0 < len(focusTasks) {
			a.tipWindow.show()

		} else {
			a.tipWindow.hide()
		}

		a.tipWindow.boundLeft = a.meterWindow.bound.Right
		a.tipWindow.update()
	}

	a.meterWindow.onMouseEnter = func() {
		a.tipWindow.show()
	}

	a.meterWindow.onMouseLeave = func() {
		a.tipWindow.hide()
	}

	a.meterWindow.onMouseRightClick = func() {
		a.contextMenu.popup(a.meterWindow.hWnd)
	}

	a.meterWindow.onPopupMenuCommand = func() {
		switch a.meterWindow.lastMenuId {
		case MID_QUIT:
			winapi.SendMessage(a.meterWindow.hWnd, winapi.WM_CLOSE, 0, 0)
		}
	}

	a.tipWindow.onPaint = func() {
		a.tipRenderer.draw(a.tipWindow.hWnd)
	}

	a.meterWindow.show()

	var msg winapi.MSG
	for winapi.GetMessage(&msg, 0, 0, 0) != 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}

	return nil
}

func (a *App) loadTasks(filename string) ([]Task, error) {
	ret := []Task{}

	if file, err := os.Open(filename); err != nil {
		return nil, err

	} else if err := json.NewDecoder(file).Decode(&ret); err != nil {
		return nil, err

	} else {
		return ret, nil
	}
}
