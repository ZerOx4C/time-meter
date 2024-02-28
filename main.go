package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/cwchiu/go-winapi"
)

type MenuId int16

const (
	MID_ZERO MenuId = iota
	MID_QUIT
)

func main() {
	if err := run(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	settings := new(Settings)
	settings.Default()
	if err := settings.LoadFile("settings.json"); err != nil {
		println(err.Error())
	}

	tasks, err := loadTasks("schedule.json")
	if err != nil {
		return err
	}

	meterWindow := new(MeterWindow)
	meterWindow.settings = settings
	if err := meterWindow.initialize(); err != nil {
		return err
	}
	defer meterWindow.finalize()

	tipWindow := new(TipWindow)
	tipWindow.settings = settings
	if err := tipWindow.initialize(); err != nil {
		return err
	}
	defer tipWindow.finalize()

	meterRenderer := new(MeterRenderer)
	meterRenderer.settings = settings
	if err := meterRenderer.initialize(); err != nil {
		return err
	}
	defer meterRenderer.finalize()

	tipRenderer := new(TipRenderer)
	tipRenderer.settings = settings
	if err := tipRenderer.initialize(); err != nil {
		return err
	}
	defer tipRenderer.finalize()

	contextMenu := new(PopupMenu)
	if err := contextMenu.initialize(); err != nil {
		return err
	}
	defer contextMenu.finalize()

	contextMenu.appendStringItem(MID_QUIT, "終了")

	meterRenderer.tasks = tasks

	meterWindow.onPaint = func() {
		meterRenderer.width = meterWindow.bound.width()
		meterRenderer.height = meterWindow.bound.height()
		meterRenderer.draw(meterWindow.hWnd)
	}

	meterWindow.onMouseMove = func() {
		var cursorPos POINT
		winapi.GetCursorPos(cursorPos.unwrap())

		focusRatio := 1 - float64(cursorPos.Y-meterWindow.bound.Top)/float64(meterWindow.bound.height())
		totalDuration := settings.FutureDuration + settings.PastDuration
		focusAt := time.Now().Add(-settings.PastDuration + time.Duration(focusRatio*float64(totalDuration)))

		var focusTasks []Task
		for _, task := range tasks {
			if task.overlapWith(focusAt, focusAt) {
				focusTasks = append(focusTasks, task)
			}
		}

		tipRenderer.tasks = focusTasks

		if 0 < len(focusTasks) {
			tipWindow.show()

		} else {
			tipWindow.hide()
		}

		tipWindow.boundLeft = meterWindow.bound.Right
		tipWindow.update()
	}

	meterWindow.onMouseEnter = func() {
		tipWindow.show()
	}

	meterWindow.onMouseLeave = func() {
		tipWindow.hide()
	}

	meterWindow.onMouseRightClick = func() {
		contextMenu.popup(meterWindow.hWnd)
	}

	meterWindow.onPopupMenuCommand = func() {
		switch meterWindow.lastMenuId {
		case MID_QUIT:
			winapi.SendMessage(meterWindow.hWnd, winapi.WM_CLOSE, 0, 0)
		}
	}

	tipWindow.onPaint = func() {
		tipRenderer.draw(tipWindow.hWnd)
	}

	meterWindow.show()

	var msg winapi.MSG
	for winapi.GetMessage(&msg, 0, 0, 0) != 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}

	return nil
}

func loadTasks(filename string) ([]Task, error) {
	ret := []Task{}

	if file, err := os.Open(filename); err != nil {
		return nil, err

	} else if err := json.NewDecoder(file).Decode(&ret); err != nil {
		return nil, err

	} else {
		return ret, nil
	}
}
