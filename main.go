package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
	"time-meter/setting"
	"time-meter/textmap"
	"time-meter/wrapped"

	"github.com/cwchiu/go-winapi"
)

type MenuId int16

const (
	MID_ZERO MenuId = iota
	MID_EDIT_SCHEDULE
	MID_QUIT
)

const SCHEDULE_FILENAME = "schedule.json"
const SETTINGS_FILENAME = "settings.json"

//go:embed embed/text.json
var embedTextJson []byte

var textMap = textmap.New()
var settings = new(setting.Settings)
var fileWatcher = new(FileWatcher)
var meterWindow = new(MeterWindow)
var tipWindow = new(TipWindow)
var meterRenderer = new(MeterRenderer)
var tipRenderer = new(TipRenderer)
var contextMenu = new(PopupMenu)
var tasks = []Task{}

func main() {
	if err := run(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	settings.Default()
	if err := settings.LoadFile(SETTINGS_FILENAME); err != nil {
		println(err.Error())
	}

	meterWindow.settings = settings
	tipWindow.settings = settings
	meterRenderer.settings = settings
	tipRenderer.textMap = textMap
	tipRenderer.settings = settings

	if err := initialize(); err != nil {
		return err
	}
	defer finalize()

	contextMenu.AppendStringItem(MID_EDIT_SCHEDULE, textMap.Of("VERB_EDIT_SCHEDULE").String())
	contextMenu.AppendStringItem(MID_QUIT, textMap.Of("VERB_QUIT").String())

	fileWatcher.filename = SCHEDULE_FILENAME
	fileWatcher.onFileChanged = func() {
		reloadSchedule()
	}

	meterWindow.onPaint = func() {
		meterRenderer.width = meterWindow.bound.Width()
		meterRenderer.height = meterWindow.bound.Height()
		meterRenderer.Draw(meterWindow.hWnd)
	}

	meterWindow.onMouseMove = func() {
		var cursorPos wrapped.POINT
		winapi.GetCursorPos(cursorPos.Unwrap())

		focusRatio := 1 - float64(cursorPos.Y-meterWindow.bound.Top)/float64(meterWindow.bound.Height())
		totalDuration := settings.FutureDuration + settings.PastDuration
		focusAt := time.Now().Add(-settings.PastDuration + time.Duration(focusRatio*float64(totalDuration)))

		var focusTasks []Task
		for _, task := range tasks {
			if task.OverlapWith(focusAt, focusAt) {
				focusTasks = append(focusTasks, task)
			}
		}

		tipRenderer.tasks = focusTasks

		if tipRenderer.errorMessage != "" {
			// NOTE: workaround.
			tipWindow.Show()

		} else if 0 < len(focusTasks) {
			tipWindow.Show()

		} else {
			tipWindow.Hide()
		}

		tipWindow.boundLeft = meterWindow.bound.Right
		tipWindow.Update()
	}

	meterWindow.onMouseEnter = func() {
		tipWindow.Show()
	}

	meterWindow.onMouseLeave = func() {
		tipWindow.Hide()
	}

	meterWindow.onMouseRightClick = func() {
		contextMenu.Popup(meterWindow.hWnd)
	}

	meterWindow.onPopupMenuCommand = func() {
		switch meterWindow.lastMenuId {
		case MID_EDIT_SCHEDULE:
			if err := handleEditSchedule(); err != nil {
				showErrorMessageBox(
					meterWindow.hWnd,
					textMap.Of("NOUN_TIME_METER").String(),
					textMap.Of("NOTIFY_FAILED_OPERATION").
						Set("detail", err.Error()).
						String(),
				)
			}

		case MID_QUIT:
			winapi.SendMessage(meterWindow.hWnd, winapi.WM_CLOSE, 0, 0)
		}
	}

	tipWindow.onPaint = func() {
		tipRenderer.Draw(tipWindow.hWnd)
	}

	fileWatcher.Watch()

	reloadSchedule()
	meterWindow.Show()

	var msg winapi.MSG
	for winapi.GetMessage(&msg, 0, 0, 0) != 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}

	return nil
}

func initialize() error {
	if err := textMap.LoadJson(bytes.NewReader(embedTextJson)); err != nil {
		return err
	}

	if err := fileWatcher.Initialize(); err != nil {
		return err
	}

	if err := meterWindow.Initialize(); err != nil {
		return err
	}

	if err := tipWindow.Initialize(); err != nil {
		return err
	}

	if err := meterRenderer.Initialize(); err != nil {
		return err
	}

	if err := tipRenderer.Initialize(); err != nil {
		return err
	}

	if err := contextMenu.Initialize(); err != nil {
		return err
	}

	return nil
}

func finalize() {
	contextMenu.Finalize()
	tipRenderer.Finalize()
	meterRenderer.Finalize()
	tipWindow.Finalize()
	meterWindow.Finalize()
	fileWatcher.Finalize()
}

func handleEditSchedule() error {
	if fileInfo, err := os.Stat(SCHEDULE_FILENAME); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if err := saveTemplateTasks(SCHEDULE_FILENAME); err != nil {
			return err
		}

	} else if fileInfo.IsDir() {
		return fmt.Errorf(`"%s" is a directory`, SCHEDULE_FILENAME)
	}

	cmd := exec.Command(settings.ScheduleEditCommand, SCHEDULE_FILENAME)
	if err := cmd.Start(); err != nil {
		return err
	}

	return nil
}

func reloadSchedule() {
	loadedTasks, err := loadTasks(SCHEDULE_FILENAME)
	if err != nil {
		tipRenderer.errorMessage = textMap.Of("NOTIFY_FAILED_SCHEDULE").
			Set("filename", SCHEDULE_FILENAME).
			String()
		return
	}

	tasks = loadedTasks
	meterRenderer.tasks = loadedTasks

	tipRenderer.errorMessage = ""
}

func saveTemplateTasks(filename string) error {
	task := Task{}
	task.Subject = textMap.Of("NOUN_SAMPLE_TASK").String()
	task.BeginAt = time.Now().Truncate(time.Minute).Add(time.Minute * 3)
	task.EndAt = task.BeginAt.Add(time.Hour)
	return saveTasks(filename, []Task{task})
}

func loadTasks(filename string) ([]Task, error) {
	ret := []Task{}

	if jsonBytes, err := os.ReadFile(filename); err != nil {
		return nil, err

	} else if err := json.NewDecoder(bytes.NewReader(jsonBytes)).Decode(&ret); err != nil {
		return nil, err

	} else {
		return ret, nil
	}
}

func saveTasks(filename string, tasks []Task) error {
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
