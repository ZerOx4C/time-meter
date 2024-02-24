package main

import (
	"github.com/cwchiu/go-winapi"
)

type App struct {
}

func (a *App) initialize() error {
	meter := new(MeterWindow)

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
