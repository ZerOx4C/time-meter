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
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	meterWindow := new(MeterWindow)
	meterWindow.tasks = tasks

	if err := meterWindow.initialize(); err != nil {
		return err
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
