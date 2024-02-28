package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/cwchiu/go-winapi"
)

type App struct {
	fileWatcher   *FileWatcher
	meterWindow   *MeterWindow
	tipWindow     *TipWindow
	meterRenderer *MeterRenderer
	tipRenderer   *TipRenderer
	contextMenu   *PopupMenu
	tasks         []Task
}

type MenuId int16

const (
	MID_ZERO MenuId = iota
	MID_QUIT
)

const SCHEDULE_FILENAME = "schedule.json"

func (a *App) Run() error {
	settings := new(Settings)
	settings.Default()
	if err := settings.LoadFile("settings.json"); err != nil {
		println(err.Error())
	}

	a.fileWatcher = new(FileWatcher)
	if err := a.fileWatcher.initialize(); err != nil {
		return err
	}
	defer a.fileWatcher.finalize()

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

	a.fileWatcher.onFileChanged = func() {
		a.reloadSchedule()
	}

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
		for _, task := range a.tasks {
			if task.overlapWith(focusAt, focusAt) {
				focusTasks = append(focusTasks, task)
			}
		}

		a.tipRenderer.tasks = focusTasks

		if a.tipRenderer.errorMessage != "" {
			// NOTE: workaround.
			a.tipWindow.show()

		} else if 0 < len(focusTasks) {
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

	a.fileWatcher.filename = SCHEDULE_FILENAME
	a.fileWatcher.watch()

	a.reloadSchedule()
	a.meterWindow.show()

	var msg winapi.MSG
	for winapi.GetMessage(&msg, 0, 0, 0) != 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}

	return nil
}

func (a *App) reloadSchedule() {
	tasks, err := a.loadTasks(SCHEDULE_FILENAME)
	if err != nil {
		a.tipRenderer.errorMessage = "schedule.json の読み込みに失敗しました"
		return
	}

	a.tasks = tasks
	a.meterRenderer.tasks = tasks

	a.tipRenderer.errorMessage = ""
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
