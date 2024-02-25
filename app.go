package main

import (
	"time"

	"github.com/cwchiu/go-winapi"
)

type App struct {
}

func (a *App) initialize() error {
	tasks, err := a.loadTasks()
	if err != nil {
		return err
	}

	meter := new(MeterWindow)
	meter.tasks = tasks

	if err := meter.initialize(); err != nil {
		return err
	}

	meter.show()

	var msg winapi.MSG
	for winapi.GetMessage(&msg, 0, 0, 0) != 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}

	meter.finalize()

	return nil
}

func (a *App) loadTasks() ([]Task, error) {
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
