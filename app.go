package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"os"
	"os/exec"
	"syscall"
	"time"
	"time-meter/textmap"

	"github.com/cwchiu/go-winapi"
)

type App struct {
	textMap       textmap.TextMap
	settings      *Settings
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
	MID_EDIT_SCHEDULE
	MID_QUIT
)

const MB_TOPMOST = 0x00040000

const SCHEDULE_FILENAME = "schedule.json"

//go:embed embed/text.json
var embedTextJson []byte

func (a *App) Run() error {
	a.textMap = textmap.New()
	if err := a.textMap.LoadJson(bytes.NewReader(embedTextJson)); err != nil {
		return err
	}

	a.settings = new(Settings)
	a.settings.Default()
	if err := a.settings.LoadFile("settings.json"); err != nil {
		println(err.Error())
	}

	a.fileWatcher = new(FileWatcher)
	if err := a.fileWatcher.initialize(); err != nil {
		return err
	}
	defer a.fileWatcher.finalize()

	a.meterWindow = new(MeterWindow)
	a.meterWindow.settings = a.settings
	if err := a.meterWindow.initialize(); err != nil {
		return err
	}
	defer a.meterWindow.finalize()

	a.tipWindow = new(TipWindow)
	a.tipWindow.settings = a.settings
	if err := a.tipWindow.initialize(); err != nil {
		return err
	}
	defer a.tipWindow.finalize()

	a.meterRenderer = new(MeterRenderer)
	a.meterRenderer.settings = a.settings
	if err := a.meterRenderer.initialize(); err != nil {
		return err
	}
	defer a.meterRenderer.finalize()

	a.tipRenderer = new(TipRenderer)
	a.tipRenderer.textMap = a.textMap
	a.tipRenderer.settings = a.settings
	if err := a.tipRenderer.initialize(); err != nil {
		return err
	}
	defer a.tipRenderer.finalize()

	a.contextMenu = new(PopupMenu)
	if err := a.contextMenu.initialize(); err != nil {
		return err
	}
	defer a.contextMenu.finalize()

	a.contextMenu.appendStringItem(MID_EDIT_SCHEDULE, a.textMap.Of("VERB_EDIT_SCHEDULE").String())
	a.contextMenu.appendStringItem(MID_QUIT, a.textMap.Of("VERB_QUIT").String())

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
		totalDuration := a.settings.FutureDuration + a.settings.PastDuration
		focusAt := time.Now().Add(-a.settings.PastDuration + time.Duration(focusRatio*float64(totalDuration)))

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
		case MID_EDIT_SCHEDULE:
			a.openSchedule()

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

func (a *App) openSchedule() {
	cmd := exec.Command(a.settings.ScheduleEditCommand, SCHEDULE_FILENAME)
	if err := cmd.Start(); err != nil {
		captionPtr, _ := syscall.UTF16PtrFromString(a.textMap.Of("NOUN_TIME_METER").String())
		messagePtr, _ := syscall.UTF16PtrFromString(a.textMap.Of("NOTIFY_FAILED_COMMAND_EDIT_SCHEDULE").
			Set("detail", err.Error()).
			String())
		winapi.MessageBox(a.meterWindow.hWnd, messagePtr, captionPtr, winapi.MB_ICONERROR|MB_TOPMOST)
	}
}

func (a *App) reloadSchedule() {
	tasks, err := a.loadTasks(SCHEDULE_FILENAME)
	if err != nil {
		a.tipRenderer.errorMessage = a.textMap.Of("NOTIFY_FAILED_SCHEDULE").
			Set("filename", SCHEDULE_FILENAME).
			String()
		return
	}

	a.tasks = tasks
	a.meterRenderer.tasks = tasks

	a.tipRenderer.errorMessage = ""
}

func (a *App) loadTasks(filename string) ([]Task, error) {
	ret := []Task{}

	if jsonBytes, err := os.ReadFile(filename); err != nil {
		return nil, err

	} else if err := json.NewDecoder(bytes.NewReader(jsonBytes)).Decode(&ret); err != nil {
		return nil, err

	} else {
		return ret, nil
	}
}
