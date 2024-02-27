package main

import (
	"os"
	"time"

	"github.com/cwchiu/go-winapi"
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

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	meterWindow := new(MeterWindow)
	meterWindow.settings = settings
	if err := meterWindow.initialize(); err != nil {
		return err
	}

	tipWindow := new(TipWindow)
	tipWindow.settings = settings
	if err := tipWindow.initialize(); err != nil {
		return err
	}

	meterRenderer := new(MeterRenderer)
	meterRenderer.settings = settings
	if err := meterRenderer.initialize(); err != nil {
		return err
	}

	tipRenderer := new(TipRenderer)
	tipRenderer.settings = settings
	if err := tipRenderer.initialize(); err != nil {
		return err
	}

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

		tipWindow.boundLeft = meterWindow.bound.Right
		tipWindow.update()
	}

	meterWindow.onMouseEnter = func() {
		tipWindow.show()
	}

	meterWindow.onMouseLeave = func() {
		tipWindow.hide()
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

	if err := meterWindow.finalize(); err != nil {
		return err
	}

	if err := tipWindow.finalize(); err != nil {
		return err
	}

	if err := meterRenderer.finalize(); err != nil {
		return err
	}

	if err := tipRenderer.finalize(); err != nil {
		return err
	}

	return nil
}

func loadTasks() ([]Task, error) {
	ret := []Task{}

	var task Task

	task.Subject = "さっき始まったタスク"
	task.BeginAt = time.Now().Add(time.Minute * -10)
	task.EndAt = time.Now().Add(time.Minute * 50)
	ret = append(ret, task)

	task.Subject = "もうすぐ始まるタスク"
	task.BeginAt = time.Now().Add(time.Minute * 10)
	task.EndAt = time.Now().Add(time.Minute * 40)
	ret = append(ret, task)

	task.Subject = "次の次の隣接タスク"
	task.BeginAt = time.Now().Add(time.Minute * 40)
	task.EndAt = time.Now().Add(time.Minute * 70)
	ret = append(ret, task)

	task.Subject = "さっき終わったタスク"
	task.BeginAt = time.Now().Add(time.Minute * -80)
	task.EndAt = time.Now().Add(time.Minute * -20)
	ret = append(ret, task)

	return ret, nil
}
