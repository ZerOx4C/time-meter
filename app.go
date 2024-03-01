package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
	if err := a.fileWatcher.Initialize(); err != nil {
		return err
	}
	defer a.fileWatcher.Finalize()

	a.meterWindow = new(MeterWindow)
	a.meterWindow.settings = a.settings
	if err := a.meterWindow.Initialize(); err != nil {
		return err
	}
	defer a.meterWindow.Finalize()

	a.tipWindow = new(TipWindow)
	a.tipWindow.settings = a.settings
	if err := a.tipWindow.Initialize(); err != nil {
		return err
	}
	defer a.tipWindow.Finalize()

	a.meterRenderer = new(MeterRenderer)
	a.meterRenderer.settings = a.settings
	if err := a.meterRenderer.Initialize(); err != nil {
		return err
	}
	defer a.meterRenderer.Finalize()

	a.tipRenderer = new(TipRenderer)
	a.tipRenderer.textMap = a.textMap
	a.tipRenderer.settings = a.settings
	if err := a.tipRenderer.Initialize(); err != nil {
		return err
	}
	defer a.tipRenderer.Finalize()

	a.contextMenu = new(PopupMenu)
	if err := a.contextMenu.Initialize(); err != nil {
		return err
	}
	defer a.contextMenu.Finalize()

	a.contextMenu.AppendStringItem(MID_EDIT_SCHEDULE, a.textMap.Of("VERB_EDIT_SCHEDULE").String())
	a.contextMenu.AppendStringItem(MID_QUIT, a.textMap.Of("VERB_QUIT").String())

	a.fileWatcher.onFileChanged = func() {
		a.reloadSchedule()
	}

	a.meterWindow.onPaint = func() {
		a.meterRenderer.width = a.meterWindow.bound.Width()
		a.meterRenderer.height = a.meterWindow.bound.Height()
		a.meterRenderer.Draw(a.meterWindow.hWnd)
	}

	a.meterWindow.onMouseMove = func() {
		var cursorPos POINT
		winapi.GetCursorPos(cursorPos.Unwrap())

		focusRatio := 1 - float64(cursorPos.Y-a.meterWindow.bound.Top)/float64(a.meterWindow.bound.Height())
		totalDuration := a.settings.FutureDuration + a.settings.PastDuration
		focusAt := time.Now().Add(-a.settings.PastDuration + time.Duration(focusRatio*float64(totalDuration)))

		var focusTasks []Task
		for _, task := range a.tasks {
			if task.OverlapWith(focusAt, focusAt) {
				focusTasks = append(focusTasks, task)
			}
		}

		a.tipRenderer.tasks = focusTasks

		if a.tipRenderer.errorMessage != "" {
			// NOTE: workaround.
			a.tipWindow.Show()

		} else if 0 < len(focusTasks) {
			a.tipWindow.Show()

		} else {
			a.tipWindow.Hide()
		}

		a.tipWindow.boundLeft = a.meterWindow.bound.Right
		a.tipWindow.Update()
	}

	a.meterWindow.onMouseEnter = func() {
		a.tipWindow.Show()
	}

	a.meterWindow.onMouseLeave = func() {
		a.tipWindow.Hide()
	}

	a.meterWindow.onMouseRightClick = func() {
		a.contextMenu.Popup(a.meterWindow.hWnd)
	}

	a.meterWindow.onPopupMenuCommand = func() {
		switch a.meterWindow.lastMenuId {
		case MID_EDIT_SCHEDULE:
			if err := a.handleEditSchedule(); err != nil {
				showErrorMessageBox(
					a.meterWindow.hWnd,
					a.textMap.Of("NOUN_TIME_METER").String(),
					a.textMap.Of("NOTIFY_FAILED_OPERATION").
						Set("detail", err.Error()).
						String(),
				)
			}

		case MID_QUIT:
			winapi.SendMessage(a.meterWindow.hWnd, winapi.WM_CLOSE, 0, 0)
		}
	}

	a.tipWindow.onPaint = func() {
		a.tipRenderer.Draw(a.tipWindow.hWnd)
	}

	a.fileWatcher.filename = SCHEDULE_FILENAME
	a.fileWatcher.Watch()

	a.reloadSchedule()
	a.meterWindow.Show()

	var msg winapi.MSG
	for winapi.GetMessage(&msg, 0, 0, 0) != 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}

	return nil
}

func (a *App) handleEditSchedule() error {
	if fileInfo, err := os.Stat(SCHEDULE_FILENAME); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if err := a.saveTemplateTasks(SCHEDULE_FILENAME); err != nil {
			return err
		}

	} else if fileInfo.IsDir() {
		return fmt.Errorf(`"%s" is a directory`, SCHEDULE_FILENAME)
	}

	cmd := exec.Command(a.settings.ScheduleEditCommand, SCHEDULE_FILENAME)
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
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

func (a *App) saveTemplateTasks(filename string) error {
	task := Task{}
	task.Subject = a.textMap.Of("NOUN_SAMPLE_TASK").String()
	task.BeginAt = time.Now().Truncate(time.Minute).Add(time.Minute * 3)
	task.EndAt = task.BeginAt.Add(time.Hour)
	return a.saveTasks(filename, []Task{task})
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

func (a *App) saveTasks(filename string, tasks []Task) error {
	jsonBuffer := bytes.NewBuffer(nil)

	encoder := json.NewEncoder(jsonBuffer)
	encoder.SetIndent("", "\t")
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(tasks); err != nil {
		return err

	} else if err := os.WriteFile(filename, jsonBuffer.Bytes(), os.ModePerm); err != nil {
		return err

	} else {
		return nil
	}
}
