package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
	"time-meter/logic"
	"time-meter/setting"
	"time-meter/textmap"
	"time-meter/ui"
	"time-meter/webapi"
)

const SCHEDULE_FILENAME = "schedule.json"
const SETTINGS_FILENAME = "settings.json"

//go:embed embed/text.json
var embedTextJson []byte

var textMap = textmap.New()
var settings = new(setting.Settings)
var webApi = webapi.New()
var uiController = ui.NewController()
var fileWatcher = new(FileWatcher)

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

	uiController.SetTextMap(textMap)
	uiController.SetSettings(settings)

	if err := initialize(); err != nil {
		return err
	}
	defer finalize()

	fileWatcher.filename = SCHEDULE_FILENAME
	fileWatcher.onFileChanged = func() {
		reloadSchedule()
	}

	webApi.OnHandled(func(t webapi.RequestType) {
		switch t {
		case webapi.PostSchedule:
			if err := logic.SaveTasksFromFile(SCHEDULE_FILENAME, webApi.PostedTasks()); err != nil {
				println(err.Error())
			}
		}
	})

	fileWatcher.Watch()

	if settings.ServerEnabled {
		mux := http.NewServeMux()
		mux.Handle("/api", http.StripPrefix("/api", webApi))
		go http.ListenAndServe(fmt.Sprintf(":%d", settings.Port), mux)
	}

	reloadSchedule()

	uiController.OnPopupMenuCommand(func(menuId ui.MenuId) {
		switch menuId {
		case ui.MID_EDIT_SCHEDULE:
			if err := handleEditSchedule(); err != nil {
				uiController.ShowErrorMessageBox(
					textMap.Of("NOTIFY_FAILED_OPERATION").
						Set("detail", err.Error()).
						String())
			}

		case ui.MID_QUIT:
			uiController.Quit()
		}
	})

	uiController.Run()

	return nil
}

func initialize() error {
	if err := textMap.LoadJson(bytes.NewReader(embedTextJson)); err != nil {
		return err
	}

	if err := fileWatcher.Initialize(); err != nil {
		return err
	}

	if err := uiController.Initialize(); err != nil {
		return err
	}

	return nil
}

func finalize() {
	uiController.Finalize()
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
	loadedTasks, err := logic.LoadTasksFromFile(SCHEDULE_FILENAME)
	if err != nil {
		uiController.SetErrorMessage(textMap.Of("NOTIFY_FAILED_SCHEDULE").
			Set("filename", SCHEDULE_FILENAME).
			String())
		return
	}

	uiController.SetTasks(loadedTasks)

	uiController.SetErrorMessage("")
}

func saveTemplateTasks(filename string) error {
	task := logic.Task{}
	task.Subject = textMap.Of("NOUN_SAMPLE_TASK").String()
	task.BeginAt = time.Now().Truncate(time.Minute).Add(time.Minute * 3)
	task.EndAt = task.BeginAt.Add(time.Hour)
	return logic.SaveTasksFromFile(filename, []logic.Task{task})
}
